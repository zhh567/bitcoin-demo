package main

import (
	"flag"
	"fmt"
	"strconv"
)

type Cli struct {
	PrintNum          int
	AddressGetBalance string
	SendCoin          bool
}

func NewCli() *Cli {
	cli := &Cli{}
	flag.IntVar(&cli.PrintNum, "print", 0, "print a specified number of blocks (0 < number < 20)")
	flag.StringVar(&cli.AddressGetBalance, "getbalance", "", "get balance of an address: -getbalance <address>")
	flag.BoolVar(&cli.SendCoin, "send", false, "send to someone: -send <from-address> <to-address> <amount>")
	flag.Parse()
	return cli
}

func (cli *Cli) Run() {
	bc, err := GetBlockChain()
	if err != nil {
		fmt.Println("can't get blockchain: ", err)
	}
	defer bc.Close()

	if cli.PrintNum > 0 {
		cli.Print(bc)
	} else if cli.AddressGetBalance != "" {
		cli.GetBalance(bc, cli.AddressGetBalance)
	} else if cli.SendCoin {
		if len(flag.Args()) != 3 {
			fmt.Println("invalid command")
			return
		}
		from := flag.Arg(0)
		to := flag.Arg(1)
		amount, err := strconv.Atoi(flag.Arg(2))
		if err != nil {
			fmt.Println("the amount must be a number")
			return
		}
		cli.Send(from, to, int64(amount), "Genesis", "some data written to sig by miner")
	} else {
		fmt.Println("invalid command")
	}
}

func (cli *Cli) Print(bc *BlockChain) {
	iter := bc.NewIterator()
	for i := 0; i < cli.PrintNum && i < 20; i++ {
		block := iter.Next()
		if block == nil {
			break
		}
		fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>[Block]>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
		fmt.Println(block.String())
		// fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
	}
}

func (cli *Cli) GetBalance(bc *BlockChain, address string) {
	_, total := bc.FindUtxo(address)
	fmt.Printf("[%s] remain utxos: %d\n", address, total)
}

func (cli *Cli) Send(from, to string, amount int64, minerPubKey string, data string) {
	bc, err := GetBlockChain()
	if err != nil {
		fmt.Println("Can't get blockchain: ", err)
	}
	defer bc.Close()

	miningTx := NewMiningTx(minerPubKey, data)
	tx := NewTransaction(from, to, amount, bc)
	if tx == nil {
		fmt.Printf("Transfer [%d] from [%s] to [%s] failed.\n", amount, from, to)
		return
	}

	bc.AddBlock([]*Transaction{miningTx, tx})
	fmt.Printf("Transfer [%d] from [%s] to [%s] success.\n", amount, from, to)
}
