package blockchain

import (
	".authenticator/encryption"
	"fmt"
	"github.com/dgraph-io/badger"
)

const (
	dbPath = "./tmp/blocks" //TODO: this should be an option

	dbFile = "./tmp/blocks/MANIFEST" //rm
)

type APIService interface { //this is the API the server should have
	AddBlock() (*Block, error)
	ValidateBlock() error
	SearchChainByHash(*Block, error)
}

type Blockchain struct {
	LastHash []byte
	Database *badger.DB
}

func InitBlockChain(dbFilePath ...string) *Blockchain { //return err if more than 1 arg?
	var lastHash []byte
	var opts badger.Options
	if len(dbFilePath) > 0 {
		opts = badger.DefaultOptions(dbFilePath[0])
	} else {
		opts = badger.DefaultOptions(dbPath)
	}

	opts.Truncate = true
	db, err := badger.Open(opts)
	Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			fmt.Println("No existing blockchain found")
			genesis := Genesis()
			fmt.Println("Genesis proved")
			err = txn.Set(genesis.Hash, genesis.Serialize())
			Handle(err)
			err = txn.Set([]byte("lh"), genesis.Hash)

			lastHash = genesis.Hash

			return err
		} else {
			item, err := txn.Get([]byte("lh"))
			Handle(err)
			err = item.Value(func(val []byte) error {
				lastHash = val
				return nil
			})
			Handle(err)
			return err
		}
	})
	Handle(err)

	blockchain := Blockchain{lastHash, db}
	return &blockchain
}

func (chain *Blockchain) AddBlock(data string, pk *encryption.PublicKey) (*Block, error){
	var lastHash []byte

	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		err = item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})
		Handle(err)
		return err
	})
	Handle(err)

	newBlock := CreateBlock(data, lastHash, pk)

	err = chain.Database.Update(func(transaction *badger.Txn) error { //rename transaction
		err := transaction.Set(newBlock.Hash, newBlock.Serialize())
		Handle(err)
		err = transaction.Set([]byte("lh"), newBlock.Hash)

		chain.LastHash = newBlock.Hash
		return err
	})
	Handle(err)
	return newBlock, err
}