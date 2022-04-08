package chain

import (
	".authenticator/internal/blockchain/block"
)

type Iterator struct {
	CurrentHash []byte
	Database    dbService
}

func (chain *Blockchain) Iterator() *Iterator {
	return &Iterator{chain.LastHash, chain.dbService}
}

func (iterator *Iterator) Next() *block.Block {
	b := iterator.Database.GetBlockByHash(iterator.CurrentHash)
	iterator.CurrentHash = b.PrevHash
	return b
}

