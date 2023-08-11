package main

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const (
	walletFile = "wallet.dat"
)

type WalletManager struct {
	Wallets map[string]*Wallet
}

func NewWalletManager() *WalletManager {
	wm := &WalletManager{Wallets: make(map[string]*Wallet)}
	wm.LoadFile()
	return wm
}

func (wm *WalletManager) CreateWallet() string {
	wallet := NewWalletKeyPair()
	address := wallet.GetAddress()
	wm.Wallets[address] = wallet
	wm.SaveFile()
	return address
}

func (wm *WalletManager) GetWallet(address string) *Wallet {
	return wm.Wallets[address]
}

func (wm *WalletManager) SaveFile() {
	var buffer bytes.Buffer
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(wm)
	if err != nil {
		fmt.Println("encode wallet file fail: ", err)
		panic(err)
	}

	f, err := os.OpenFile(walletFile, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("open wallet file fail: ", err)
		panic(err)
	}
	defer f.Close()

	_, err2 := f.Write(buffer.Bytes())
	if err2 != nil {
		fmt.Println("write wallet file fail: ", err2)
		panic(err2)
	}
}

func (wm *WalletManager) LoadFile() {
	if !IsFileExist(walletFile) {
		return
	}
	// open wallet file
	fp, err := os.Open(walletFile)
	if err != nil {
		fmt.Println("open wallet file fail: ", err)
		panic(err)
	}
	defer fp.Close()

	b, err2 := io.ReadAll(fp)
	if err2 != nil {
		fmt.Println("read wallet file fail: ", err2)
		panic(err2)
	}
	// unserialize wallet file
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(b))
	err3 := decoder.Decode(&wm)
	if err3 != nil {
		fmt.Println("decode wallet file fail: ", err3)
		panic(err3)
	}
	log.Printf("Load wallet:\n%s\n", wm)
}

func (wm *WalletManager) ListAllAddresses() []string {
	addresses := make([]string, 0)
	for address, wallet := range wm.Wallets {
		addresses = append(addresses, address+" : "+fmt.Sprintf("%X", wallet.PriKey))
	}
	return addresses
}

func (wm *WalletManager) String() string {
	str := strings.Builder{}
	for address, wallet := range wm.Wallets {
		str.WriteString(fmt.Sprintf("Address: %s\n", address))
		str.WriteString(fmt.Sprintf("Public Key: %X\n", wallet.PubKey))
		str.WriteString(fmt.Sprintf("Private Key: %X\n", wallet.PriKey))
	}
	return str.String()
}
