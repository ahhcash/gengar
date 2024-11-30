package crypto

import (
	"fmt"
	"github.com/open-quantum-safe/liboqs-go/oqs"
)

type KyberWrapper struct {
	// OQS supports three kyber variants - Kyber512, 768 adn 1024, this tracks the specific variant chosen
	kemType string
	kem     oqs.KeyEncapsulation
}

type KeyPair struct {
	publicKey  []byte
	privateKey []byte
}

func NewKyberWrapper(kemType string) (*KyberWrapper, error) {
	kyberWrapper := &KyberWrapper{
		kemType: kemType,
	}

	if err := kyberWrapper.kem.Init(kemType, nil); err != nil {
		return nil, err
	}

	return kyberWrapper, nil
}

func (k *KyberWrapper) Clean() {
	k.kem.Clean()
}

func (k *KyberWrapper) Details() string {
	return k.kem.Details().String()
}

func (k *KyberWrapper) GenerateKeyPair() (*KeyPair, error) {
	publicKey, err := k.kem.GenerateKeyPair()
	if err != nil {
		return nil, fmt.Errorf("error when generating key pair: %v", err)
	}

	privateKey := k.kem.ExportSecretKey()
	return &KeyPair{
		publicKey:  publicKey,
		privateKey: privateKey,
	}, nil
}

func (k *KyberWrapper) EncapsulateSecret(publicKey []byte) ([]byte, []byte, error) {
	ciphertext, sharedSecret, err := k.kem.EncapSecret(publicKey)
	if err != nil {
		return nil, nil, fmt.Errorf("error encapsulating secret: %v", err)
	}

	return ciphertext, sharedSecret, nil
}

func (k *KyberWrapper) DecapsulateSecret(ciphertext []byte) ([]byte, error) {
	// the kem object of the wrapper must have it's secret key established

	sharedSecret, err := k.kem.DecapSecret(ciphertext)

	if err != nil {
		return nil, fmt.Errorf("error decapsulating secret: %v", err)
	}

	return sharedSecret, nil
}
