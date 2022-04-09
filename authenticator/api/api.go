package api

import (
	".authenticator/cryptography"
	".authenticator/internal/blockchain/block"
)

type Service interface { //this is the API the server should have
	CreateUser(data string, key cryptography.PublicKey) (*block.Block, error)
	AuthenticateUser(hash []byte) bool
}
