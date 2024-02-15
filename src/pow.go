package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/big"
	"strconv"
	"time"
)

type PoW struct {
	//blockHeader *BlockHeader
	block  *BlockBroad
	target *big.Int
}

type PoWArray struct {
	block  *BlockBroadArray
	target *big.Int
}
type PoWTree struct {
	block  *BlockBroadTree
	target *big.Int
}

var targetBits = 24
var size = 8
var nodeId string

//func NewBlock(hash []byte) *BlockBroad {
//	blockBroad := &BlockBroad{BlockHeader{time.Now().UnixNano(), hash}, pool.head.Transaction, []byte{}, 0}
//	return blockBroad
//}

func NewBlockArray(hash []byte, chain *BlockChainArray, genesis bool) (*BlockBroadArray, []Request) {
	var requests []Request
	id, _ := strconv.Atoi(node)
	//fmt.Println(wallets.GetAddresses())
	nodeId = wallets.GetAddresses()[id-16]
	cointx := NewCoinbaseTx(nodeId, "")
	blockBroad := BlockBroadArray{BlockHeader{time.Now().UnixNano(), hash}, []transaction{*cointx}, []byte{}, 0}
	block := &BlockArray{blockBroad, 0, nil}
	if genesis == false {
		idx := chain.AddBlockArray(block)
		requestList := pool.getHead(size)
		for _, request := range requestList {
			tx := SendUTXO(request.From, request.To, request.Amount, chain)
			if tx != nil {
				ver := chain.VerifyTransaction(tx)
				fmt.Println("Verify", tx, ver)
				blockBroad.Data = append(blockBroad.Data, *tx)
				chain.tails[idx].blockBroad.Data = append(chain.tails[idx].blockBroad.Data, *tx)
			}
			requests = append(requests, request)
		}
		chain.Delete(idx)
	}
	return &blockBroad, requests
}

func (block *BlockBroad) NewPoW() *PoW {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))
	pow := &PoW{block, target}
	return pow
}

func (block *BlockBroadArray) NewPoW() *PoWArray {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))
	pow := &PoWArray{block, target}
	return pow
}

func IntToBytes(n int64) []byte {
	byteBuf := bytes.NewBuffer([]byte{})
	binary.Write(byteBuf, binary.BigEndian, n)
	return byteBuf.Bytes()
}

func IntToBytes_(n int) []byte {
	byteBuf := bytes.NewBuffer([]byte{})
	binary.Write(byteBuf, binary.BigEndian, n)
	return byteBuf.Bytes()
}

//func (pow *PoW) proof() bool {
//	header := bytes.Join(
//		[][]byte{
//			IntToBytes(pow.block.BlockHeader.TimeStamp),
//			pow.block.BlockHeader.PrevHash,
//			IntToBytes(pow.block.Data.TimeStamp),
//			[]byte(pow.block.Data.Transaction),
//		},
//		[]byte{},
//	)
//	nonce := int64(0)
//	//maxNonce := int64(256)
//	for breakPow == false {
//		var tmp []byte
//		tmp = append(header)
//		tmp = append(IntToBytes(nonce))
//		hash := sha256.Sum256(tmp)
//		var hashInt big.Int
//		hashInt.SetBytes(hash[:])
//		if hashInt.Cmp(pow.target) == -1 {
//			pow.block.Nonce = nonce
//			pow.block.Hash = hash[:]
//			return true
//		} else {
//			nonce++
//		}
//	}
//	breakPow = false
//	return false
//}

func (pow *PoWArray) proof(gen bool) bool {
	var header []byte
	header = bytes.Join(
		[][]byte{
			header,
			IntToBytes(pow.block.BlockHeader.TimeStamp),
			pow.block.BlockHeader.PrevHash,
		},
		[]byte{},
	)
	for i := 0; i < len(pow.block.Data); i++ {
		header = bytes.Join(
			[][]byte{
				header,
				IntToBytes(pow.block.Data[i].TimeStamp),
				pow.block.Data[i].Transaction.id,
			},
			[]byte{},
		)
		for j := 0; j < len(pow.block.Data[i].Transaction.Vin); j++ {
			header = bytes.Join(
				[][]byte{
					header,
					pow.block.Data[i].Transaction.Vin[j].id,
					IntToBytes_(pow.block.Data[i].Transaction.Vin[j].Vout),
					pow.block.Data[i].Transaction.Vin[j].Sign,
					pow.block.Data[i].Transaction.Vin[j].PubKey,
				},
				[]byte{},
			)
		}
		for j := 0; j < len(pow.block.Data[i].Transaction.Vout); j++ {
			header = bytes.Join(
				[][]byte{
					header,
					IntToBytes_(pow.block.Data[i].Transaction.Vout[j].value),
					pow.block.Data[i].Transaction.Vout[j].PubKeyHash,
				},
				[]byte{},
			)
		}
	}

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
