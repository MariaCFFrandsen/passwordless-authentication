package blockchain

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

// Difficulty /*People shouldn't be able to create Blocks willy nilly */
/*The idea is that lots of work has to be done to create a new block on the chain.
But there is another piece; we need to be able to verify that the work has actually been done.
PoW in a nutshell says that performing the algorithm is hard "work",
but proving that the "work" was done should be easy.
*/
const Difficulty = 12

/*
	1. Grab the data from the block
	2. Create a nonce
	3. Create a hash of the data from the block + the nonce
	4. Check the proceeding hash to see if it meets the requirements we listed above. (the first few bytes must be 0's).
*/
type ProofOfWork struct {
	Block *Block
	Target *big.Int
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-Difficulty))
	/* We're shifting the first arg target,
	however many units left we set our 2nd arg.
	The closer we get to 256, the easier the computation will be.
	Increasing our difficulty will increase the runtime of our algorithm.*/

	pow := &ProofOfWork{b, target}

	return pow
}

// InitNonce - Nonce is very smart.
func (pow *ProofOfWork) InitNonce(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.Block.PrevHash,
			pow.Block.Data,
			ToHex(int64(nonce)),
			ToHex(int64(Difficulty)),
		},
		[]byte{},
	)
	return data
}

func (pow *ProofOfWork) Validate() bool {
	var intHash big.Int

	data := pow.InitNonce(pow.Block.Nonce)

	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	return intHash.Cmp(pow.Target) == -1
}

func (pow ProofOfWork) Run() (int, []byte) {
	var intHash big.Int
	var hash [32]byte

	nonce := 0
	// This is essentially an infinite loop due to how large
	// MaxInt64 is.
	for nonce < math.MaxInt64 {
		data := pow.InitNonce(nonce)
		hash = sha256.Sum256(data)

		fmt.Printf("\r%x", hash)
		intHash.SetBytes(hash[:])

		if intHash.Cmp(pow.Target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Println()

	return nonce, hash[:]
}
