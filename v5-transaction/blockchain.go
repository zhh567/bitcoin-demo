// blockchain store in bolt database using KV pair `[]byte(Hash): []byte(data)` format
// The last block's hash correspond to `lastBlockHashKey`
package main

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/boltdb/bolt"
)

const (
	genesisInfo = "Init blockchain demo"

	dbName     = "blockchain.db"
	bucketName = "bitcoin"

	lastBlockHashKey = "lastBlockHash"
)

var (
	blockchain *BlockChain = nil
)

type BlockChain struct {
	// Blocks []*Block
	db   *bolt.DB
	tail []byte
}

// single mode, don't consider concurrent sence
func GetBlockChain() (*BlockChain, error) {
	if blockchain != nil {
		return blockchain, nil
	}

	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		return nil, err
	}

	var lastHash []byte
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		// bucket is empty, init
		if bucket == nil {
			bucket, err = tx.CreateBucket([]byte(bucketName))
			if err != nil {
				return fmt.Errorf("create bucket fail: %e", err)
			}

			b := [32]byte{0}
			// genesis block
			miningTx := NewMiningTx("Genesis", genesisInfo)
			genesisBlock := NewBlock([]*Transaction{miningTx}, b[:])
			blcokBytes, err2 := genesisBlock.Serialize()
			if err2 != nil {
				return fmt.Errorf("serialize Block fail: %e", err2)
			}
			bucket.Put(genesisBlock.Hash, blcokBytes)
			bucket.Put([]byte(lastBlockHashKey), genesisBlock.Hash)

			lastHash = genesisBlock.Hash
		} else {
			// bucket exists, read data
			lastHash = bucket.Get([]byte(lastBlockHashKey))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	blockchain = &BlockChain{
		db:   db,
		tail: lastHash,
	}

	return blockchain, nil
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

// ////////////////////////////////////////////////////////////////////////////////////////////
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

// /////////////////////////////////////////////////////////////////////////

type UTXOInfo struct {
	TxId   []byte
	Index  int64
	Output TxOutput
}

func (bc *BlockChain) FindUtxo(address string) ([]UTXOInfo, int64) {
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
				if output.ScriptPubKey == address {
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
			// traversal all inputs
			for _, input := range tx.TxInputs {
				// if input is related to address, add it to container
				if string(input.ScriptSig) == address {
					spentKey := string(input.TxId)
					spentOutputs[spentKey] = append(spentOutputs[spentKey], input.Index)
				}
			}
		}
	}

	return utxos, total
}

func (bc *BlockChain) FindNeededUtxo(address string, amount int64) (int64, map[string][]int64) {
	utxoInfos, total := bc.FindUtxo(address)
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
