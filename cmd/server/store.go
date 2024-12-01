package main

import (
	"fmt"
	"github.com/ahhcash/gengar/internal/crypto"
	"github.com/ahhcash/gengar/internal/types"
	"github.com/google/uuid"
	"sync"
	"time"
)

const (
	kemVariant = "Kyber768"
)

type DocumentStore struct {
	// let's keep it thread safe
	documents sync.Map

	encryptor *crypto.DocEncryptor
}

func NewDocumentStore() (*DocumentStore, error) {
	encryptor, err := crypto.NewDocEncryptor(kemVariant)
	if err != nil {
		return nil, fmt.Errorf("error initializing encryptor: %v", err)
	}

	return &DocumentStore{
		documents: sync.Map{},
		encryptor: encryptor,
	}, nil
}

func (ds *DocumentStore) Store(doc *types.Document, clientPublicKey []byte) (uuid.UUID, error) {
	docId := uuid.New()
	doc.Id = docId

	now := time.Now()
	doc.CreatedAt, doc.UpdatedAt = now, now

	if err := ds.encryptor.Encrypt(doc, clientPublicKey); err != nil {
		return uuid.Nil, err
	}

	ds.documents.Store(docId, doc)
	return docId, nil
}

func (ds *DocumentStore) Get(id uuid.UUID) (*types.Document, error) {
	val, exists := ds.documents.Load(id)
	if !exists {
		return nil, fmt.Errorf("%s is not a valid document id", id.String())
	}

	doc, ok := val.(*types.Document)
	if !ok {
		return nil, fmt.Errorf("invalid document format for id: %s", id.String())
	}

	return doc, nil
}

func (ds *DocumentStore) List() []*types.DocumentMetadata {
	var metadata []*types.DocumentMetadata
	ds.documents.Range(func(key, val interface{}) bool {
		if doc, ok := val.(*types.Document); ok {
			metadata = append(metadata, doc.GetMetadata())
		}

		return true
	})

	return metadata
}
