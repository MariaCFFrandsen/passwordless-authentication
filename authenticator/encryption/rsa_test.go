package encryption

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestPublicKeyEncryption(t *testing.T) {
	var(
		keyPair = GenerateKeyPair()
		cipher = "Super confidential message"
	)

	t.Run("Encryption", func(t *testing.T) {
		encrypt, err := Encrypt(cipher, keyPair.PublicKey)
		assert.NoErrorf(t, err, "error occurred when encrypting")
		msg, err := Decrypt(encrypt, keyPair.PrivateKey)
		assert.NoErrorf(t, err, "error occurred when decrypting")
		assert.Equal(t, cipher, msg)
	})

	t.Run("Acceptable nonce size", func(t *testing.T) {
		hashPublicKey := PublicKeyToNonce(keyPair.PublicKey)
		assert.True(t, hashPublicKey < math.MaxInt64)
	})
}
