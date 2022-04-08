package acceptance_tests

import (
	".authenticator/blockchain/block"
	".authenticator/blockchain/chain"
	".authenticator/encryption"
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"strconv"
	"testing"
)

const (
	dbPath = "..\\tmp\\blocks"
)

func TestPrint(t *testing.T) {
	t.Run("Print added block", func(t *testing.T) {
		var (
			bc      = chain.InitBlockChain(dbPath)
			keyPair = encryption.GenerateKeyPair()
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
}
