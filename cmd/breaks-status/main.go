// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen
//
// breaks-status parses `breaks.log` and reports outstanding BREAK_OK
// entries with deadline status.
//
// breaks.log entry format (append-only):
//
//	## <BREAK_ID> opened <RFC3339> by commit <sha>
//	reason: <text>
//	restore_by: 2026-05-13
//	disappeared_count: 87
//	disappeared_set: <newline-separated symbols, indented 4 spaces>
//	---END
//
// BREAK_FIXED appends a single line:
//
//	## <BREAK_ID> closed <RFC3339> by commit <sha>
//
// BREAK_EXTEND appends:
//
//	## <BREAK_ID> extended <RFC3339> restore_by <new-iso-date>
//
// Exit non-zero if any open break has expired (today > restore_by).
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

type breakEntry struct {
	id            string
	openedAt      time.Time
	openedCommit  string
	reason        string
	restoreBy     time.Time
	closed        bool
	closedAt      time.Time
	closedCommit  string
	extensions    []time.Time
	disappeared   int
}

var (
	reOpened   = regexp.MustCompile(`^## (\S+) opened (\S+) by commit (\S+)$`)
	reClosed   = regexp.MustCompile(`^## (\S+) closed (\S+) by commit (\S+)$`)
	reExtended = regexp.MustCompile(`^## (\S+) extended (\S+) restore_by (\S+)$`)
)

func parseBreaksLog(path string) ([]*breakEntry, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()
	byID := map[string]*breakEntry{}
	var current *breakEntry
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
	for sc.Scan() {
		line := sc.Text()
		if m := reOpened.FindStringSubmatch(line); m != nil {
			ts, _ := time.Parse(time.RFC3339, m[2])
			be := &breakEntry{id: m[1], openedAt: ts, openedCommit: m[3]}
			byID[m[1]] = be
			current = be
			continue
		}
		if m := reClosed.FindStringSubmatch(line); m != nil {
			be, ok := byID[m[1]]
			if !ok {
				continue
			}
			ts, _ := time.Parse(time.RFC3339, m[2])
			be.closed = true
			be.closedAt = ts
			be.closedCommit = m[3]
			current = be
			continue
		}
		if m := reExtended.FindStringSubmatch(line); m != nil {
			be, ok := byID[m[1]]
			if !ok {
				continue
			}
			ts, _ := time.Parse(time.RFC3339, m[2])
			rb, _ := time.Parse("2006-01-02", m[3])
			be.extensions = append(be.extensions, ts)
			be.restoreBy = rb
			current = be
			continue
		}
		if current == nil {
			continue
		}
		if strings.HasPrefix(line, "reason:") {
			current.reason = strings.TrimSpace(strings.TrimPrefix(line, "reason:"))
		}
		if strings.HasPrefix(line, "restore_by:") {
			rb := strings.TrimSpace(strings.TrimPrefix(line, "restore_by:"))
			t, err := time.Parse("2006-01-02", rb)
			if err == nil {
				current.restoreBy = t
			}
		}
		if strings.HasPrefix(line, "disappeared_count:") {
			fmt.Sscanf(strings.TrimSpace(strings.TrimPrefix(line, "disappeared_count:")), "%d", &current.disappeared)
		}
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}

	out := make([]*breakEntry, 0, len(byID))
	for _, be := range byID {
		out = append(out, be)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].openedAt.Before(out[j].openedAt)
	})
	return out, nil
}

func main() {
	repoRoot := flag.String("repo", "/data/p/demangle", "repo root")
	flag.Parse()

	path := filepath.Join(*repoRoot, "breaks.log")
	entries, err := parseBreaksLog(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse breaks.log: %v\n", err)
		os.Exit(2)
	}

	if len(entries) == 0 {
		fmt.Println("breaks-status: no breaks recorded")
		os.Exit(0)
	}

	now := time.Now().UTC().Truncate(24 * time.Hour)
	openCount, expiredCount := 0, 0
	fmt.Printf("breaks-status: %d total entries\n\n", len(entries))
	for _, be := range entries {
		state := "open"
		if be.closed {
			state = "closed " + be.closedAt.Format("2006-01-02")
		}
		daysRemaining := 0
		expired := false
		if !be.closed {
			openCount++
			if !be.restoreBy.IsZero() {
				daysRemaining = int(be.restoreBy.Sub(now).Hours() / 24)
				if daysRemaining < 0 {
					expired = true
					expiredCount++
				}
			}
		}
		marker := "  "
		if expired {
			marker = "❌"
		} else if !be.closed {
			marker = "🟡"
		} else {
			marker = "✓ "
		}
		fmt.Printf("%s %s — %s\n", marker, be.id, state)
		fmt.Printf("    opened: %s by %s\n", be.openedAt.Format("2006-01-02"), be.openedCommit)
		fmt.Printf("    reason: %s\n", be.reason)
		if !be.closed {
			fmt.Printf("    restore_by: %s (%d days remaining)\n",
				be.restoreBy.Format("2006-01-02"), daysRemaining)
		}
		if be.disappeared > 0 {
			fmt.Printf("    disappeared_count: %d\n", be.disappeared)
		}
		if len(be.extensions) > 0 {
			fmt.Printf("    extensions: %d\n", len(be.extensions))
		}
		fmt.Println()
	}
	fmt.Printf("summary: %d open, %d expired, %d closed\n",
		openCount, expiredCount, len(entries)-openCount)
	if expiredCount > 0 {
		os.Exit(1)
	}
}
