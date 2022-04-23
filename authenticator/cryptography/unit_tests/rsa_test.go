package unit_tests

import (
	"github.com/MariaCFFrandsen/passwordless-authentication/authenticator/cryptography"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestPublicKeyEncryption(t *testing.T) {
	var (
		keyPair = cryptography.GenerateKeyPair()
		cipher  = "Super confidential message"
	)

	t.Run("Encryption", func(t *testing.T) {
		encrypt, err := cryptography.Encrypt(cipher, keyPair.PublicKey)
		assert.NoErrorf(t, err, "error occurred when encrypting")
		msg, err := cryptography.Decrypt(encrypt, keyPair.PrivateKey)
		assert.NoErrorf(t, err, "error occurred when decrypting")
		assert.Equal(t, cipher, msg)
	})

	t.Run("Acceptable nonce size", func(t *testing.T) {
		hashPublicKey := cryptography.PublicKeyToNonce(keyPair.PublicKey.PublicKey)
		assert.True(t, hashPublicKey < math.MaxInt64)
	})
}
