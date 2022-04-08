package api

import (
	".authenticator/blockchain/block"
	".authenticator/encryption"
)

type Service interface { //this is the API the server should have
	AddBlock(data string, key encryption.PublicKey) (*block.Block, error)
	ValidateBlock(hash []byte) error
	SearchChainByHash(hash []byte) (*block.Block, error)
}
