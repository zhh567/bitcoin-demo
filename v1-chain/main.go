package main

import (
	"fmt"
	"time"
)

func main() {
	bc := NewBlockChain()
	time.Sleep(1 * time.Second)

	bc.AddBlock("this is second blcok")
	time.Sleep(1 * time.Second)

	bc.AddBlock("this is third block")
	time.Sleep(1 * time.Second)

	fmt.Println(bc)
}
