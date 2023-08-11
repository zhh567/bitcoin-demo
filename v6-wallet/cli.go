package main

import (
	"flag"
	"fmt"
	"sort"
	"strconv"
)

type Cli struct {
	Create            bool
	PrintNum          int
	AddressGetBalance string
	SendCoin          bool

	CreateWallet     bool
	ListAllAddresses bool
}

func NewCli() *Cli {
	cli := &Cli{}
	flag.BoolVar(&cli.Create, "create", false, "create a new blockchain: -create <miner-address>")
	flag.IntVar(&cli.PrintNum, "print", 0, "print a specified number of blocks (0 < number < 20): -print <number>")
	flag.StringVar(&cli.AddressGetBalance, "getbalance", "", "get balance of an address: -getbalance <address>")
	flag.BoolVar(&cli.SendCoin, "send", false, "send to someone: -send <from-address> <to-address> <amount> <miner-address> <data>")
	flag.BoolVar(&cli.CreateWallet, "createwallet", false, "create a new wallet")
	flag.BoolVar(&cli.ListAllAddresses, "listAllAddresses", false, "list all addresses (and private key) in wallet")
	flag.Parse()
	return cli
}

func (cli *Cli) Run() {
	if cli.CreateWallet {
		wm := NewWalletManager()
		address := wm.CreateWallet()
		fmt.Printf("New wallet created: %s\n", address)
		return
	}
	if cli.ListAllAddresses {
		wm := NewWalletManager()
		addresses := wm.ListAllAddresses()
		fmt.Printf("All addresses in wallet:\n")
		sort.Strings(addresses)
		for _, address := range addresses {
			fmt.Println(address)
		}
		return
	}

	if cli.Create {
		if len(flag.Args()) != 2 {
			fmt.Println("invalid command, command format: -create <miner-address> <genesis-info>")
			return
		}
		if !IsValidAddress(flag.Arg(0)) {
			fmt.Println("invalid address: ", flag.Arg(0))
			return
		}
		err2 := CreateBlockChain(flag.Arg(0), flag.Arg(1))
		if err2 != nil {
			fmt.Println("create blockchain fail: ", err2)
		}
		return
	}

	// following command will operate the blockchain.db file
	bc, err := GetBlockChain()
	if err != nil {
		fmt.Println("can't get blockchain: ", err)
	}
	defer bc.Close()

	if cli.PrintNum > 0 {
		cli.Print(bc)
		return
	}
	if cli.AddressGetBalance != "" {
		if !IsValidAddress(cli.AddressGetBalance) {
			fmt.Println("invalid address: ", cli.AddressGetBalance)
			return
		}
		cli.GetBalance(bc, cli.AddressGetBalance)
		return
	}
	if cli.SendCoin {
		if len(flag.Args()) != 5 {
			fmt.Println("invalid command")
			return
		}
		if !IsValidAddress(flag.Arg(0)) {
			fmt.Println("invalid address: ", flag.Arg(0))
			return
		}
		if !IsValidAddress(flag.Arg(1)) {
			fmt.Println("invalid address: ", flag.Arg(1))
			return
		}
		if !IsValidAddress(flag.Arg(3)) {
			fmt.Println("invalid address: ", flag.Arg(3))
			return
		}
		amount, err := strconv.Atoi(flag.Arg(2))
		if err != nil {
			fmt.Println("the amount must be a number")
			return
		}
		cli.Send(bc, flag.Arg(0), flag.Arg(1), int64(amount), flag.Arg(3), flag.Arg(4))
		return
	}
	fmt.Println("invalid command")
}

func (cli *Cli) Print(bc *BlockChain) {
	iter := bc.NewIterator()
	for i := 0; i < cli.PrintNum && i < 20; i++ {
		block := iter.Next()
		if block == nil {
			break
		}
		fmt.Println("===============================================[Block]===============================================")
		fmt.Println(block.String())
	}
}

func (cli *Cli) GetBalance(bc *BlockChain, address string) {
	pubKeyHash, err := GetPubKeyHashFromAddress(address)
	if err != nil {
		fmt.Println("invalid address: ", address)
		return
	}
	_, total := bc.FindUtxo(pubKeyHash)
	fmt.Printf("[%s] remain utxos: %d\n", address, total)
}

func (cli *Cli) Send(bc *BlockChain, from, to string, amount int64, minerPubKey string, data string) {
	// bc, err := GetBlockChain()
	// if err != nil {
	// 	fmt.Println("Can't get blockchain: ", err)
	// }
	// defer bc.Close()

	miningTx := NewMiningTx(minerPubKey, data)
	tx, err := NewTransaction(from, to, amount, bc)
	if err != nil {
		fmt.Printf("Transfer [%d] from [%s] to [%s] failed: %s\n", amount, from, to, err)
		return
	}

	err = bc.AddBlock([]*Transaction{miningTx, tx})
	if err != nil {
		fmt.Printf("Transfer [%d] from [%s] to [%s] failed: %s\n", amount, from, to, err)
		return
	}
	fmt.Printf("Transfer [%d] from [%s] to [%s] success.\n", amount, from, to)
}
