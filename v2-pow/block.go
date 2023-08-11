package main

import (
	"fmt"
	"time"
)

type Block struct {
	Version    uint64
	PrevHash   []byte
	MerkleRoot []byte
	TimeStamp  uint64
	Bits       uint64 // complex level
	Nonce      uint64
	Hash       []byte // add it for simplify, BTC don't have this field
	Data       []byte
}

func NewBlock(data string, prevHash []byte) *Block {
	b := Block{
		Version:    0,
		PrevHash:   prevHash,
		MerkleRoot: nil, // TODO
		TimeStamp:  uint64(time.Now().Unix()),
		Nonce:      0,
		Hash:       nil,
		Data:       []byte(data),
	}
	pow := NewProofOfWork(&b)
	pow.Run()
	return &b
}

// except `Hash []byte` field
// func (b *Block) SetHash() {
// 	s := [][]byte{
// 		UintToByte(b.Version),
// 		b.PrevHash,
// 		b.MerkleRoot,
// 		UintToByte(b.TimeStamp),
// 		UintToByte(b.Bits),
// 		UintToByte(b.Nonce),
// 		b.Data,
// 	}
// 	b2 := bytes.Join(s, nil)
// 	b3 := sha256.Sum256(b2)
// 	b.Hash = b3[:]
// }

func (b Block) String() string {
	format := `Version : %d
PrevHash    : %x
MerkleRoot  : %x
TimeStamp   : %d
Bits        : %d
Nonce       : %d
Hash        : %x
Data        : %s
`
	return fmt.Sprintf(format, b.Version, b.PrevHash, b.MerkleRoot,
		b.TimeStamp, b.Bits, b.Nonce, b.Hash, b.Data)
}
