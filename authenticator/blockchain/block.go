package blockchain

import (
	".authenticator/encryption"
	".authenticator/utils"
	"bytes"
	"encoding/gob"
)

type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
	Nonce    int
	//pk       encryption.PublicKey //could this a pointer?
}

func CreateBlock(data string, prevHash []byte, pk *encryption.PublicKey) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash, 0}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
} //private

func genesis() *Block { //should return err
	pair := encryption.GenerateKeyPair()
	err := encryption.CreateCertificate(pair)
	if err != nil {
		return nil
	}
	return CreateBlock("Genesis", []byte{}, &encryption.PublicKey{
		PublicKey: pair.PublicKey.PublicKey,
	})
}

func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(b)

	utils.Handle(err)

	return res.Bytes()
}

func Deserialize(data []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)

	utils.Handle(err)

	return &block
}
