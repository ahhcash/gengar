package types

import (
	"github.com/google/uuid"
	"time"
)

type Document struct {
	Id         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Content    []byte    `json:"content"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	Ciphertext []byte    `json:"ciphertext"`
}

type DocumentMetadata struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func (d *Document) GetMetadata() *DocumentMetadata {
	return &DocumentMetadata{
		Id:        d.Id.String(),
		Name:      d.Name,
		CreatedAt: d.CreatedAt.String(),
		UpdatedAt: d.UpdatedAt.String(),
	}
}
