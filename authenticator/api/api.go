package api

import (
	".authenticator/cryptography"
	".authenticator/internal/blockchain/block"
)

type Service interface { //this is the API the server should have
	AddBlock(data string, key cryptography.PublicKey) (*block.Block, error)
	ValidateBlock(hash []byte) error
	SearchChainByHash(hash []byte) (*block.Block, error)
}
