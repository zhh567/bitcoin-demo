package main

import (
	"bytes"
	"encoding/gob"
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

func (b Block) String() string {
	format := `Version     : %d
PrevHash    : %x
MerkleRoot  : %x
TimeStamp   : %d
Bits        : %d
Nonce       : %d
Hash        : %x
Data        : %s`
	return fmt.Sprintf(format, b.Version, b.PrevHash, b.MerkleRoot,
		b.TimeStamp, b.Bits, b.Nonce, b.Hash, b.Data)
}

func (b *Block) Serialize() ([]byte, error) {
	var buf = bytes.Buffer{}
	e := gob.NewEncoder(&buf)
	err := e.Encode(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Deserialize(data []byte) (*Block, error) {
	var block Block
	d := gob.NewDecoder(bytes.NewReader(data))
	err := d.Decode(&block)
	if err != nil {
		return nil, err
	}
	return &block, nil
}
