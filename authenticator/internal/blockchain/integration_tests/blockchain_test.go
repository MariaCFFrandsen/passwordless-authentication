package integration_tests

import (
	crypto "crypto/x509"
	"fmt"
	"github.com/passwordless-authentication/authenticator/cryptography"
	blockchain "github.com/passwordless-authentication/authenticator/internal/blockchain/chain"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	dbPath = "..\\..\\..\\tmp\\blocks"
)

func TestBlockchain(t *testing.T) {
	t.Run("Init Blockchain", func(t *testing.T) {
		bc := blockchain.InitBlockChain()
		assert.NotNil(t, bc)
	})

	t.Run("Correct last hash", func(t *testing.T) {
		var (
			bc = blockchain.InitBlockChain()
		)
		h := bc.LastHash
		fmt.Printf(string(h))
	})

	t.Run("Add Block", func(t *testing.T) {
		var (
			bc = blockchain.InitBlockChain()
		)
		bc.AddBlock("test block", nil)
	})

	t.Run("Find Block by Hash", func(t *testing.T) {
		var (
		//bc = blockchain.InitBlockChain()
		)

	})

	t.Run("Find Block by Public Key", func(t *testing.T) {
		var (
			iterator = blockchain.InitBlockChain(dbPath).Iterator()
			block    = iterator.Next()
		)
		iterator.Database.Close()
		iterator = blockchain.InitBlockChain(dbPath).Iterator()
		key, _ := crypto.ParsePKCS1PublicKey(block.PublicKey)
		foundBlock, found := iterator.SearchBlockchainByPublicKey(&cryptography.PublicKey{PublicKey: key})
		assert.True(t, found)
		assert.NotNil(t, foundBlock)
		assert.Equal(t, block.Hash, foundBlock.Hash)
		assert.Equal(t, block.PrevHash, foundBlock.PrevHash)
		assert.Equal(t, block.Nonce, foundBlock.Nonce)
		assert.Equal(t, block.Data, foundBlock.Data)
	})
}
