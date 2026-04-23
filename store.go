// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package demangle

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	_ "modernc.org/sqlite"
)

// OpenContextStore opens a SQLite-backed ContextStore at path. Creates
// the file + schema if it doesn't exist. Applies the standard pragmas
// (WAL, NORMAL, foreign_keys, busy_timeout, mmap, cache_size) on every
// connection.
//
// path == ":memory:" opens an in-memory database (single process, not
// across goroutines unless you also pass ?cache=shared). For tests
// prefer InMemoryContextStore which hides that detail.
func OpenContextStore(path string) (ContextStore, error) {
	dsn := path
	if !strings.Contains(dsn, "?") {
		dsn += "?_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)&_pragma=foreign_keys(ON)&_pragma=busy_timeout(5000)&_pragma=temp_store(MEMORY)&_pragma=mmap_size(268435456)&_pragma=cache_size(-32768)"
	}
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("demangle: open context store: %w", err)
	}
	db.SetMaxOpenConns(8)
	db.SetMaxIdleConns(4)
	db.SetConnMaxIdleTime(5 * time.Minute)

	if err := initStoreSchema(db); err != nil {
		db.Close()
		return nil, err
	}
	return &sqliteStore{db: db, writerLock: &sync.Mutex{}}, nil
}

// InMemoryContextStore returns an ephemeral ContextStore that lives
// entirely in process memory. Ideal for tests and short-lived
// subprocesses (Lambda, CI runners).
func InMemoryContextStore() ContextStore {
	return &memStore{blobs: map[string]*memEntry{}}
}

// --- sqlite-backed store -----------------------------------------

type sqliteStore struct {
	db         *sql.DB
	writerLock *sync.Mutex // serialise writers; WAL handles reader concurrency
}

const storeSchema = `
CREATE TABLE IF NOT EXISTS contexts (
    sha256         TEXT PRIMARY KEY,
    kind           TEXT NOT NULL,
    byte_size      INTEGER NOT NULL,
    blob           BLOB NOT NULL,
    metadata_json  TEXT NOT NULL DEFAULT '{}',
    uploaded_ts    INTEGER NOT NULL,
    last_access_ts INTEGER NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_contexts_kind ON contexts(kind);
CREATE INDEX IF NOT EXISTS idx_contexts_last_access ON contexts(last_access_ts);
`

func initStoreSchema(db *sql.DB) error {
	for _, stmt := range splitStatements(storeSchema) {
		if _, err := db.Exec(stmt); err != nil {
			return fmt.Errorf("demangle: init schema: %w", err)
		}
	}
	return nil
}

func splitStatements(s string) []string {
	parts := strings.Split(s, ";")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}

func (s *sqliteStore) Put(ctx context.Context, kind string, blob []byte, meta map[string]string) (string, error) {
	sum := sha256.Sum256(blob)
	shaHex := hex.EncodeToString(sum[:])
	metaJSON, err := marshalMeta(meta)
	if err != nil {
		return "", err
	}
	now := time.Now().Unix()

	s.writerLock.Lock()
	defer s.writerLock.Unlock()

	// Idempotent: duplicate sha256 updates last_access_ts only.
	res, err := s.db.ExecContext(ctx, `
		INSERT INTO contexts(sha256, kind, byte_size, blob, metadata_json, uploaded_ts, last_access_ts)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(sha256) DO UPDATE SET last_access_ts = excluded.last_access_ts
	`, shaHex, kind, len(blob), blob, metaJSON, now, now)
	if err != nil {
		return "", fmt.Errorf("demangle: put context: %w", err)
	}
	_ = res
	return shaHex, nil
}

func (s *sqliteStore) Get(ctx context.Context, sha string) (Context, error) {
	var (
		kind     string
		blob     []byte
		metaJSON string
	)
	row := s.db.QueryRowContext(ctx, `
		SELECT kind, blob, metadata_json FROM contexts WHERE sha256 = ?
	`, sha)
	if err := row.Scan(&kind, &blob, &metaJSON); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &Error{Kind: ErrNeedsContext, Expected: "stored context", Got: sha, Offset: -1}
		}
		return nil, fmt.Errorf("demangle: get context: %w", err)
	}
	// Update last_access_ts asynchronously to keep the read fast.
	go func() { _ = s.Touch(context.Background(), sha) }()

	meta, err := unmarshalMeta(metaJSON)
	if err != nil {
		return nil, err
	}
	return &blobContext{kind: kind, sha: sha, blob: blob, meta: meta}, nil
}

func (s *sqliteStore) List(ctx context.Context, kind string) ([]ContextInfo, error) {
	var rows *sql.Rows
	var err error
	if kind == "" {
		rows, err = s.db.QueryContext(ctx, `
			SELECT kind, sha256, byte_size, metadata_json, uploaded_ts, last_access_ts
			FROM contexts ORDER BY last_access_ts DESC
		`)
	} else {
		rows, err = s.db.QueryContext(ctx, `
			SELECT kind, sha256, byte_size, metadata_json, uploaded_ts, last_access_ts
			FROM contexts WHERE kind = ? ORDER BY last_access_ts DESC
		`, kind)
	}
	if err != nil {
		return nil, fmt.Errorf("demangle: list contexts: %w", err)
	}
	defer rows.Close()

	out := make([]ContextInfo, 0, 16)
	for rows.Next() {
		var (
			info       ContextInfo
			metaJSON   string
			uploaded   int64
			lastAccess int64
		)
		if err := rows.Scan(&info.Kind, &info.SHA256, &info.ByteSize, &metaJSON, &uploaded, &lastAccess); err != nil {
			return nil, err
		}
		info.UploadedTS = time.Unix(uploaded, 0)
		info.LastAccessTS = time.Unix(lastAccess, 0)
		info.Metadata, _ = unmarshalMeta(metaJSON)
		out = append(out, info)
	}
	return out, rows.Err()
}

func (s *sqliteStore) Delete(ctx context.Context, sha string) error {
	s.writerLock.Lock()
	defer s.writerLock.Unlock()
	_, err := s.db.ExecContext(ctx, `DELETE FROM contexts WHERE sha256 = ?`, sha)
	return err
}

func (s *sqliteStore) Touch(ctx context.Context, sha string) error {
	s.writerLock.Lock()
	defer s.writerLock.Unlock()
	_, err := s.db.ExecContext(ctx, `UPDATE contexts SET last_access_ts = ? WHERE sha256 = ?`, time.Now().Unix(), sha)
	return err
}

func (s *sqliteStore) Close() error { return s.db.Close() }

// --- in-memory store ---------------------------------------------

type memStore struct {
	mu    sync.Mutex
	blobs map[string]*memEntry
}

type memEntry struct {
	kind         string
	blob         []byte
	meta         map[string]string
	uploadedTS   time.Time
	lastAccessTS time.Time
}

func (m *memStore) Put(_ context.Context, kind string, blob []byte, meta map[string]string) (string, error) {
	sum := sha256.Sum256(blob)
	sha := hex.EncodeToString(sum[:])
	now := time.Now()
	m.mu.Lock()
	defer m.mu.Unlock()
	if existing, ok := m.blobs[sha]; ok {
		existing.lastAccessTS = now
		return sha, nil
	}
	metaCopy := map[string]string{}
	for k, v := range meta {
		metaCopy[k] = v
	}
	m.blobs[sha] = &memEntry{
		kind: kind, blob: append([]byte(nil), blob...),
		meta: metaCopy, uploadedTS: now, lastAccessTS: now,
	}
	return sha, nil
}

func (m *memStore) Get(_ context.Context, sha string) (Context, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	e, ok := m.blobs[sha]
	if !ok {
		return nil, &Error{Kind: ErrNeedsContext, Expected: "stored context", Got: sha, Offset: -1}
	}
	e.lastAccessTS = time.Now()
	metaCopy := map[string]string{}
	for k, v := range e.meta {
		metaCopy[k] = v
	}
	return &blobContext{kind: e.kind, sha: sha, blob: append([]byte(nil), e.blob...), meta: metaCopy}, nil
}

func (m *memStore) List(_ context.Context, kind string) ([]ContextInfo, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]ContextInfo, 0, len(m.blobs))
	for sha, e := range m.blobs {
		if kind != "" && e.kind != kind {
			continue
		}
		metaCopy := map[string]string{}
		for k, v := range e.meta {
			metaCopy[k] = v
		}
		out = append(out, ContextInfo{
			Kind: e.kind, SHA256: sha, ByteSize: int64(len(e.blob)),
			UploadedTS: e.uploadedTS, LastAccessTS: e.lastAccessTS,
			Metadata: metaCopy,
		})
	}
	return out, nil
}

func (m *memStore) Delete(_ context.Context, sha string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.blobs, sha)
	return nil
}

func (m *memStore) Touch(_ context.Context, sha string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if e, ok := m.blobs[sha]; ok {
		e.lastAccessTS = time.Now()
	}
	return nil
}

func (m *memStore) Close() error { return nil }

// --- blobContext (returned by Get) -------------------------------

type blobContext struct {
	kind string
	sha  string
	blob []byte
	meta map[string]string
}

func (b *blobContext) Kind() string                { return b.kind }
func (b *blobContext) SHA256() string              { return b.sha }
func (b *blobContext) Metadata() map[string]string { return b.meta }
func (b *blobContext) Lookup(key string) (string, bool) {
	// Default Lookup: scan metadata. Specialised context types wrap
	// blobContext + override Lookup to parse the blob.
	v, ok := b.meta[key]
	return v, ok
}
func (b *blobContext) Reader() (io.ReadCloser, error) {
	return io.NopCloser(bytesReader(b.blob)), nil
}

// Blob exposes the stored bytes for schemes that need to parse them
// (e.g. ProGuard map parser, JS source map parser). Callers MUST NOT
// mutate the returned slice; the store guarantees a defensive copy
// was made on retrieval.
func (b *blobContext) Blob() []byte { return b.blob }

type bytesReader []byte

func (b bytesReader) Read(p []byte) (int, error) {
	if len(b) == 0 {
		return 0, io.EOF
	}
	n := copy(p, b)
	return n, nil
}

// --- helpers -----------------------------------------------------

func marshalMeta(m map[string]string) (string, error) {
	if m == nil {
		return "{}", nil
	}
	b, err := json.Marshal(m)
	if err != nil {
		return "", fmt.Errorf("demangle: marshal metadata: %w", err)
	}
	return string(b), nil
}

func unmarshalMeta(s string) (map[string]string, error) {
	if s == "" {
		return map[string]string{}, nil
	}
	var m map[string]string
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		return nil, fmt.Errorf("demangle: unmarshal metadata: %w", err)
	}
	return m, nil
}
