package integration_tests

import (
	blockchain2 ".authenticator/blockchain"
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBlockchain(t *testing.T) {
	//TODO: make some test with database

	t.Run("Init Blockchain", func(t *testing.T) {
		blockchain := blockchain2.InitBlockChain()
		assert.NotNil(t, blockchain)
	})

	t.Run("Correct last hash", func(t *testing.T) {
		var (
			blockchain = blockchain2.InitBlockChain()
		)
		h := blockchain.LastHash
		fmt.Printf(string(h))
	})

	t.Run("Add Block", func(t *testing.T) {
		var (
			blockchain = blockchain2.InitBlockChain()
		)
		blockchain.AddBlock("test block", nil)
	})

	t.Run("", func(t *testing.T) {

	})

	t.Run("Find Block by Hash", func(t *testing.T) {
		var (
			blockchain = blockchain2.InitBlockChain()
		)
		iterator := blockchain.Iterator()
		findHash := []byte{}
		for iterator.Next() != nil {
			if bytes.Compare(findHash, iterator.CurrentHash) == 0 {
				//hash found
				//this is a shitty and slow method but now we can one
			}
		}
	})
}
