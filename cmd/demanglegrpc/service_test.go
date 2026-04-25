// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package main

import (
	"context"
	"io"
	"net"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/jelmerdehen/demangle"
	pb "github.com/jelmerdehen/demangle/cmd/demanglegrpc/proto/demanglepb"
	_ "github.com/jelmerdehen/demangle/scheme/all"
)

// bring up an in-process server + client for end-to-end tests.
func startServer(t *testing.T) (pb.DemangleClient, func()) {
	t.Helper()
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	store := demangle.InMemoryContextStore()

	srv := grpc.NewServer()
	pb.RegisterDemangleServer(srv, newService(demangle.Default, store))

	done := make(chan struct{})
	go func() {
		srv.Serve(lis) //nolint:errcheck
		close(done)
	}()

	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	cleanup := func() {
		conn.Close()
		srv.GracefulStop()
		<-done
		store.Close()
	}
	return pb.NewDemangleClient(conn), cleanup
}

func TestWrapErrNonStructured(t *testing.T) {
	t.Parallel()
	err := io.ErrUnexpectedEOF
	pe := wrapErr(err)
	if pe == nil {
		t.Fatal("nil Error")
	}
	if pe.Kind != int32(demangle.ErrInternal) {
		t.Fatalf("Kind = %d want ErrInternal", pe.Kind)
	}
	if pe.Message == "" {
		t.Fatal("empty Message")
	}
}

func TestWrapErrStructured(t *testing.T) {
	t.Parallel()
	err := &demangle.Error{
		Kind:     demangle.ErrGrammarViolation,
		Scheme:   "cpp-itanium",
		Offset:   7,
		Expected: "identifier",
		Got:      "'Z'",
		Window:   "_ZN4llvm",
	}
	pe := wrapErr(err)
	if pe.Kind != int32(demangle.ErrGrammarViolation) {
		t.Fatalf("Kind = %d", pe.Kind)
	}
	if pe.Scheme != "cpp-itanium" || pe.Offset != 7 || pe.Expected != "identifier" {
		t.Fatalf("wrong wrapping: %+v", pe)
	}
}

func TestErrorRespShape(t *testing.T) {
	t.Parallel()
	resp := errorResp(42, "rust", "boom")
	if resp.Id != 42 {
		t.Fatalf("Id = %d", resp.Id)
	}
	if resp.Scheme != "rust" {
		t.Fatalf("Scheme = %q", resp.Scheme)
	}
	if resp.Err == nil || resp.Err.Message != "boom" {
		t.Fatalf("Err = %+v", resp.Err)
	}
	if resp.Err.Kind != int32(demangle.ErrInternal) {
		t.Fatalf("Kind = %d want ErrInternal", resp.Err.Kind)
	}
}

func TestGRPCDemangle(t *testing.T) {
	t.Parallel()
	client, cleanup := startServer(t)
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Itanium via auto-detect.
	r, err := client.Demangle(ctx, &pb.Request{Input: "_ZN4llvm5Value4dumpEv"})
	if err != nil {
		t.Fatalf("Demangle: %v", err)
	}
	if r.GetErr() != nil {
		t.Fatalf("Demangle err: %+v", r.GetErr())
	}
	if r.GetOutput() != "llvm::Value::dump()" {
		t.Fatalf("output = %q", r.GetOutput())
	}
	if r.GetScheme() != "cpp-itanium" {
		t.Fatalf("scheme = %q", r.GetScheme())
	}

	// Swift stable function entity.
	r, err = client.Demangle(ctx, &pb.Request{Input: "$s4main3fooyyF"})
	if err != nil {
		t.Fatalf("Demangle swift: %v", err)
	}
	if r.GetOutput() != "main.foo() -> ()" {
		t.Fatalf("swift output = %q", r.GetOutput())
	}

	// JNI.
	r, err = client.Demangle(ctx, &pb.Request{Input: "Java_com_example_Foo_bar"})
	if err != nil {
		t.Fatalf("Demangle jni: %v", err)
	}
	if r.GetOutput() != "com.example.Foo.bar" {
		t.Fatalf("jni output = %q", r.GetOutput())
	}
}

func TestGRPCSchemes(t *testing.T) {
	t.Parallel()
	client, cleanup := startServer(t)
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.Schemes(ctx, &pb.Empty{})
	if err != nil {
		t.Fatalf("Schemes: %v", err)
	}
	if len(resp.GetSchemes()) < 8 {
		t.Fatalf("expected ≥8 schemes, got %d", len(resp.GetSchemes()))
	}
	// Verify at least one scheme has ImplementsMangler=true.
	sawMangler := false
	for _, s := range resp.GetSchemes() {
		if s.GetImplementsMangler() {
			sawMangler = true
			break
		}
	}
	if !sawMangler {
		t.Fatalf("no scheme advertised Mangler; expected jni / kotlin / scala2 / dex / proguard to satisfy it")
	}
}

func TestGRPCDetect(t *testing.T) {
	t.Parallel()
	client, cleanup := startServer(t)
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.Detect(ctx, &pb.DetectRequest{Input: "_ZN4llvm5Value4dumpEv"})
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}
	if len(resp.GetCandidates()) == 0 {
		t.Fatalf("no candidates")
	}
	if resp.GetCandidates()[0].GetScheme() != "cpp-itanium" {
		t.Fatalf("top candidate = %q", resp.GetCandidates()[0].GetScheme())
	}
}

func TestGRPCUploadAndResolveProGuard(t *testing.T) {
	t.Parallel()
	client, cleanup := startServer(t)
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mapBlob := []byte(`com.example.Foo -> a:
    void bar(int) -> b
`)
	up, err := client.UploadContext(ctx, &pb.UploadContextRequest{
		Kind: "proguard_map",
		Blob: mapBlob,
	})
	if err != nil {
		t.Fatalf("UploadContext: %v", err)
	}
	if up.GetSha256() == "" {
		t.Fatalf("empty sha")
	}

	// Demangle with --context-sha.
	r, err := client.Demangle(ctx, &pb.Request{
		Input:  "a.b",
		Scheme: "proguard-map",
		Options: &pb.Options{
			ContextSha256: up.GetSha256(),
		},
	})
	if err != nil {
		t.Fatalf("Demangle: %v", err)
	}
	if r.GetErr() != nil {
		t.Fatalf("err: %+v", r.GetErr())
	}
	if r.GetOutput() != "com.example.Foo.bar" {
		t.Fatalf("output = %q", r.GetOutput())
	}

	// List + delete.
	ll, err := client.ListContexts(ctx, &pb.ListContextsRequest{Kind: "proguard_map"})
	if err != nil || len(ll.GetContexts()) != 1 {
		t.Fatalf("list: %v %+v", err, ll)
	}
	if _, err := client.DeleteContext(ctx, &pb.DeleteContextRequest{Sha256: up.GetSha256()}); err != nil {
		t.Fatalf("delete: %v", err)
	}
}

func TestGRPCMangle(t *testing.T) {
	t.Parallel()
	client, cleanup := startServer(t)
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// JNI round-trip: pass the already-mangled symbol as input; the server
	// demangling-parses it to a Node tree and re-mangles it back.
	// "Java_com_example_Foo_bar" → demangle → tree → mangle → "Java_com_example_Foo_bar"
	r, err := client.Mangle(ctx, &pb.MangleRequest{
		Scheme: "jni",
		Input:  "Java_com_example_Foo_bar",
	})
	if err != nil {
		t.Fatalf("Mangle jni: %v", err)
	}
	if r.GetError() != "" {
		t.Fatalf("Mangle jni error: %q", r.GetError())
	}
	if r.GetMangled() != "Java_com_example_Foo_bar" {
		t.Fatalf("jni mangled = %q, want %q", r.GetMangled(), "Java_com_example_Foo_bar")
	}

	// JVM descriptor round-trip: "Ljava/lang/String;" → tree → mangle → same.
	r2, err := client.Mangle(ctx, &pb.MangleRequest{
		Scheme: "jvmdesc",
		Input:  "Ljava/lang/String;",
	})
	if err != nil {
		t.Fatalf("Mangle jvmdesc: %v", err)
	}
	if r2.GetError() != "" {
		t.Fatalf("Mangle jvmdesc error: %q", r2.GetError())
	}
	if r2.GetMangled() != "Ljava/lang/String;" {
		t.Fatalf("jvmdesc mangled = %q, want %q", r2.GetMangled(), "Ljava/lang/String;")
	}

	// Unknown scheme returns a MangleResponse with a non-empty Error field,
	// not a gRPC-level error.
	r3, err := client.Mangle(ctx, &pb.MangleRequest{
		Scheme: "no-such-scheme",
		Input:  "whatever",
	})
	if err != nil {
		t.Fatalf("Mangle unknown scheme: unexpected gRPC error %v", err)
	}
	if r3.GetError() == "" {
		t.Fatal("Mangle unknown scheme: expected Error in response")
	}

	// Empty scheme returns an error in the response.
	r4, err := client.Mangle(ctx, &pb.MangleRequest{Input: "whatever"})
	if err != nil {
		t.Fatalf("Mangle empty scheme: unexpected gRPC error %v", err)
	}
	if r4.GetError() == "" {
		t.Fatal("Mangle empty scheme: expected Error in response")
	}
}

func TestGRPCDemangleStream(t *testing.T) {
	t.Parallel()
	client, cleanup := startServer(t)
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream, err := client.DemangleStream(ctx)
	if err != nil {
		t.Fatalf("open stream: %v", err)
	}
	inputs := []string{
		"_ZN4llvm5Value4dumpEv",
		"$s4main3fooyyF",
		"Java_com_example_Foo_bar",
	}
	go func() {
		for i, in := range inputs {
			if err := stream.Send(&pb.Request{Id: uint64(i), Input: in}); err != nil {
				return
			}
		}
		stream.CloseSend() //nolint:errcheck
	}()

	got := map[uint64]string{}
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("recv: %v", err)
		}
		if resp.GetErr() != nil {
			t.Fatalf("stream err: %+v", resp.GetErr())
		}
		got[resp.GetId()] = resp.GetOutput()
	}
	if got[0] != "llvm::Value::dump()" {
		t.Fatalf("0 = %q", got[0])
	}
	if got[1] != "main.foo() -> ()" {
		t.Fatalf("1 = %q", got[1])
	}
	if got[2] != "com.example.Foo.bar" {
		t.Fatalf("2 = %q", got[2])
	}
}
