package main

import (
	"context"
	"fmt"
	"github.com/ahhcash/gengar/internal/types"
	pb "github.com/ahhcash/gengar/proto/generated/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type DocumentServer struct {
	pb.UnimplementedDocumentServiceServer
	store *DocumentStore
}

func (s *DocumentServer) UploadDocument(_ context.Context, req *pb.UploadRequest) (*pb.UploadResponse, error) {
	doc := &types.Document{
		Name:    req.Name,
		Content: req.Content,
	}

	docId, err := s.store.Store(doc, req.ClientPublicKey)

	if err != nil {
		return nil, fmt.Errorf("error storing document: %v", err)
	}

	return &pb.UploadResponse{
		DocumentId: docId.String(),
	}, nil
}

func (s *DocumentServer) DownloadDocument(_ context.Context, req *pb.DownloadRequest) (*pb.DownloadResponse, error) {
	docId, err := uuid.Parse(req.DocumentId)
	if err != nil {
		return nil, fmt.Errorf("error parsig document ID %s, %v", req.DocumentId, err)
	}

	doc, err := s.store.Get(docId)
	if err != nil {
		return nil, err
	}

	metadata := doc.GetMetadata()

	return &pb.DownloadResponse{
		Metadata: &pb.DocumentMetadata{
			Id:        metadata.Id,
			Name:      metadata.Name,
			CreatedAt: metadata.CreatedAt,
			UpdatedAt: metadata.UpdatedAt,
		},
		EncryptedContent: doc.Content,
		Ciphertext:       doc.Ciphertext,
	}, nil
}

func (s *DocumentServer) ListDocuments(context.Context, *pb.Empty) (*pb.ListResponse, error) {
	metadata := s.store.List()

	pbMetadata := make([]*pb.DocumentMetadata, 0, len(metadata))

	for _, m := range metadata {
		pbMetadata = append(pbMetadata, &pb.DocumentMetadata{
			Id:        m.Id,
			Name:      m.Name,
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
		})
	}

	return &pb.ListResponse{
		Documents: pbMetadata,
	}, nil
}

func Run(port int) error {
	documnetStore, err := NewDocumentStore()
	if err != nil {
		return fmt.Errorf("failed to initiaslize document stpre")
	}
	defer documnetStore.Clean()

	list, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("could not start port %d: %v", port, err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterDocumentServiceServer(grpcServer, &DocumentServer{
		store: documnetStore,
	})

	go func() {
		// capture any interrupts and other signals
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		grpcServer.GracefulStop()
	}()

	if err := grpcServer.Serve(list); err != nil {
		return fmt.Errorf("failed to serve on port %d: %v", port, err)
	}
	return nil
}
