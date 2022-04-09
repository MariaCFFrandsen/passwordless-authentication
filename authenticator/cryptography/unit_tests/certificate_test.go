package unit_tests

import (
	".authenticator/cryptography"
	crypto "crypto/x509"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCertificate(t *testing.T) {
	t.Run("Marshalling/Unmarshalling of certificate", func(t *testing.T) {
		var (
			text = "hello my friend"
		)
		c := cryptography.Certificate{Text: text}
		toBytes := cryptography.ToBytes(c)
		cryptography.FromBytes(toBytes)
		assert.Equal(t, text, c.Text)
	})

	t.Run("Create certificate", func(t *testing.T) {
		var (
			keyPair     = cryptography.GenerateKeyPair()
			certificate = cryptography.InitCertificate(*keyPair)
		)

		cryptography.CreateCertificate(certificate)
	})

	t.Run("Decrypt certificate", func(t *testing.T) {
		c := cryptography.ReadCertificate() //if we experience problems, try to change to txt
		assert.Equal(t, "we try", c.Text)
		assert.Equal(t, cryptography.GetMacAddr(), c.MacAddress)
	})

	t.Run("DUAl test", func(t *testing.T) {
		var (
			keyPair     = cryptography.GenerateKeyPair()
			certificate = cryptography.InitCertificate(*keyPair)
		)
		cryptography.CreateCertificate(certificate)
		c := cryptography.ReadCertificate()
		assert.Equal(t, "we try", c.Text)
		assert.Equal(t, cryptography.GetMacAddr(), c.MacAddress)
		assert.True(t, keyPair.PrivateKey.PrivateKey.Equal(c.PRK))
		assert.True(t, keyPair.PrivateKey.PrivateKey.PublicKey.Equal(c.PUK))

		c2 := cryptography.ReadCertificate()
		assert.Equal(t, "we try", c2.Text)
		assert.Equal(t, cryptography.GetMacAddr(), c2.MacAddress)
		assert.True(t, keyPair.PrivateKey.PrivateKey.Equal(c.PRK))
		assert.True(t, keyPair.PrivateKey.PrivateKey.PublicKey.Equal(c.PUK))
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

	t.Run("length of certificate", func(t *testing.T) {
		var (
			keyPair     = *cryptography.GenerateKeyPair()
			mac         = cryptography.GetMacAddr()
			txt         = "worked!"
			pk          = crypto.MarshalPKCS1PrivateKey(keyPair.PrivateKey.PrivateKey)
			puk         = crypto.MarshalPKCS1PublicKey(keyPair.PublicKey.PublicKey)
			certificate = cryptography.Certificate{
				PrivateKey: pk,
				PublicKey:  puk,
				MacAddress: mac,
				Text:       "worked!",
			}
		)

		bytes := cryptography.ToBytes(certificate)
		fromBytes := cryptography.FromBytes(bytes)
		privateKey, _ := crypto.ParsePKCS1PrivateKey(fromBytes.PrivateKey)
		publickey, _ := crypto.ParsePKCS1PublicKey(fromBytes.PublicKey)
		assert.True(t, keyPair.PrivateKey.PrivateKey.Equal(privateKey))
		assert.True(t, keyPair.PrivateKey.PrivateKey.PublicKey.Equal(publickey))
		assert.Equal(t, mac, fromBytes.MacAddress)
		assert.Equal(t, txt, fromBytes.Text)
		assert.True(t, 4096 > len(bytes))
		var (
			msg = "this is a secret"
		)
		cipher, _ := cryptography.Encrypt(msg, &cryptography.PublicKey{
			PublicKey: publickey,
		})

		decrypt, _ := cryptography.Decrypt(cipher, &cryptography.PrivateKey{
			PrivateKey: privateKey,
		})
		assert.Equal(t, msg, decrypt)
	})
}
