package block

import (
	"authenticator/cryptography"
	"authenticator/internal/utils"
	"bytes"
	crypto "crypto/x509"
	"encoding/gob"
)

type Block struct {
	Hash      []byte
	Data      []byte
	PrevHash  []byte
	Nonce     int
	PublicKey []byte
}

func CreateBlock(data string, prevHash []byte, pk *cryptography.PublicKey) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash, 0, crypto.MarshalPKCS1PublicKey(pk.PublicKey)}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
} //private

func (b *Block) Serialize() []byte { //figure out what do with the publickey
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)
	err := encoder.Encode(b)
	utils.Handle(err)
	return res.Bytes()
}

func Deserialize(data []byte) *Block { //figure out what do with the publickey
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	utils.Handle(err)
	return &block
}
