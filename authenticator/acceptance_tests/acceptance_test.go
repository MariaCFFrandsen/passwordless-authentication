package acceptance_tests

import (
	".authenticator/blockchain"
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
			bc      = blockchain.InitBlockChain(dbPath)
			keyPair = encryption.GenerateKeyPair()
			rn      = rand.Intn(100)
		)

		defer bc.Database.Close()

		block, err := bc.AddBlock(fmt.Sprintf("test %d", rn), keyPair.PublicKey)
		assert.NoErrorf(t, err, "error occurred creating block")
		fmt.Println("Added Block!")
		fmt.Printf("Previous hash: %x\n", block.PrevHash)
		fmt.Printf("data: %s\n", block.Data)
		fmt.Printf("hash: %x\n", block.Hash)
		pow := blockchain.NewProofOfWork(block)
		fmt.Printf("Pow: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Printf("Nonce: %d\n", block.Nonce)
	})

	t.Run("Print blockchain", func(t *testing.T) {
		var (
			bc       = blockchain.InitBlockChain(dbPath)
			iterator = bc.Iterator()
		)
		defer bc.Database.Close()

		for {
			block := iterator.Next()
			fmt.Printf("Previous hash: %x\n", block.PrevHash)
			fmt.Printf("data: %s\n", block.Data)
			fmt.Printf("hash: %x\n", block.Hash)
			pow := blockchain.NewProofOfWork(block)
			fmt.Printf("Pow: %s\n", strconv.FormatBool(pow.Validate()))
			fmt.Printf("Nonce: %d\n", block.Nonce)
			fmt.Println()
			if len(block.PrevHash) == 0 {
				break
			}
		}
	})
}
