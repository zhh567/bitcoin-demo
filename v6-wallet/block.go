package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"time"
)

type Block struct {
	Version      uint64
	PrevHash     []byte
	MerkleRoot   []byte
	TimeStamp    uint64
	Bits         uint64 // complex level
	Nonce        uint64
	Hash         []byte // add it for simplify, BTC don't have this field
	Transactions []*Transaction
}

func NewBlock(txs []*Transaction, prevHash []byte) *Block {
	b := Block{
		Version:      0,
		PrevHash:     prevHash,
		MerkleRoot:   nil,
		TimeStamp:    uint64(time.Now().Unix()),
		Nonce:        0,
		Hash:         nil,
		Transactions: txs,
	}
	// fill in MerkleRoot
	b.HashTransactionsMerkleRoot()
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
Txs         : 
%v`
	return fmt.Sprintf(format, b.Version, b.PrevHash, b.MerkleRoot,
		b.TimeStamp, b.Bits, b.Nonce, b.Hash, b.Transactions)
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

// not implement true merkle tree, just calculate SHA256 of all transactions
func (b *Block) HashTransactionsMerkleRoot() {
	var txHashes [][]byte
	for _, tx := range b.Transactions {
		txHash := tx.Id
		txHashes = append(txHashes, txHash)
	}
	value := bytes.Join(txHashes, []byte{})
	hash := sha256.Sum256(value)
	b.MerkleRoot = hash[:]
}
