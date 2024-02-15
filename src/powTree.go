package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
	"strconv"
	"time"
)

//func NewBlock(hash []byte) *BlockBroad {
//	blockBroad := &BlockBroad{BlockHeader{time.Now().UnixNano(), hash}, pool.head.Transaction, []byte{}, 0}
//	return blockBroad
//}

func NewBlockTree(hash []byte, chain *BlockChainTree, genesis bool) (*BlockBroadTree, []Request) {
	var requests []Request
	id, _ := strconv.Atoi(node)
	//fmt.Println(wallets.GetAddresses())
	nodeId = wallets.GetAddresses()[id-16]
	cointx := NewCoinbaseTx(nodeId, "")
	tree, root := GenTree([]transaction{*cointx})
	blockBroad := BlockBroadTree{BlockHeader{time.Now().UnixNano(), hash}, []transaction{*cointx}, []byte{}, 0, root, tree}
	block := &BlockTree{blockBroad, 0, nil}
	if genesis == false {
		idx := chain.AddBlockTree(block)
		requestList := pool.getHead(size)
		for _, request := range requestList {
			tx := SendUTXOTree(request.From, request.To, request.Amount, chain)
			if tx != nil {
				ver := chain.VerifyTransactionTree(tx)
				fmt.Println("Verify", tx, ver)
				blockBroad.Data = append(blockBroad.Data, *tx)
				chain.tails[idx].blockBroad.Data = append(chain.tails[idx].blockBroad.Data, *tx)
			}
			requests = append(requests, request)
		}
		chain.DeleteTree(idx)
		blockBroad.tree, blockBroad.Root = GenTree(blockBroad.Data)
		//if blockBroad.tree == nil {
		//	fmt.Println("NewBlockTree err")
		//}
	}
	return &blockBroad, requests
}

func (block *BlockBroadTree) NewPoWTree() *PoWTree {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))
	pow := &PoWTree{block, target}
	return pow
}

func (pow *PoWTree) proofTree(gen bool) bool {
	timestamp := []byte(strconv.FormatInt(pow.block.BlockHeader.TimeStamp, 10))
	var data []byte
	data = bytes.Join([][]byte{data, pow.block.Root}, []byte{})
	header := bytes.Join([][]byte{pow.block.BlockHeader.PrevHash, data, timestamp}, []byte{})

	//var header []byte
	//header = bytes.Join(
	//	[][]byte{
	//		header,
	//		IntToBytes(pow.block.BlockHeader.TimeStamp),
	//		pow.block.BlockHeader.PrevHash,
	//	},
	//	[]byte{},
	//)
	//for i := 0; i < len(pow.block.Data); i++ {
	//	header = bytes.Join(
	//		[][]byte{
	//			header,
	//			IntToBytes(pow.block.Data[i].TimeStamp),
	//			pow.block.Data[i].Transaction.id,
	//		},
	//		[]byte{},
	//	)
	//	for j := 0; j < len(pow.block.Data[i].Transaction.Vin); j++ {
	//		header = bytes.Join(
	//			[][]byte{
	//				header,
	//				pow.block.Data[i].Transaction.Vin[j].id,
	//				IntToBytes_(pow.block.Data[i].Transaction.Vin[j].Vout),
	//				pow.block.Data[i].Transaction.Vin[j].Sign,
	//				pow.block.Data[i].Transaction.Vin[j].PubKey,
	//			},
	//			[]byte{},
	//		)
	//	}
	//	for j := 0; j < len(pow.block.Data[i].Transaction.Vout); j++ {
	//		header = bytes.Join(
	//			[][]byte{
	//				header,
	//				IntToBytes_(pow.block.Data[i].Transaction.Vout[j].value),
	//				pow.block.Data[i].Transaction.Vout[j].PubKeyHash,
	//			},
	//			[]byte{},
	//		)
	//	}
	//}

	nonce := int64(0)
	//maxNonce := int64(256)
	for breakPow == false {
		if gen && IsGen == false {
			break
		}
		var tmp []byte
		//tmp = append(header)
		//tmp = append(header, IntToBytes(nonce))
		tmp = bytes.Join(
			[][]byte{
				header,
				IntToBytes(nonce),
			},
			[]byte{},
		)
		//var tmp []byte
		//tmp = append(header)
		//tmp = append(IntToBytes(nonce))
		hash := sha256.Sum256(tmp)
		var hashInt big.Int
		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(pow.target) == -1 {
			pow.block.Nonce = nonce
			pow.block.Hash = hash[:]
			return true
		} else {
			nonce++
		}
	}
	breakPow = false
	return false
}
