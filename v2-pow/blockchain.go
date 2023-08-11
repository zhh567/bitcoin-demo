package main

import (
	"fmt"
	"strings"
)

const genesisInfo = "Init blockchain demo"

type BlockChain struct {
	Blocks []*Block
}

func NewBlockChain() *BlockChain {
	b := [32]byte{0}
	init := NewBlock(genesisInfo, b[:])
	return &BlockChain{
		[]*Block{init},
	}
}

func (bc *BlockChain) AddBlock(data string) {
	preHash := bc.Blocks[len(bc.Blocks)-1].Hash
	curBlock := NewBlock(data, preHash)
	bc.Blocks = append(bc.Blocks, curBlock)
}

func (bc BlockChain) String() string {
	sb := strings.Builder{}
	for i, block := range bc.Blocks {
		sb.WriteString(fmt.Sprintf("\n++++++++++++++++++ Current block height: %d ++++++++++++++++++\n", i))
		sb.WriteString(block.String())
	}
	return sb.String()
}
