package database

import (
	".authenticator/blockchain/block"
	".authenticator/utils"
	"fmt"
	"github.com/dgraph-io/badger"
)

type Access struct {
	Db *badger.DB
}

func InitDb(dbPath string) *Access { //receiver?
	var opts badger.Options
	opts = badger.DefaultOptions(dbPath)
	opts.Truncate = true
	db, err := badger.Open(opts)
	utils.Handle(err)
	return &Access{db}
}

func (a *Access) Insert(newBlock *block.Block) []byte {
	var lastHash []byte
	err := a.Db.Update(func(transaction *badger.Txn) error { //rename transaction
		err := transaction.Set(newBlock.Hash, newBlock.Serialize())
		utils.Handle(err)
		err = transaction.Set([]byte("lh"), newBlock.Hash)

		lastHash = newBlock.Hash
		return err
	})
	utils.Handle(err)
	return lastHash
}

func (a *Access) LastHash() []byte {
	var lastHash []byte
	err := a.Db.View(func(txn *badger.Txn) error {
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
	return lastHash
}

func (a *Access) GetBlockByHash(currentHash []byte) *block.Block {
	var b *block.Block
	err := a.Db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(currentHash)
		utils.Handle(err)

		err = item.Value(func(val []byte) error {
			b = block.Deserialize(val)
			return nil
		})
		utils.Handle(err)
		return err
	})
	utils.Handle(err)
	return b
}

func (a *Access) CreateOrFindGenesis(genesis *block.Block) []byte {
	var lastHash []byte
	err := a.Db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			fmt.Println("No existing blockchain found")
			genBlock := genesis
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
	return lastHash
}

func (a *Access) Close() {
	a.Db.Close()
}
