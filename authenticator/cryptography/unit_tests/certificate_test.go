package unit_tests

import (
	".authenticator/cryptography"
	"bytes"
	"encoding/gob"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	key = []byte("example key 1234")
)

func TestCertificate(t *testing.T) {
	t.Run("Marshalling/Unmarshalling of certificate", func(t *testing.T) {
		var (
			text = "hello my friend"
		)
		c := cryptography.Certificate{Text: text}
		buffer := &bytes.Buffer{}
		gob.NewEncoder(buffer).Encode(c)
		byteSlice := buffer.Bytes()

		var unmarshal cryptography.Certificate
		bf := bytes.NewBuffer(byteSlice)
		gob.NewDecoder(bf).Decode(&unmarshal)
		assert.Equal(t, text, c.Text)
	})

	t.Run("Marshalling/Unmarshalling of string", func(t *testing.T) {
		var (
			text = "hello my friend"
		)
		var s string
		buffer := &bytes.Buffer{}
		gob.NewEncoder(buffer).Encode(text)
		byteSlice := buffer.Bytes()

		bf := bytes.NewBuffer(byteSlice)
		gob.NewDecoder(bf).Decode(&s)
		assert.Equal(t, text, s)
	})

	t.Run("Unmarshalling a string", func(t *testing.T) {
		var (
			text = "hello my friend"
		)
		byteSlice := []byte(text)
		assert.Equal(t, 15, len(byteSlice))
	})

	t.Run("Create certificate", func(t *testing.T) {
		var (
			text = "hello my friend"
		)
		cryptography.CreateCertificate(text)
	})

	t.Run("Decrypt certificate", func(t *testing.T) {
		var (
			text = "hello my friend"
		)
		plaintext := cryptography.ReadCertificate()
		assert.Equal(t, text, plaintext)
	})

	t.Run("DUAl test", func(t *testing.T) {
		var (
			text = "hello my friend"
		)
		cryptography.CreateCertificate(text)
		plaintext := cryptography.ReadCertificate()
		assert.Equal(t, text, plaintext)
	})

	t.Run("Symmetric Key Generation", func(t *testing.T) {
		sk := cryptography.CreateSymmetricKey()
		assert.NotEmpty(t, sk)
		assert.Equal(t, 16, len(sk))
	})

	t.Run("write and read key file", func(t *testing.T) {
		symmetricKey := cryptography.CreateSymmetricKey()
		readKey := cryptography.ReadKeyFile()
		assert.Equal(t, symmetricKey, readKey)
	})

}
