package chain

import (
	".authenticator/blockchain/block"
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

func (iterator *Iterator) Next() *block.Block {
	var b *block.Block

	err := iterator.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iterator.CurrentHash)
		utils.Handle(err)

		err = item.Value(func(val []byte) error {
			b = block.Deserialize(val)
			return nil
		})
		utils.Handle(err)
		return err
	})
	utils.Handle(err)

	iterator.CurrentHash = b.PrevHash

	return b
}

