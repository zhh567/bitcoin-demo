package main

import (
	"flag"
	"fmt"
)

type Cli struct {
	PrintNum     int
	AddBlockData string
}

func NewCli() *Cli {
	cli := &Cli{}
	flag.IntVar(&cli.PrintNum, "print", 0, "print a specified number of blocks (0 < number < 20)")
	flag.StringVar(&cli.AddBlockData, "add", "", "add some data to blockchain")
	flag.Parse()
	return cli
}

func (cli *Cli) Run() {
	bc, err := GetBlockChain()
	if err != nil {
		fmt.Println("can't get blockchain: ", err)
	}
	defer bc.Close()

	if cli.AddBlockData != "" {
		cli.AddBlock(bc)
	}

	if cli.PrintNum > 0 {
		cli.Print(bc)
	}
}

func (cli *Cli) Print(bc *BlockChain) {
	iter := bc.NewIterator()
	for i := 0; i < cli.PrintNum && i < 20; i++ {
		block := iter.Next()
		if block == nil {
			break
		}
		fmt.Println("\n*************************************************************")
		fmt.Println(block.String())
	}
}

func (cli *Cli) AddBlock(bc *BlockChain) {
	err := bc.AddBlock(cli.AddBlockData)
	if err != nil {
		fmt.Println("add block fail: ", err)
		return
	}
}
