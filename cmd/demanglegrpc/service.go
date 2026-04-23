// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 Jelmer de Hen

package main

import (
	"context"
	"errors"
	"io"

	"github.com/jelmerdehen/demangle"
	pb "github.com/jelmerdehen/demangle/cmd/demanglegrpc/proto/demanglepb"
)

// service is the gRPC adapter over the library. It's a thin shim —
// every call translates the proto message into a library call and
// translates the result back. No business logic here; all behaviour
// lives in the library so GraphQL (on skynet) and gRPC stay in sync.
type service struct {
	pb.UnimplementedDemangleServer
	cat   *demangle.Catalog
	store demangle.ContextStore
}

func newService(cat *demangle.Catalog, store demangle.ContextStore) *service {
	return &service{cat: cat, store: store}
}

func (s *service) Demangle(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	opts, err := s.buildOptions(ctx, req.GetOptions())
	if err != nil {
		return nil, err
	}
	var (
		r   *demangle.Result
		dErr error
	)
	if sn := req.GetScheme(); sn != "" {
		sch, ok := s.cat.Scheme(sn)
		if !ok {
			return errorResp(req.GetId(), sn, "unknown scheme"), nil
		}
		r, dErr = sch.Demangle(ctx, req.GetInput(), *opts)
	} else {
		r, dErr = s.cat.Demangle(ctx, req.GetInput(), opts)
	}
	return wrapResult(req.GetId(), req.GetInput(), r, dErr), nil
}

func (s *service) Detect(ctx context.Context, req *pb.DetectRequest) (*pb.DetectResponse, error) {
	cands := s.cat.Detect(req.GetInput(), demangle.DetectOptions{
		MaxCandidates:    int(req.GetMaxCandidates()),
		MinConfidence:    int(req.GetMinConfidence()),
		IncludeWeak:      req.GetIncludeWeak(),
		SchemeHintFamily: req.GetSchemeHintFamily(),
		Strict:           req.GetStrict(),
		AmbiguityWindow:  int(req.GetAmbiguityWindow()),
	})
	out := &pb.DetectResponse{}
	for _, c := range cands {
		out.Candidates = append(out.Candidates, &pb.Candidate{
			Scheme:     c.Scheme,
			Confidence: int32(c.Confidence),
			Signals:    c.Signals,
			Negatives:  c.Negatives,
			Diagnostic: c.Diagnostic,
		})
	}
	return out, nil
}

func (s *service) Schemes(_ context.Context, _ *pb.Empty) (*pb.SchemesResponse, error) {
	infos := s.cat.Schemes()
	out := &pb.SchemesResponse{}
	for _, info := range infos {
		sch, _ := s.cat.Scheme(info.Name)
		_, isMangler := sch.(demangle.Mangler)
		caps := sch.Capabilities()
		out.Schemes = append(out.Schemes, &pb.SchemeInfo{
			Name:              info.Name,
			Family:            info.Family,
			Version:           info.Version,
			Description:       info.Description,
			Stability:         info.Stability.String(),
			MangleFidelity:    info.MangleFidelity.String(),
			ImplementsMangler: isMangler,
			MaxInputBytes:     int32(caps.MaxInputBytes),
			RequiresContext:   info.RequiresContext,
		})
	}
	return out, nil
}

func (s *service) DemangleStream(stream pb.Demangle_DemangleStreamServer) error {
	// Producer: receives from the wire, feeds the library's batch API.
	in := make(chan demangle.BatchRequest, 256)
	out := make(chan demangle.BatchResponse, 256)

	errCh := make(chan error, 1)
	go func() {
		var optsCopy *demangle.Options
		for {
			req, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				close(in)
				return
			}
			if err != nil {
				errCh <- err
				close(in)
				return
			}
			if optsCopy == nil {
				o, oerr := s.buildOptions(stream.Context(), req.GetOptions())
				if oerr != nil {
					errCh <- oerr
					close(in)
					return
				}
				optsCopy = o
			}
			_ = optsCopy // options applied per-batch in Stage 6.5
			in <- demangle.BatchRequest{ID: req.GetId(), Input: req.GetInput()}
		}
	}()

	done := make(chan demangle.BatchSummary, 1)
	go func() {
		done <- s.cat.DemangleBatch(stream.Context(), in, out, demangle.BatchOptions{})
	}()

	for resp := range out {
		msg := wrapBatchResp(resp)
		if err := stream.Send(msg); err != nil {
			return err
		}
	}
	select {
	case e := <-errCh:
		return e
	default:
	}
	<-done
	return nil
}

func (s *service) UploadContext(ctx context.Context, req *pb.UploadContextRequest) (*pb.UploadContextResponse, error) {
	if s.store == nil {
		return nil, errors.New("demangle: no context store configured")
	}
	sha, err := s.store.Put(ctx, req.GetKind(), req.GetBlob(), req.GetMetadata())
	if err != nil {
		return nil, err
	}
	return &pb.UploadContextResponse{Sha256: sha, ByteSize: int64(len(req.GetBlob()))}, nil
}

func (s *service) ListContexts(ctx context.Context, req *pb.ListContextsRequest) (*pb.ListContextsResponse, error) {
	if s.store == nil {
		return &pb.ListContextsResponse{}, nil
	}
	infos, err := s.store.List(ctx, req.GetKind())
	if err != nil {
		return nil, err
	}
	out := &pb.ListContextsResponse{}
	for _, i := range infos {
		out.Contexts = append(out.Contexts, &pb.ContextInfo{
			Kind:          i.Kind,
			Sha256:        i.SHA256,
			ByteSize:      i.ByteSize,
			UploadedTs:    i.UploadedTS.Unix(),
			LastAccessTs:  i.LastAccessTS.Unix(),
			Metadata:      i.Metadata,
		})
	}
	return out, nil
}

func (s *service) DeleteContext(ctx context.Context, req *pb.DeleteContextRequest) (*pb.Empty, error) {
	if s.store != nil {
		if err := s.store.Delete(ctx, req.GetSha256()); err != nil {
			return nil, err
		}
	}
	return &pb.Empty{}, nil
}

// --- helpers -----------------------------------------------------

func (s *service) buildOptions(ctx context.Context, proto *pb.Options) (*demangle.Options, error) {
	if proto == nil {
		return &demangle.Options{}, nil
	}
	o := &demangle.Options{
		Simplified:                    proto.GetSimplified(),
		SynthesizeSugar:               proto.GetSynthesizeSugar(),
		QualifyEntities:               proto.GetQualifyEntities(),
		DisplayGenericSpecialisations: proto.GetDisplayGenericSpecialisations(),
		DisplayThunks:                 proto.GetDisplayThunks(),
		ReturnTree:                    proto.GetReturnTree(),
		VerifyRoundTrip:               proto.GetVerifyRoundTrip(),
		AllowLegacy:                   proto.GetAllowLegacy(),
		BestEffortMangle:              proto.GetBestEffortMangle(),
		SchemeSpecific:                proto.GetSchemeSpecific(),
	}
	if sha := proto.GetContextSha256(); sha != "" && s.store != nil {
		c, err := s.store.Get(ctx, sha)
		if err != nil {
			return nil, err
		}
		o.Context = c
	}
	return o, nil
}

func wrapResult(id uint64, input string, r *demangle.Result, dErr error) *pb.Response {
	resp := &pb.Response{Id: id, Input: input}
	if dErr != nil {
		resp.Err = wrapErr(dErr)
		return resp
	}
	if r != nil {
		resp.Scheme = r.Scheme
		resp.Output = r.Output
		resp.Confidence = int32(r.Confidence)
		resp.Partial = r.Partial
		resp.LostInfo = r.LostInfo
		resp.Tree = wrapNode(r.Tree)
		for k, v := range r.Annotations {
			resp.Annotations = append(resp.Annotations, &pb.Annotation{Key: k, Value: v})
		}
	}
	return resp
}

func wrapBatchResp(r demangle.BatchResponse) *pb.Response {
	out := &pb.Response{Id: r.ID, Input: r.Input, Scheme: r.Scheme}
	if r.Err != nil {
		out.Err = wrapErr(r.Err)
		return out
	}
	if r.Result != nil {
		out.Output = r.Result.Output
		out.Confidence = int32(r.Result.Confidence)
		out.Partial = r.Result.Partial
		out.LostInfo = r.Result.LostInfo
		out.Tree = wrapNode(r.Result.Tree)
		for k, v := range r.Result.Annotations {
			out.Annotations = append(out.Annotations, &pb.Annotation{Key: k, Value: v})
		}
	}
	return out
}

func wrapErr(err error) *pb.Error {
	var e *demangle.Error
	if errors.As(err, &e) {
		return &pb.Error{
			Kind: int32(e.Kind), Scheme: e.Scheme, Offset: int32(e.Offset),
			Expected: e.Expected, Got: e.Got, Window: e.Window,
			Message: e.Error(),
		}
	}
	return &pb.Error{Kind: int32(demangle.ErrInternal), Message: err.Error()}
}

func wrapNode(n *demangle.Node) *pb.Node {
	if n == nil {
		return nil
	}
	out := &pb.Node{Scheme: n.Scheme, Kind: n.Kind, Text: n.Text, Index: n.Index}
	for k, v := range n.Attrs {
		out.Attrs = append(out.Attrs, &pb.Annotation{Key: k, Value: v})
	}
	for _, c := range n.Children {
		out.Children = append(out.Children, wrapNode(c))
	}
	return out
}

func errorResp(id uint64, scheme, message string) *pb.Response {
	return &pb.Response{
		Id:     id,
		Scheme: scheme,
		Err: &pb.Error{
			Kind:    int32(demangle.ErrInternal),
			Scheme:  scheme,
			Message: message,
		},
	}
}
