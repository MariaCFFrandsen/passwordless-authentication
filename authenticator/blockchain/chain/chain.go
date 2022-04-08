package chain

import (
	".authenticator/blockchain/block"
	".authenticator/blockchain/db_access"
	".authenticator/encryption"
	".authenticator/utils"
	"fmt"
	"github.com/dgraph-io/badger"
)

const (
	dbPath = "./tmp/blocks" //TODO: this should be an option

	dbFile = "./tmp/blocks/MANIFEST" //rm
)

type Service interface {
	//InsertGenesis() []byte
	//Insert(b *block.Block) error
	//LastHash() error
	//NextHash(currentHash []byte) *block.Block
	Close()
}

type Blockchain struct {
	LastHash  []byte
	Database  *badger.DB
	dbService *db_access.Access
}

func InitBlockChain(dbFilePath ...string) *Blockchain { //return err if more than 1 arg?
	var lastHash []byte
	var connect *db_access.Access
	var db *badger.DB
	if len(dbFilePath) > 0 { //switch
		connect, db = db_access.Connect(dbFilePath[0])
	} else {
		connect, db = db_access.Connect(dbPath)
	}

	err := db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			fmt.Println("No existing blockchain found")
			genBlock := genesis()
			fmt.Println("Genesis proved")
			err = txn.Set(genBlock.Hash, genBlock.Serialize())
			utils.Handle(err)
			err = txn.Set([]byte("lh"), genBlock.Hash)

			lastHash = genBlock.Hash

			return err
		} else {
			item, err := txn.Get([]byte("lh"))
			utils.Handle(err)
			err = item.Value(func(val []byte) error {
				lastHash = val
				return nil
			})
			utils.Handle(err)
			return err
		}
	})
	utils.Handle(err)

	blockchain := Blockchain{lastHash, db, connect}
	return &blockchain
}

func (chain *Blockchain) AddBlock(data string, pk *encryption.PublicKey) (*block.Block, error) {
	var lastHash []byte

	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		utils.Handle(err)
		err = item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})
		utils.Handle(err)
		return err
	})
	utils.Handle(err)

	newBlock := block.CreateBlock(data, lastHash, pk)

	err = chain.Database.Update(func(transaction *badger.Txn) error { //rename transaction
		err := transaction.Set(newBlock.Hash, newBlock.Serialize())
		utils.Handle(err)
		err = transaction.Set([]byte("lh"), newBlock.Hash)

		chain.LastHash = newBlock.Hash
		return err
	})
	utils.Handle(err)
	return newBlock, err
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
