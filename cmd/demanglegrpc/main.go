// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

// Command demanglegrpc runs the library behind a gRPC service.
//
// Stage 6 ships this binary; Stage 6.5 adds the lux deploy, systemd
// unit, healthz, and Prometheus metrics when a concrete non-Go
// non-skynet caller materialises. Until then, this binary is
// in-repo + in-tests only.
//
// Usage:
//
//	demanglegrpc --listen 127.0.0.1:50061 [--context-db /var/lib/demangle/contexts.db]
//
// For development + tests, pass --listen 127.0.0.1:0 (choose-a-port).
package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"

	"google.golang.org/grpc"

	"github.com/jelmerdehen/demangle"
	pb "github.com/jelmerdehen/demangle/cmd/demanglegrpc/proto/demanglepb"

	_ "github.com/jelmerdehen/demangle/scheme/all"
)

func main() {
	listen := flag.String("listen", "127.0.0.1:50061", "address to listen on")
	storePath := flag.String("context-db", "", "path to context SQLite DB; defaults to temp-file")
	flag.Parse()

	if *storePath == "" {
		*storePath = filepath.Join(os.TempDir(), "demangle-contexts.db")
	}

	store, err := demangle.OpenContextStore(*storePath)
	if err != nil {
		log.Fatalf("demangle: open context store: %v", err)
	}
	defer store.Close()

	lis, err := net.Listen("tcp", *listen)
	if err != nil {
		log.Fatalf("demangle: listen: %v", err)
	}

	srv := grpc.NewServer()
	pb.RegisterDemangleServer(srv, newService(demangle.Default, store))

	fmt.Fprintf(os.Stderr, "demanglegrpc: listening on %s\n", lis.Addr())
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("demangle: serve: %v", err)
	}
}
