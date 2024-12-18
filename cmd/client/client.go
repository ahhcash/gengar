package main

import (
	"context"
	"fmt"
	"github.com/ahhcash/gengar/internal/crypto"
	"github.com/ahhcash/gengar/internal/types"
	pb "github.com/ahhcash/gengar/proto/generated/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

const (
	kemVariant = "Kyber768"
)

type DocumentClient struct {
	conn      *grpc.ClientConn
	client    pb.DocumentServiceClient
	encryptor *crypto.DocEncryptor
}

func NewDocumentClient(serverAddr string) (*DocumentClient, error) {
	conn, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("error connecting to server: %v", err)
	}

	encryptor, err := crypto.NewDocEncryptor(kemVariant)
	if err != nil {
		return nil, fmt.Errorf("could not initialize encryptor for grpc client: %v", err)
	}

	client := pb.NewDocumentServiceClient(conn)

	// is the server even running lol
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = client.ListDocuments(ctx, &pb.Empty{})
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("could not connect to server: %v", err)
	}

	return &DocumentClient{
		conn:      conn,
		client:    client,
		encryptor: encryptor,
	}, nil
}

func (c *DocumentClient) Close() {
	if c.conn != nil {
		_ = c.conn.Close()
	}
}

func (c *DocumentClient) UploadDocument(name string, content []byte) (string, error) {

	keyPair := c.encryptor.GetKeyPair()

	resp, err := c.client.UploadDocument(context.Background(), &pb.UploadRequest{
		Name:            name,
		Content:         content,
		ClientPublicKey: keyPair.PublicKey,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload document: %v", err)
	}

	return resp.DocumentId, nil
}

func (c *DocumentClient) DownloadDocument(documentID string) (*types.Document, error) {
	resp, err := c.client.DownloadDocument(context.Background(), &pb.DownloadRequest{
		DocumentId: documentID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to download document: %v", err)
	}

	createdAt, _ := time.Parse(time.DateTime, resp.Metadata.CreatedAt)
	updatedAt, _ := time.Parse(time.DateTime, resp.Metadata.UpdatedAt)

	doc := &types.Document{
		Content:    resp.EncryptedContent,
		Ciphertext: resp.Ciphertext,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}

	// Decrypt the document
	if err := c.encryptor.Decrypt(doc); err != nil {
		return nil, fmt.Errorf("failed to decrypt document: %v", err)
	}

	return doc, nil
}

func (c *DocumentClient) ListDocuments() ([]*types.DocumentMetadata, error) {
	resp, err := c.client.ListDocuments(context.Background(), &pb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("failed to list documents: %v", err)
	}

	metadata := make([]*types.DocumentMetadata, len(resp.Documents))
	for i, doc := range resp.Documents {
		metadata[i] = &types.DocumentMetadata{
			Id:        doc.Id,
			Name:      doc.Name,
			CreatedAt: doc.CreatedAt,
			UpdatedAt: doc.UpdatedAt,
		}
	}

	return metadata, nil
}

func (c *DocumentClient) ViewDocument(documentID string) (*types.Document, error) {
	resp, err := c.client.ViewDocument(context.Background(), &pb.ViewRequest{
		DocumentId: documentID,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to download document: %v", err)
	}

	createdAt, _ := time.Parse(time.DateTime, resp.Metadata.CreatedAt)
	updatedAt, _ := time.Parse(time.DateTime, resp.Metadata.UpdatedAt)

	parsedId, _ := uuid.Parse(resp.Metadata.Id)
	doc := &types.Document{
		Id:        parsedId,
		Content:   resp.Contents,
		Name:      resp.Metadata.Name,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	return doc, nil
}
