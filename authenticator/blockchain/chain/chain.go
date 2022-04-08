package chain

import (
	".authenticator/blockchain/block"
	".authenticator/blockchain/dbaccess"
	".authenticator/encryption"
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
	var db *dbaccess.Access
	if len(dbFilePath) > 0 { //switch
		db = dbaccess.Connect(dbFilePath[0])
	} else {
		db = dbaccess.Connect(dbPath)
	}
	lastHash := db.CreateOrFindGenesis(genesis())
	blockchain := Blockchain{lastHash, db}
	return &blockchain
}

func (chain *Blockchain) AddBlock(data string, pk *encryption.PublicKey) (*block.Block, error) {
	lastHash := chain.dbService.LastHash()
	newBlock := block.CreateBlock(data, lastHash, pk)
	chain.LastHash = chain.dbService.Insert(newBlock)
	return newBlock, nil
}

func genesis() *block.Block { //should return err
	pair := encryption.GenerateKeyPair()
	err := encryption.CreateCertificate(pair)
	if err != nil {
		return nil
	}
	return block.CreateBlock("Genesis", []byte{}, &encryption.PublicKey{
		PublicKey: pair.PublicKey.PublicKey,
	})
}
