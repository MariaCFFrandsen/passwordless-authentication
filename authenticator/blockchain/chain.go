package blockchain

import (
	"fmt"
	"github.com/dgraph-io/badger"
)

const (
	dbPath = "./tmp/blocks"

	// This can be used to verify that the blockchain exists
	dbFile = "./tmp/blocks/MANIFEST"

	// This is arbitrary data for our genesis block
	genesisData = "First Transaction from Genesis"
)

type Blockchain struct {
	LastHash []byte
	Database *badger.DB
}

func InitBlockChain() *Blockchain {
	var lastHash []byte

	opts := badger.DefaultOptions(dbPath)

	db, err := badger.Open(opts)
	Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		// "lh" stand for last hash
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

func (chain *Blockchain) AddBlock(data string) {
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

	newBlock := CreateBlock(data, lastHash)

	err = chain.Database.Update(func(transaction *badger.Txn) error {
		err := transaction.Set(newBlock.Hash, newBlock.Serialize())
		Handle(err)
		err = transaction.Set([]byte("lh"), newBlock.Hash)

		chain.LastHash = newBlock.Hash
		return err
	})
	Handle(err)
}

type Iterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func (chain *Blockchain) Iterator() *Iterator {
	iterator := Iterator{chain.LastHash, chain.Database}

	return &iterator
}

func (iterator *Iterator) Next() *Block {
	var block *Block

	err := iterator.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iterator.CurrentHash)
		Handle(err)

		err = item.Value(func(val []byte) error {
			block = Deserialize(val)
			return nil
		})
		Handle(err)
		return err
	})
	Handle(err)

	iterator.CurrentHash = block.PrevHash

	return block
}
