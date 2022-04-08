package integration_tests

import (
	".authenticator/blockchain"
	".authenticator/encryption"
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

const (
	dbFile = "..\\..\\tmp\\blocks"
)

func TestBlock(t *testing.T) {
	var (
		bc      = blockchain.InitBlockChain(dbFile)
		keyPair = encryption.GenerateKeyPair()
		rn      = rand.Intn(100)
		block = blockchain.CreateBlock(fmt.Sprintf("test block %d", rn),
			bc.LastHash,
			keyPair.PublicKey)
	)
	t.Run("Create Block", func(t *testing.T) {
		createdBlock := blockchain.CreateBlock(fmt.Sprintf("test block %d", rn),
			bc.LastHash,
			keyPair.PublicKey)
		assert.NotNil(t, createdBlock)
	})

	t.Run("Block to []byte", func(t *testing.T) {
		pb := []byte(fmt.Sprintf("%v", keyPair.PublicKey))
		assert.NotNil(t, block)
		assert.NotEmpty(t, pb, "")
	})

	t.Run("Run PoW", func(t *testing.T) {
		proofOfWork := blockchain.NewProofOfWork(block)
		proofOfWork.Run()
	})

	t.Run("Verify hash", func(t *testing.T) {
		var (
			proofOfWork = blockchain.NewProofOfWork(block)
		)
		proofOfWork.Validate()

	})
}
