package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"github.com/ahhcash/gengar/internal/types"
	"io"
)

type DocEncryptor struct {
	kyber   *KyberWrapper
	keyPair *KeyPair
}

func NewDocEncryptor(kemType string) (*DocEncryptor, error) {
	wrapper, err := NewKyberWrapper(kemType)
	if err != nil {
		return nil, err
	}

	keyPair, err := wrapper.GenerateKeyPair()
	if err != nil {
		wrapper.Clean()
		return nil, err
	}

	return &DocEncryptor{
		wrapper,
		keyPair,
	}, nil
}

func (e *DocEncryptor) Encrypt(doc *types.Document, publicKey []byte) error {
	ciphertext, sharedSecret, err := e.kyber.EncapsulateSecret(publicKey)
	defer e.kyber.Clean()

	if err != nil {
		return err
	}

	// kyber is just a key encapsulation mechanism. it must be used in conjunction with AES.
	// the keys and the data can now be exchanged securely

	block, err := aes.NewCipher(sharedSecret)
	if err != nil {
		return fmt.Errorf("error creating aes cipher: %v", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("error creating GCM: %v", err)
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return fmt.Errorf("error creating nonce: %v", err)
	}

	encryptedDoc := gcm.Seal(nonce, nonce, doc.Content, nil)
	doc.Content = encryptedDoc
	doc.Ciphertext = ciphertext

	return nil
}

func (e *DocEncryptor) Decrypt(doc *types.Document) error {
	sharedSecret, err := e.kyber.DecapsulateSecret(doc.Ciphertext)
	defer e.kyber.Clean()

	if err != nil {
		return fmt.Errorf("error decryption secret: %v", err)
	}

	block, err := aes.NewCipher(sharedSecret)
	if err != nil {
		return fmt.Errorf("error creating AES cippher: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("error creating GCM: %v", err)
	}

	nonceSize := gcm.NonceSize()
	if nonceSize > len(doc.Content) {
		return fmt.Errorf("encrypted content is corrupted")
	}

	nonce, contentCipherText := doc.Content[:nonceSize], doc.Content[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, contentCipherText, nil)
	if err != nil {
		return fmt.Errorf("failed to decrypt document: %v", err)
	}

	doc.Content = plaintext

	return nil
}
