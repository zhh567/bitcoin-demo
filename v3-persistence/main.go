package main

import (
	"fmt"
)

func main() {
	// get blockchain
	bc, err := GetBlockChain()
	if err != nil {
		fmt.Println("can't get blockchain: ", err)
	}
	defer bc.Close()
	// add block
	// bc.AddBlock("111")
	// bc.AddBlock("222")
	// bc.AddBlock("333")
	bc.AddBlock("444")

	// print
	iter := bc.NewIterator()
	for {
		block := iter.Next()
		if block == nil {
			break
		}
		fmt.Println("\n*************************************************************")
		fmt.Println(block.String())
	}
}
