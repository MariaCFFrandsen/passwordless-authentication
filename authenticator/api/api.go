package api

import (
	".authenticator/blockchain/block"
)

type Service interface { //this is the API the server should have
	AddBlock() (*block.Block, error)
	ValidateBlock() error
	SearchChainByHash(*block.Block, error)
}
