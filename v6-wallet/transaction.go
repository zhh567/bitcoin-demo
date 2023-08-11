package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"math/big"
	"time"
)

const (
	Reward int64 = 17
)

// 1. 交易id
// 2. 交易输出input，由历史中某个output转换而来（可有多个）
//  1. 引用的交易id
//  2. 此交易中对应output的索引
//  3. **解锁脚本**（发起者的签名、公钥）
//
// 3. 交易输出output， 表明流向（可有多个）
//  1. **锁定脚本**（收款人的地址，可反推出公钥哈希）
//  2. 转账金额
//
// 4. 时间戳
type Transaction struct {
	Id        []byte
	TxInputs  []TxInput
	TxOutputs []TxOutput
	TimeStamp int64
}

type TxInput struct {
	TxId  []byte
	Index int64

	ScriptSig []byte
	PubKey    []byte
}

type TxOutput struct {
	ScriptPubKeyHash []byte // receiver's public key hash
	Value            int64
}

func (t *Transaction) SetHash() {
	t.Id = nil // if don't set it, multiple calls will get different results

	var buf bytes.Buffer
	e := gob.NewEncoder(&buf)

	err := e.Encode(&t)
	if err != nil {
		fmt.Println("encode transaction fail: ", err)
		panic(err)
	}
	hashBytes := sha256.Sum256(buf.Bytes())

	t.Id = hashBytes[:]
}

func NewMiningTx(
	address string, // miner's public key
	data string, // mining reward have no input, write data to sig
) *Transaction {
	log.Println("Start creating new mining transaction")
	minerPubKeyHash, err := GetPubKeyHashFromAddress(address)
	if err != nil {
		panic("invalid address")
	}

	txInput := TxInput{nil, 0, []byte(data), nil}
	txOutput := TxOutput{minerPubKeyHash, Reward}

	tx := &Transaction{
		TxInputs:  []TxInput{txInput},
		TxOutputs: []TxOutput{txOutput},
		TimeStamp: time.Now().Unix(),
	}
	tx.SetHash()

	log.Println("New mining transaction")
	return tx
}

func NewTransaction(
	from string, // sender's address
	to string, // receiver's address
	amount int64, // transfer amount
	bc *BlockChain,
) (*Transaction, error) {
	// 1. 遍历账本，找到关于from的utxo集合，返回总金额
	// 2. 金额不足，创建失败
	// 3. 拼接 inputs
	//     遍历utxo集合，每个output都转换为一个input
	// 4. 拼接 outputs
	//     创建一个属于to的output
	//     如果总金额大于转账金额，给from创建找零output
	// 5. 设置hash
	wm := NewWalletManager()
	if wm == nil {
		return nil, errors.New("can't get wallet manager")
	}
	wallet, ok := wm.Wallets[from]
	if !ok {
		return nil, errors.New("can't find sender's wallet")
	}

	fromPubKeyHash, err := GetPubKeyHashFromAddress(from)
	if err != nil {
		return nil, errors.New("invalid address")
	}
	toPubKeyHash, err := GetPubKeyHashFromAddress(to)
	if err != nil {
		return nil, errors.New("invalid address")
	}
	total, utxoInfos := bc.FindNeededUtxo(fromPubKeyHash, amount)
	if total < amount {
		log.Printf("Transfer %s to %s: not enough money\n", from, to)
		return nil, errors.New("not enough money")
	}
	inputs := make([]TxInput, 0)
	outputs := make([]TxOutput, 0)

	for txId, indexes := range utxoInfos {
		for _, index := range indexes {
			input := TxInput{[]byte(txId), index, nil, wallet.PubKey}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, TxOutput{toPubKeyHash, amount})
	if total > amount {
		outputs = append(outputs, TxOutput{fromPubKeyHash, (total - amount)})
	}

	tx := &Transaction{
		TxInputs:  inputs,
		TxOutputs: outputs,
		TimeStamp: time.Now().Unix(),
	}

	tx.SetHash()

	// i, i2 := elliptic.P256().ScalarBaseMult(wallet.PriKey)
	priKey := ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: elliptic.P256(),
			X:     new(big.Int).SetBytes(wallet.PubKey[:len(wallet.PubKey)/2]),
			Y:     new(big.Int).SetBytes(wallet.PubKey[len(wallet.PubKey)/2:]),
		},
		D: new(big.Int).SetBytes(wallet.PriKey),
	}
	log.Printf("Created private key:\n\t%X\n\t%#X\n\t%#X", priKey.D.Bytes(), priKey.PublicKey.X.Bytes(), priKey.PublicKey.Y.Bytes())
	if !bc.SignTransaction(tx, &priKey) {
		log.Println("sign transaction failed")
		return nil, errors.New("sign transaction failed")
	}

	log.Println("Create new transaction")
	return tx, nil
}

func (tx *Transaction) IsMiningTx() bool {
	return len(tx.TxInputs) == 1 && len(tx.TxInputs[0].TxId) == 0 && tx.TxInputs[0].Index == 0
}

// copy transaction, remove signature and public key
func (tx *Transaction) TrimmedCopy() *Transaction {
	inputs := make([]TxInput, 0)
	outputs := make([]TxOutput, 0)

	for _, input := range tx.TxInputs {
		inputs = append(inputs, TxInput{input.TxId, input.Index, nil, nil})
	}
	for _, output := range tx.TxOutputs {
		outputs = append(outputs, TxOutput{output.ScriptPubKeyHash, output.Value})
	}

	return &Transaction{
		Id:        tx.Id,
		TxInputs:  inputs,
		TxOutputs: outputs,
		TimeStamp: tx.TimeStamp,
	}
}

func (tx *Transaction) Sign(priKey *ecdsa.PrivateKey, referencedTxs map[string]*Transaction) bool {
	if tx.IsMiningTx() {
		return true
	}
	log.Println("Start Transaction.Sign()")
	txCopy := tx.TrimmedCopy()
	for i, input := range txCopy.TxInputs {
		refedTx := referencedTxs[string(input.TxId)]
		if refedTx == nil {
			return false
		}

		refedOutput := refedTx.TxOutputs[input.Index]
		txCopy.TxInputs[i].PubKey = refedOutput.ScriptPubKeyHash
		txCopy.SetHash()
		txCopy.TxInputs[i].PubKey = nil
		hashData := txCopy.Id
		log.Printf("In Sign() hashData: %X\n", hashData)

		r, s, err := ecdsa.Sign(rand.Reader, priKey, hashData)
		if err != nil {
			return false
		}
		tx.TxInputs[i].ScriptSig = append(r.Bytes(), s.Bytes()...)
		log.Printf("In Sign() signature: [%X]\n", tx.TxInputs[i].ScriptSig)
		log.Printf("In Sign() pubKey: %X\n", priKey.PublicKey)
	}
	return true
}

func (tx *Transaction) Verify(refedTxs map[string]*Transaction) bool {
	log.Println("Start Transaction.Verify()")
	// copy a transaction, remove signature and public key
	txCopy := tx.TrimmedCopy()
	// traverse all inputs
	for i, input := range tx.TxInputs {
		// get referenced transaction
		refedTx := refedTxs[string(input.TxId)]
		if refedTx == nil {
			log.Printf("can't find referenced transaction: %X\n", input.TxId)
			return false
		}
		// get referenced output, get public key hash, restore pubKey field
		refedOutput := refedTx.TxOutputs[input.Index]
		txCopy.TxInputs[i].PubKey = refedOutput.ScriptPubKeyHash
		txCopy.SetHash()
		txCopy.TxInputs[i].PubKey = nil
		hashData := txCopy.Id
		log.Printf("In Verify() hashData: %X\n", hashData)
		log.Printf("In Verify() signature: [%X]\n", input.ScriptSig)

		// verify signature
		signature := input.ScriptSig
		var r, s, x, y big.Int
		r.SetBytes(signature[:len(signature)/2])
		s.SetBytes(signature[len(signature)/2:])

		x.SetBytes(input.PubKey[:32])
		y.SetBytes(input.PubKey[32:])

		pubKeyNew := ecdsa.PublicKey{Curve: elliptic.P256(), X: &x, Y: &y}
		log.Printf("In Verify() pubKey: %X\n", pubKeyNew)

		if !ecdsa.Verify(&pubKeyNew, hashData, &r, &s) {
			log.Println("verify signature failed: ecdsa.Verify() return false")
			return false
		}
	}

	return true
}

///////////////////////////////////////////////////////////////////

func (t *TxInput) String() string {
	var format string
	if t.TxId != nil {
		// common transaction
		format = "TxID:\t\t%X\nIndex:\t\t%d\nPubKey:\t\t%X\nScriptSig:\t%X"
		return fmt.Sprintf(format, t.TxId, t.Index, t.PubKey, t.ScriptSig)
	} else {
		// coinbase
		return fmt.Sprintf("<Coinbase Transaction>\nside message:%s", t.ScriptSig)
	}
}
func (t *TxOutput) String() string {
	format := `%d => %X`
	return fmt.Sprintf(format, t.Value, t.ScriptPubKeyHash)
}
func (t *Transaction) String() string {
	strInputs := ""
	strOutputs := ""

	for _, input := range t.TxInputs {
		strInputs += input.String() + "\n"
	}
	for _, output := range t.TxOutputs {
		strOutputs += output.String() + "\n"
	}

	//-------------------------------------[Transaction]-------------------------------------
	format := `
                                     [Transaction]                                     
[Id]
%X
[TxInputs]
%s[TxOutputs]
%s[TimeStamp]
%d
`
	return fmt.Sprintf(format, t.Id, strInputs, strOutputs, t.TimeStamp)
}
