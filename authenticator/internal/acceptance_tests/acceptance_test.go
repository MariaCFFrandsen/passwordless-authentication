package acceptance_tests

import (
	".authenticator/cryptography"
	".authenticator/internal/blockchain/block"
	".authenticator/internal/blockchain/chain"
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
	t.Run("Print added block", func(t *testing.T) {
		var (
			bc      = chain.InitBlockChain(dbPath)
			keyPair = cryptography.GenerateKeyPair()
			rn      = rand.Intn(100)
		)

		defer bc.Iterator().Database.Close()

		b, err := bc.AddBlock(fmt.Sprintf("test %d", rn), keyPair.PublicKey)
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
			fmt.Println()
			if len(b.PrevHash) == 0 {
				break
			}
		}
	})

	t.Run("acceptance test, happy path", func(t *testing.T) {
		var (
			bc      = chain.InitBlockChain(dbPath)
			keyPair = cryptography.GenerateKeyPair()
			rn      = rand.Intn(100)
			mac     = cryptography.GetMacAddr()
			keyPath = "genesis-key.txt"
			msg     = "this is a secret message"
		)
		defer bc.Iterator().Database.Close()                                 //should close on iterator? prob not
		b, err := bc.AddBlock(fmt.Sprintf("test %d", rn), keyPair.PublicKey) //genesis instead, or should there be one for genesis and one for additional block?
		assert.NoErrorf(t, err, "error occurred creating b")                 //---this should be method
		fmt.Println("Added Block!")
		fmt.Printf("Previous hash: %x\n", b.PrevHash)
		fmt.Printf("data: %s\n", b.Data)
		fmt.Printf("hash: %x\n", b.Hash)
		pow := block.NewProofOfWork(b)
		fmt.Printf("Pow: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Printf("Nonce: %d\n", b.Nonce)                //this--this should be method
		certificate := cryptography.RetrieveCertificate() //"genesis", remember to change save certificate
		assert.True(t, keyPair.PrivateKey.PrivateKey.Equal(certificate.PRK))
		assert.True(t, keyPair.PrivateKey.PrivateKey.PublicKey.Equal(certificate.PUK))
		assert.Equal(t, mac, certificate.MacAddress)
		assert.Equal(t, keyPath, certificate.Text)
		cipher, _ := cryptography.Encrypt(msg, &cryptography.PublicKey{
			PublicKey: certificate.PUK,
		})
		decrypt, _ := cryptography.Decrypt(cipher, &cryptography.PrivateKey{
			PrivateKey: certificate.PRK,
		})
		assert.Equal(t, msg, decrypt)
	})

	t.Run("AddBlock with certificate", func(t *testing.T) {

	})
}
