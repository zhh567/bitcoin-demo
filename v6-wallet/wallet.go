package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

type Wallet struct {
	PriKey []byte
	PubKey []byte
}

// NewWalletKeyPair creates a new wallet with a key pair
func NewWalletKeyPair() *Wallet {
	priKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	x, y := priKey.PublicKey.X.Bytes(), priKey.PublicKey.Y.Bytes()

	return &Wallet{priKey.D.Bytes(), append(x, y...)}
}

func (w *Wallet) GetAddress() string {
	pubKeyHash := GetPubKeyHashFromPubKey(w.PubKey)
	versionedPayload := append([]byte{0x00}, pubKeyHash...)
	checksum := Checksum(versionedPayload)
	fullPayload := append(versionedPayload, checksum...)
	address := base58.Encode(fullPayload)
	return address
}

// Checksum calculates the checksum of payload, return the first 4 bytes
func Checksum(payload []byte) []byte {
	first := sha256.Sum256(payload)
	second := sha256.Sum256(first[:])
	return second[:4]
}

// Crypto public key using SHA256 and RIPEMD160
func GetPubKeyHashFromPubKey(pubKey []byte) []byte {
	sha256Hash := sha256.Sum256(pubKey)
	ripemd160Hash := ripemd160.New()
	ripemd160Hash.Write(sha256Hash[:])
	return ripemd160Hash.Sum(nil)
}
func GetPubKeyHashFromAddress(address string) ([]byte, error) {
	fullPayload := base58.Decode(address)
	if len(fullPayload) != 25 {
		return nil, errors.New("address'length is not 25, invalid address")
	}
	return fullPayload[1 : len(fullPayload)-4], nil
}
func IsValidAddress(address string) bool {
	fullPayload := base58.Decode(address)
	if len(fullPayload) != 25 {
		return false
	}
	versionedPayload := fullPayload[:len(fullPayload)-4]
	checksum := fullPayload[len(fullPayload)-4:]
	return bytes.Equal(Checksum(versionedPayload), checksum)
}

func (w *Wallet) String() string {
	return fmt.Sprintf("Address: %s\nPubKey: %#X\nPriKey:%#X", w.GetAddress(), w.PubKey, w.PriKey)
}
