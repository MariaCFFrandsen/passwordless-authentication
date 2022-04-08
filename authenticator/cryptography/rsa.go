package cryptography

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

type KeyPair struct {
	PrivateKey *PrivateKey
	PublicKey  *PublicKey
}

type PublicKey struct {
	PublicKey *rsa.PublicKey
}

type PrivateKey struct {
	PrivateKey *rsa.PrivateKey
}

func GenerateKeyPair() *KeyPair {
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048) // 1024 - 4096 supported
	return &KeyPair{
		PrivateKey: &PrivateKey{
			privateKey,
		},
		PublicKey:  &PublicKey{
			PublicKey: &privateKey.PublicKey,
		},
	}
}

func Encrypt(secretMessage string, publicKey *PublicKey) (string, error) {
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey.PublicKey, []byte(secretMessage), []byte(""))
	//also takes a label, do not know if we need it
	fmt.Errorf("cryptography caused error: %s", err) //TODO: define logger
	return base64.StdEncoding.EncodeToString(ciphertext), err
}

func Decrypt(cipherText string, privateKey *PrivateKey) (string, error) {
	ct, _ := base64.StdEncoding.DecodeString(cipherText)
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey.PrivateKey, ct, []byte("")) //[]byte("") this is a label, perhaps we can use for create/authenticate
	fmt.Errorf("decryption caused error: %s", err)
	return string(plaintext), err
}

func CreateCertificate(pair *KeyPair) error {
	//not implemented
	return nil
}

func PublicKeyToNonce(publicKey *rsa.PublicKey) int {
	//this is super wrong, but I don't really understand how the public key modulus/exponent works
	//there is a task to change it
	//we also need to expand block with rsa.PublicKey(probably)
	return int(publicKey.N.Uint64())
}
