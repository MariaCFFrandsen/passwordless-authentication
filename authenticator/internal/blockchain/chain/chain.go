package chain

import (
	".authenticator/cryptography"
	".authenticator/internal/blockchain/block"
	".authenticator/internal/blockchain/database"
)

const (
	dbPath = "./tmp/blocks" //TODO: this should be an option

	dbFile = "./tmp/blocks/MANIFEST" //rm
)

type dbService interface {
	Insert(newBlock *block.Block) []byte
	LastHash() []byte
	GetBlockByHash(currentHash []byte) *block.Block
	CreateOrFindGenesis(genesis *block.Block) []byte
	Close()
}

type Blockchain struct {
	LastHash  []byte
	dbService dbService
}

func InitBlockChain(dbFilePath ...string) *Blockchain { //return err if more than 1 arg?
	var db *database.Access
	if len(dbFilePath) > 0 { //switch
		db = database.InitDb(dbFilePath[0])
	} else {
		db = database.InitDb(dbPath)
	}
	lastHash := db.CreateOrFindGenesis(Genesis())
	blockchain := Blockchain{lastHash, db}
	return &blockchain
}

func (chain *Blockchain) AddBlock(data string, pk *cryptography.PublicKey) (*block.Block, error) {
	lastHash := chain.dbService.LastHash()
	newBlock := block.CreateBlock(data, lastHash, pk)
	chain.LastHash = chain.dbService.Insert(newBlock)
	return newBlock, nil
}

func Genesis() *block.Block { //should return err
	pair := cryptography.GenerateKeyPair()
	cryptography.SaveCertificate(cryptography.CreateCertificate(*pair), "genesis")
	return block.CreateBlock("Genesis", []byte{}, &cryptography.PublicKey{
		PublicKey: pair.PublicKey.PublicKey,
	})
}
