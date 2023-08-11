// blockchain store in bolt database using KV pair `[]byte(Hash): []byte(data)` format
// The last block's hash correspond to `lastBlockHashKey`
package main

import (
	"bytes"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/boltdb/bolt"
)

const (
	dbName     = "blockchain.db"
	bucketName = "bitcoin"

	lastBlockHashKey = "lastBlockHash"
)

type BlockChain struct {
	// Blocks []*Block
	db   *bolt.DB
	tail []byte
}

func CreateBlockChain(address, genesisInfo string) error {
	if IsFileExist(dbName) {
		return errors.New("blockchain store file exists")
	}
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			bucket, err = tx.CreateBucket([]byte(bucketName))
			if err != nil {
				return fmt.Errorf("create bucket fail: %e", err)
			}
		}
		// mining transaction
		miningTx := NewMiningTx(address, genesisInfo)
		// genesis block
		genesisBlock := NewBlock([]*Transaction{miningTx}, []byte{})
		// serialize
		blcokBytes, err2 := genesisBlock.Serialize()
		if err2 != nil {
			return fmt.Errorf("serialize Block fail: %e", err2)
		}
		// store bytes and hash to db
		err = bucket.Put(genesisBlock.Hash, blcokBytes)
		if err != nil {
			return err
		}
		err = bucket.Put([]byte(lastBlockHashKey), genesisBlock.Hash)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func GetBlockChain() (*BlockChain, error) {
	if !IsFileExist(dbName) {
		return nil, errors.New("blockchain store file not exists")
	}

	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("open db fail: %e", err)
	}

	var lastHash []byte
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return errors.New("bucket not exists")
		} else {
			lastHash = bucket.Get([]byte(lastBlockHashKey))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &BlockChain{db, lastHash}, nil
}

func (bc *BlockChain) Close() error {
	return bc.db.Close()
}

func (bc *BlockChain) AddBlock(txs []*Transaction) error {
	lastHash := bc.tail
	block := NewBlock(txs, lastHash)
	blockBytes, err := block.Serialize()
	if err != nil {
		return err
	}

	// varify transactions
	for _, tx := range txs {
		if !bc.VerifyTransaction(tx) {
			log.Println("verify transaction fail")
			return errors.New("invalid transaction")
		}
	}

	err = bc.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return errors.New("bucket not exists")
		}
		err = bucket.Put(block.Hash, blockBytes)
		if err != nil {
			return err
		}
		err = bucket.Put([]byte(lastBlockHashKey), block.Hash)
		if err != nil {
			return err
		}
		bc.tail = block.Hash
		return nil
	})

	return err
}

func (bc *BlockChain) FindTransaction(txid []byte) *Transaction {
	iter := bc.NewIterator()
	for block := iter.Next(); block != nil; block = iter.Next() {
		for _, tx := range block.Transactions {
			if bytes.Equal(tx.Id, txid) {
				return tx
			}
		}
	}
	return nil
}

func (bc *BlockChain) SignTransaction(tx *Transaction, priKey *ecdsa.PrivateKey) bool {
	if tx.IsMiningTx() {
		return true
	}
	log.Println("Start SignTransaction()")
	refedTxs := make(map[string]*Transaction)

	for _, input := range tx.TxInputs {
		refedTx := bc.FindTransaction(input.TxId)
		if refedTx == nil {
			return false
		}
		refedTxs[string(input.TxId)] = refedTx
	}

	return tx.Sign(priKey, refedTxs)
}
func (bc *BlockChain) VerifyTransaction(tx *Transaction) bool {
	if tx.IsMiningTx() {
		return true
	}
	log.Printf("Start VerifyTransaction(%X)\n", tx.Id)
	refedTxs := make(map[string]*Transaction)

	for _, input := range tx.TxInputs {
		refedTx := bc.FindTransaction(input.TxId)
		if refedTx == nil {
			return false
		}
		refedTxs[string(input.TxId)] = refedTx
	}

	return tx.Verify(refedTxs)
}

///////////////////////////////////////////////////////////////////////////

type Iterator struct {
	db          *bolt.DB
	currentHash []byte
}

func (bc *BlockChain) NewIterator() *Iterator {
	return &Iterator{
		db:          bc.db,
		currentHash: bc.tail,
	}
}
func (i *Iterator) Next() *Block {
	// if pointer points to the place before the first
	if reflect.DeepEqual(i.currentHash, make([]byte, 32)) {
		return nil
	}

	var block *Block
	i.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return errors.New("bucket don't exists")
		}
		// get current block
		blockBytes := bucket.Get(i.currentHash)
		b, err := Deserialize(blockBytes)
		if err != nil {
			return err
		}
		block = b
		// update current hash
		i.currentHash = b.PrevHash
		return nil
	})
	return block
}

///////////////////////////////////////////////////////////////////////////

type UTXOInfo struct {
	TxId   []byte
	Index  int64
	Output TxOutput
}

func (bc *BlockChain) FindUtxo(pubKeyHash []byte) ([]UTXOInfo, int64) {
	var utxos []UTXOInfo
	var total int64 = 0

	var spentOutputs = make(map[string][]int64)

	// traversal blockchain
	iter := bc.NewIterator()
	for block := iter.Next(); block != nil; block = iter.Next() {
		// traversal all transactions
		for _, tx := range block.Transactions {
			// traversal all outputs
		TRAVERSAL_OUTPUTS:
			for outputIdx, output := range tx.TxOutputs {
				// is output related to address
				if bytes.Equal(output.ScriptPubKeyHash, pubKeyHash) {
					for _, spentIdx := range spentOutputs[string(tx.Id)] {
						if int64(outputIdx) == spentIdx {
							// if output is spent, continue
							continue TRAVERSAL_OUTPUTS
						}
					}
					// if output is unspent, add it to utxo
					utxos = append(utxos, UTXOInfo{tx.Id, int64(outputIdx), output})
					total += output.Value
				}
			}

			// traversal all inputs except mining transaction
			if tx.IsMiningTx() {
				continue
			}
			for _, input := range tx.TxInputs {
				// if input is related to address, add it to container
				if bytes.Equal(pubKeyHash, GetPubKeyHashFromPubKey(input.PubKey)) {
					spentKey := string(input.TxId)
					spentOutputs[spentKey] = append(spentOutputs[spentKey], input.Index)
				}
			}
		}
	}

	return utxos, total
}

func (bc *BlockChain) FindNeededUtxo(pubKeyHash []byte, amount int64) (int64, map[string][]int64) {
	utxoInfos, total := bc.FindUtxo(pubKeyHash)
	if total < amount {
		return total, nil
	}
	retUtxoInfos := make(map[string][]int64)
	var retTotal int64 = 0
	// traversal utxo, compare with amount
	for _, utxoInfo := range utxoInfos {
		retTotal += utxoInfo.Output.Value
		key := string(utxoInfo.TxId)
		retUtxoInfos[key] = append(retUtxoInfos[key], utxoInfo.Index)
		// if utxo is enough, return
		if total >= amount {
			break
		}
	}
	return retTotal, retUtxoInfos
}
