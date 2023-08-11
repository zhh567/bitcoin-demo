package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

type ProofOfWork struct {
	block  *Block
	target *big.Int // system provided
}

func NewProofOfWork(block *Block) *ProofOfWork {
	targetInt := new(big.Int)
	targetInt.SetString("0010000000000000000000000000000000000000000000000000000000000000", 16)

	return &ProofOfWork{
		block:  block,
		target: targetInt,
	}
}

func (pow *ProofOfWork) Run() uint64 {
	fmt.Printf("Finding nounce: \n")

	var nounce uint64 = 0
	for ; true; nounce++ {
		pow.block.Nonce = nounce
		hashBytes, isValid := pow.IsValid()

		fmt.Printf("\r%x", hashBytes)
		if isValid {
			fmt.Print("\n")
			pow.block.Hash = hashBytes
			return nounce
		}
	}
	panic("Fail to find a nounce for create a block")
}

// if block's SHA256 less than boundary?
func (pow *ProofOfWork) IsValid() ([]byte, bool) {
	data := pow.prepareData(pow.block.Nonce)

	hashBytes := sha256.Sum256(data)
	hashInt := new(big.Int)
	hashInt.SetBytes(hashBytes[:])

	return hashBytes[:], (hashInt.Cmp(pow.target) == -1)
}

// convert block struct to byte slice
func (pow *ProofOfWork) prepareData(nonce uint64) []byte {
	b := pow.block
	s := [][]byte{
		UintToByte(b.Version),
		b.PrevHash,
		b.MerkleRoot,
		UintToByte(b.TimeStamp),
		UintToByte(b.Bits),
		UintToByte(nonce),
		b.Data,
	}
	return bytes.Join(s, []byte{})
}
