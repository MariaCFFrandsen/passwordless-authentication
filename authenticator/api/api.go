package api

import ".authenticator/blockchain"

type Service interface { //this is the API the server should have
	AddBlock() (*blockchain.Block, error)
	ValidateBlock() error
	SearchChainByHash(*blockchain.Block, error)
}
