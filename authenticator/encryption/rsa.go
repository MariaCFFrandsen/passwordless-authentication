package encryption

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"fmt"
)

func GenerateKeyPair() *rsa.PrivateKey {
	keyPair, _ := rsa.GenerateKey(rand.Reader, 2048) // 1024 - 4096 supported
	return keyPair
}

func Encrypt(secretMessage string, publicKey rsa.PublicKey) (string, error) {
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, &publicKey, []byte(secretMessage), []byte(""))
	//also takes a label, do not know if we need it
	fmt.Errorf("encryption caused error: %s", err)
	return base64.StdEncoding.EncodeToString(ciphertext), err
}

func Decrypt(cipherText string, privateKey rsa.PrivateKey) string {
	ct, _ := base64.StdEncoding.DecodeString(cipherText)
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, &privateKey, ct, []byte(""))
	fmt.Errorf("decryption caused error: %s", err)
	fmt.Println("Plaintext:", string(plaintext)) //delete this later
	return string(plaintext)
}

func PublicKeyToNonce(publicKey rsa.PublicKey) int {
	//this is super wrong but I don't really understand how the public key modulus/exponent works
	//this is a task to change it
	//we also need to expand block with rsa.PublicKey(probably)

	b := make([]byte, 8)
	n, err := rand.Read(b)
	if n != 8 {
		panic(n)
	} else if err != nil {
		panic(err)
	}
	return int(binary.BigEndian.Uint64(b)) % n
}