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
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"

	"github.com/jelmerdehen/demangle"
	pb "github.com/jelmerdehen/demangle/cmd/demanglegrpc/proto/demanglepb"

	_ "github.com/jelmerdehen/demangle/scheme/all"
)

func main() {
	listen := flag.String("listen", "127.0.0.1:50061", "gRPC address to listen on")
	healthListen := flag.String("health-listen", "127.0.0.1:50062", "HTTP health/metrics address")
	storePath := flag.String("context-db", "", "path to context SQLite DB; defaults to temp-file")
	tlsCert := flag.String("tls-cert", "", "TLS certificate file (PEM); empty = plaintext")
	tlsKey := flag.String("tls-key", "", "TLS key file (PEM); required when --tls-cert is set")
	maxRecvMB := flag.Int("max-recv-mb", 16, "max gRPC incoming message size (MB)")
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

	opts := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(*maxRecvMB * 1024 * 1024),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     15 * time.Minute,
			MaxConnectionAge:      30 * time.Minute,
			MaxConnectionAgeGrace: 30 * time.Second,
			Time:                  5 * time.Minute,
			Timeout:               1 * time.Minute,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             30 * time.Second,
			PermitWithoutStream: true,
		}),
	}
	if *tlsCert != "" {
		if *tlsKey == "" {
			log.Fatalf("demangle: --tls-key is required when --tls-cert is set")
		}
		creds, err := credentials.NewServerTLSFromFile(*tlsCert, *tlsKey)
		if err != nil {
			log.Fatalf("demangle: load TLS cert: %v", err)
		}
		opts = append(opts, grpc.Creds(creds))
	}

	srv := grpc.NewServer(opts...)
	svc := newService(demangle.Default, store)
	pb.RegisterDemangleServer(srv, svc)

	// Start health + metrics HTTP server on a separate port.
	_ = startHealthEndpoint(*healthListen, svc.health)

	proto := "plaintext"
	if *tlsCert != "" {
		proto = "TLS"
	}
	fmt.Fprintf(os.Stderr, "demanglegrpc: gRPC listening on %s (%s)\n", lis.Addr(), proto)
	fmt.Fprintf(os.Stderr, "demanglegrpc: health+metrics on http://%s/healthz + /metrics\n", *healthListen)
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("demangle: serve: %v", err)
	}
}
