package acceptance_tests

import (
	"authenticator/cryptography"
	"authenticator/internal/blockchain/block"
	"authenticator/internal/blockchain/chain"
	crypto "crypto/x509"
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"strconv"
	"testing"
)

const (
	dbPath = "..\\..\\tmp\\blocks"
)

func TestPrint(t *testing.T) {
	t.Run("Add block without certificate", func(t *testing.T) { //make this work
		var (
			bc      = chain.InitBlockChain(dbPath)
			keyPair = cryptography.GenerateKeyPair()
			rn      = rand.Intn(100)
		)

		defer bc.Iterator().Database.Close()

		b, err := bc.AddBlock(fmt.Sprintf("test-%d", rn), keyPair.PublicKey)
		assert.NoErrorf(t, err, "error occurred creating b")
		fmt.Println("Added Block!")
		fmt.Printf("Previous hash: %x\n", b.PrevHash)
		fmt.Printf("data: %s\n", b.Data)
		fmt.Printf("hash: %x\n", b.Hash)
		pow := block.NewProofOfWork(b)
		fmt.Printf("Pow: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Printf("Nonce: %d\n", b.Nonce)
	})

	t.Run("Print blockchain", func(t *testing.T) {
		var (
			bc       = chain.InitBlockChain(dbPath)
			iterator = bc.Iterator()
		)
		defer bc.Iterator().Database.Close()

		for {
			b := iterator.Next()
			fmt.Printf("Previous hash: %x\n", b.PrevHash)
			fmt.Printf("data: %s\n", b.Data)
			fmt.Printf("hash: %x\n", b.Hash)
			pow := block.NewProofOfWork(b)
			fmt.Printf("Pow: %s\n", strconv.FormatBool(pow.Validate()))
			fmt.Printf("Nonce: %d\n", b.Nonce)
			publicKey, _ := crypto.ParsePKCS1PublicKey(b.PublicKey)
			fmt.Printf("Public key: %d\n", publicKey.E)
			fmt.Printf("Public key 2: %d\n", publicKey.N)
			fmt.Println()
			if len(b.PrevHash) == 0 {
				break
			}
		}
	})

	t.Run("acceptance test, happy path", func(t *testing.T) { //make this work
		var (
			bc      = chain.InitBlockChain(dbPath)
			keyPair = cryptography.GenerateKeyPair()
			kp      = cryptography.KeyPair{
				PrivateKey: keyPair.PrivateKey,
				PublicKey:  keyPair.PublicKey,
			}
			rn  = rand.Intn(100)
			mac = cryptography.GetMacAddr()
			msg = "this is a secret message"
		)
		defer bc.Iterator().Database.Close() //should close on iterator? prob not
		b, err := bc.AddBlock(fmt.Sprintf("test-%d", rn), keyPair.PublicKey)
		assert.NoErrorf(t, err, "error occurred creating b") //---this should be method
		fmt.Println("Added Block!")
		fmt.Printf("Previous hash: %x\n", b.PrevHash)
		fmt.Printf("data: %s\n", b.Data)
		fmt.Printf("hash: %x\n", b.Hash)
		pow := block.NewProofOfWork(b)
		fmt.Printf("Pow: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Printf("Nonce: %d\n", b.Nonce) //this--this should be method
		publicKey, err := crypto.ParsePKCS1PublicKey(b.PublicKey)
		fmt.Printf("Public key: %d\n", publicKey.E)
		cryptography.SaveCertificate(cryptography.CreateCertificate(kp), fmt.Sprintf("test-%d", rn))
		certificate := cryptography.RetrieveCertificate(fmt.Sprintf("test-%d", rn))
		assert.True(t, keyPair.PrivateKey.PrivateKey.Equal(certificate.PRK))
		assert.True(t, keyPair.PrivateKey.PrivateKey.PublicKey.Equal(certificate.PUK))
		assert.Equal(t, mac, certificate.MacAddress)
		cipher, _ := cryptography.Encrypt(msg, &cryptography.PublicKey{
			PublicKey: certificate.PUK,
		})
		decrypt, _ := cryptography.Decrypt(cipher, &cryptography.PrivateKey{
			PrivateKey: certificate.PRK,
		})
		assert.Equal(t, msg, decrypt)
	})

	t.Run("Genesis with certificate", func(t *testing.T) {
		//database must be empty
		bc := chain.InitBlockChain(dbPath)
		defer bc.Iterator().Database.Close()
	})

}
