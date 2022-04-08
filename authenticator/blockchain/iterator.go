package blockchain

import (
	".authenticator/utils"
	"github.com/dgraph-io/badger"
)

type Iterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func (chain *Blockchain) Iterator() *Iterator {
	return &Iterator{chain.LastHash, chain.Database}
}

func (iterator *Iterator) Next() *Block {
	var block *Block

	err := iterator.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iterator.CurrentHash)
		utils.Handle(err)

		err = item.Value(func(val []byte) error {
			block = Deserialize(val)
			return nil
		})
		utils.Handle(err)
		return err
	})
	utils.Handle(err)

	iterator.CurrentHash = block.PrevHash

	return block
}

