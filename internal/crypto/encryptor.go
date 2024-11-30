package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"github.com/ahhcash/gengar/internal/types"
	"io"
)

var kyber *KyberWrapper

func init() {
	// hardcoding this to kyber1024, we may want to customize this later
	kyber, _ = NewKyberWrapper("Kyber1024")
}

func Encrypt(doc *types.Document, publicKey []byte) (*types.Document, []byte, []byte, error) {
	ciphertext, sharedSecret, err := kyber.EncapsulateSecret(publicKey)
	defer kyber.Clean()

	if err != nil {
		return nil, nil, nil, err
	}

	// kyber is just a key encapsulation mechanism. it must be used in conjunction with AES.
	// the keys and the data can now be exchanged securely

	block, err := aes.NewCipher(sharedSecret)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error creating aes cipher: %v", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error creating GCM: %v", err)
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, nil, fmt.Errorf("error creating nonce: %v", err)
	}

	encryptedDoc := gcm.Seal(nonce, nonce, doc.Content, nil)
	doc.Content = encryptedDoc
	doc.PrivateKey = ciphertext

	return doc, ciphertext, sharedSecret, nil
}

func Decrypt(doc *types.Document, ciphertext []byte) ([]byte, error) {
	sharedSecret, err := kyber.DecapsulateSecret(ciphertext)
	defer kyber.Clean()

}
