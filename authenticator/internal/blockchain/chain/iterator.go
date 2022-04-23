package chain

import (
	"authenticator/cryptography"
	"authenticator/internal/blockchain/block"
	crypto "crypto/x509"
	"fmt"
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

func (iterator *Iterator) SearchBlockchainByPublicKey(pk *cryptography.PublicKey) (*block.Block, bool) { //tmp
	continueLoop := true
	found := false
	var foundBlock *block.Block
	for continueLoop {
		b := iterator.Database.GetBlockByHash(iterator.CurrentHash)
		fmt.Printf("hash: %x\n", b.Hash)
		key, _ := crypto.ParsePKCS1PublicKey(b.PublicKey)
		if found = pk.PublicKey.Equal(key); found {
			foundBlock = b
			continueLoop = false
		}
		iterator.CurrentHash = b.PrevHash
	}
	return foundBlock, found
}
