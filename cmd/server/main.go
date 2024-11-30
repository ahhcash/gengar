package server

import (
	"context"
	"fmt"
	"github.com/ahhcash/gengar/internal/types"
	pb "github.com/ahhcash/gengar/proto/generated/proto"
	"github.com/google/uuid"
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

	pbMetadata := make([]*pb.DocumentMetadata, len(metadata))

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
