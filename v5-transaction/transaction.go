package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
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
	TxId      []byte
	Index     int64
	ScriptSig []byte
}

type TxOutput struct {
	ScriptPubKey string
	Value        int64
}

func (t *Transaction) SetHash() {
	t.Id = nil // if don't set it, multiple calls will get different results

	var buf bytes.Buffer
	e := gob.NewEncoder(&buf)

	err := e.Encode(&t)
	if err != nil {
		panic(err)
	}
	hashBytes := sha256.Sum256(buf.Bytes())

	t.Id = hashBytes[:]
}

func NewMiningTx(
	minerPubKey string, // miner's public key
	data string, // mining reward have no input, write data to sig
) *Transaction {
	txInput := TxInput{nil, 0, []byte(data)}
	txOutput := TxOutput{minerPubKey, Reward}

	tx := &Transaction{
		TxInputs:  []TxInput{txInput},
		TxOutputs: []TxOutput{txOutput},
		TimeStamp: time.Now().Unix(),
	}
	tx.SetHash()

	return tx
}

func NewTransaction(
	from string, // sender's address
	to string, // receiver's address
	amount int64, // transfer amount
	bc *BlockChain,
) *Transaction {
	// 1. 遍历账本，找到关于from的utxo集合，返回总金额
	// 2. 金额不足，创建失败
	// 3. 拼接 inputs
	//     遍历utxo集合，每个output都转换为一个input
	// 4. 拼接 outputs
	//     创建一个属于to的output
	//     如果总金额大于转账金额，给from创建找零output
	// 5. 设置hash
	total, utxoInfos := bc.FindNeededUtxo(from, amount)
	if total < amount {
		fmt.Printf("Transfer %s to %s: not enough money\n", from, to)
		return nil
	}
	inputs := make([]TxInput, 0)
	outputs := make([]TxOutput, 0)

	for txId, indexes := range utxoInfos {
		for _, index := range indexes {
			input := TxInput{[]byte(txId), index, []byte(from)}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, TxOutput{to, amount})
	if total > amount {
		outputs = append(outputs, TxOutput{from, (total - amount)})
	}

	tx := &Transaction{
		TxInputs:  inputs,
		TxOutputs: outputs,
		TimeStamp: time.Now().Unix(),
	}

	tx.SetHash()

	return tx
}

func (tx *Transaction) IsMiningTx() bool {
	return len(tx.TxInputs) == 1 && len(tx.TxInputs[0].TxId) == 0 && tx.TxInputs[0].Index == 0
}

///////////////////////////////////////////////////////////////////

func (t *TxInput) String() string {
	format := `%X : %d : %s`
	return fmt.Sprintf(format, t.TxId, t.Index, t.ScriptSig)
}
func (t *TxOutput) String() string {
	format := `%d => %s`
	return fmt.Sprintf(format, t.Value, t.ScriptPubKey)
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
