package main

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"log"
	"time"
)

func NewGenesisBlock() *BlockBroadArray {
	block, _ := NewBlockArray([]byte{}, nil, true)
	pow := block.NewPoW()
	res := pow.proof(true)
	if res {
		return block
	}
	return nil
}

//var wallets, _ = NewWallets()

var chain = &BlockChainArray{nil, nil, 0}

func NewBlockChain(block *BlockBroadArray) {
	Genesis := &BlockArray{*block, 0, nil}
	chain.start = Genesis
	chain.tails = append(chain.tails, Genesis)
	chain.numTail = 1
	//var tails = []*BlockArray{Genesis}
	//chain = &BlockChainArray{Genesis, tails, 1}
}

func (chain *BlockChainArray) FindUnspentTransactions(pubKeyHash []byte) []transaction {
	var unspentTXs []transaction
	spentTXs := make(map[string][]int)

	longest := 0
	longestid := 0
	for i := 0; i < chain.numTail; i++ {
		if chain.tails[i].num > longest {
			longestid = i
			longest = chain.tails[i].num
		}
	}
	block := chain.tails[longestid]
	for i := 0; i < block.num; i++ {
		for _, tx := range block.blockBroad.Data {
			txID := hex.EncodeToString(tx.Transaction.id)
		Outputs:
			for outIdx, out := range tx.Transaction.Vout {
				if spentTXs[txID] != nil {
					for _, spentOutIdx := range spentTXs[txID] {
						if spentOutIdx == outIdx {
							continue Outputs
						}
					}
				}
				if out.IsLockedWithKey(pubKeyHash) {
					unspentTXs = append(unspentTXs, tx)
				}
			}
			if tx.IsCoinbase() == false {
				for _, in := range tx.Transaction.Vin {
					if in.UsesKey(pubKeyHash) {
						inTxID := hex.EncodeToString(in.id)
						spentTXs[inTxID] = append(spentTXs[inTxID], in.Vout)
					}
				}
			}
		}
		block = block.prevPtr
	}
	return unspentTXs
}

func (chain *BlockChainArray) FindUTXO(pubKeyHash []byte) []txOutput {
	var UTXOs []txOutput
	unspentTransactions := chain.FindUnspentTransactions(pubKeyHash)
	for _, tx := range unspentTransactions {
		for _, out := range tx.Transaction.Vout {
			if out.IsLockedWithKey(pubKeyHash) {
				UTXOs = append(UTXOs, out)
			}
		}
	}
	return UTXOs
}

func (chain *BlockChainArray) FindSpendableOutputs(pubKeyHash []byte, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	unspentTXs := chain.FindUnspentTransactions(pubKeyHash)
	total := 0

Work:
	for _, tx := range unspentTXs {
		txID := hex.EncodeToString(tx.Transaction.id)
		for outIdx, out := range tx.Transaction.Vout {
			if out.IsLockedWithKey(pubKeyHash) && total < amount {
				total += out.value
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)

				if total >= amount {
					break Work
				}
			}
		}
	}
	return total, unspentOutputs
}

func SendUTXO(from, to string, amount int, chain *BlockChainArray) *transaction {
	var inputs []txInput
	var outputs []txOutput

	wallet := wallets.GetWallet(from)
	pubKeyHash := HashPubKey(wallet.PublicKey)
	tot, validOutputs := chain.FindSpendableOutputs(pubKeyHash, amount)
	if tot < amount {
		//log.Panic("ERROR:Not enough tokens...")
		return nil
	}
	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}
		for _, out := range outs {
			input := txInput{txID, out, nil, wallet.PublicKey}
			inputs = append(inputs, input)
		}
	}
	outputs = append(outputs, *NewTXOutput(amount, to))
	if tot > amount {
		outputs = append(outputs, *NewTXOutput(tot-amount, from))
	}
	tx := transaction{time.Now().UnixNano(), Transaction{nil, inputs, outputs}}
	tx.Transaction.id = tx.Transaction.Hash()
	chain.SignTransaction(&tx, wallet.PrivateKey)

	return &tx
}

func (chain *BlockChainArray) FindTransaction(ID []byte) (transaction, error) {
	longest := 0
	longestid := 0
	for i := 0; i < chain.numTail; i++ {
		if chain.tails[i].num > longest {
			longestid = i
			longest = chain.tails[i].num
		}
	}
	block := chain.tails[longestid]
	for i := 0; i < block.num; i++ {
		for _, tx := range block.blockBroad.Data {
			if bytes.Compare(tx.Transaction.id, ID) == 0 {
				return tx, nil
			}
		}
		block = block.prevPtr
	}
	return transaction{}, errors.New("Transaction is not found.")
}

func (chain *BlockChainArray) SignTransaction(tx *transaction, privKey ecdsa.PrivateKey) {
	prevTXs := make(map[string]transaction)
	for _, vin := range tx.Transaction.Vin {
		prevTX, err := chain.FindTransaction(vin.id)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.Transaction.id)] = prevTX
	}
	tx.Sign(privKey, prevTXs)
}

func (chain *BlockChainArray) VerifyTransaction(tx *transaction) bool {
	prevTXs := make(map[string]transaction)
	for _, vin := range tx.Transaction.Vin {
		prevTX, err := chain.FindTransaction(vin.id)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.Transaction.id)] = prevTX
	}
	return tx.Verify(prevTXs)
}
