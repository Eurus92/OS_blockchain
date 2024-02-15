package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"strconv"
)

type BlockChainNew struct {
	start *BlockNew
	tails []*BlockNew
	//numTail int64
	numTail int
}
type BlockChainArray struct {
	start   *BlockArray
	tails   []*BlockArray
	numTail int
}
type BlockChainTree struct {
	start   *BlockTree
	tails   []*BlockTree
	numTail int
}

func (b *BlockNew) CompHashNew() []byte {
	timestamp := []byte(strconv.FormatInt(b.blockBroad.BlockHeader.TimeStamp, 10))
	var data []byte
	//var trans transaction
	//for trans = range b.Data {
	//	data = bytes.Join([][]byte{data, []byte(strconv.FormatInt(trans.TimeStamp, 10)), []byte(trans.Transaction)}, []byte{})
	//}
	h := b.blockBroad.Data.Transaction.Hash()
	data = bytes.Join([][]byte{data, []byte(strconv.FormatInt(b.blockBroad.Data.TimeStamp, 10)), h}, []byte{})
	headers := bytes.Join([][]byte{b.blockBroad.BlockHeader.PrevHash, data, timestamp}, []byte{})
	Hash := sha256.Sum256(headers)
	hash := Hash[:]
	return hash
}

func (b *BlockArray) CompHashArray() []byte {
	//timestamp := []byte(strconv.FormatInt(b.blockBroad.BlockHeader.TimeStamp, 10))
	//var data []byte
	//for i := 0; i < len(b.blockBroad.Data); i++ {
	//	h := b.blockBroad.Data[i].Transaction.Hash()
	//	data = bytes.Join([][]byte{data, []byte(strconv.FormatInt(b.blockBroad.Data[i].TimeStamp, 10)), h}, []byte{})
	//}
	//headers := bytes.Join([][]byte{b.blockBroad.BlockHeader.PrevHash, data, timestamp}, []byte{})
	block := b.blockBroad
	var header []byte
	header = bytes.Join(
		[][]byte{
			header,
			IntToBytes(block.BlockHeader.TimeStamp),
			block.BlockHeader.PrevHash,
		},
		[]byte{},
	)
	for i := 0; i < len(block.Data); i++ {
		header = bytes.Join(
			[][]byte{
				header,
				IntToBytes(block.Data[i].TimeStamp),
				block.Data[i].Transaction.id,
			},
			[]byte{},
		)
		for j := 0; j < len(block.Data[i].Transaction.Vin); j++ {
			header = bytes.Join(
				[][]byte{
					header,
					block.Data[i].Transaction.Vin[j].id,
					IntToBytes_(block.Data[i].Transaction.Vin[j].Vout),
					block.Data[i].Transaction.Vin[j].Sign,
					block.Data[i].Transaction.Vin[j].PubKey,
				},
				[]byte{},
			)
		}
		for j := 0; j < len(block.Data[i].Transaction.Vout); j++ {
			header = bytes.Join(
				[][]byte{
					header,
					IntToBytes_(block.Data[i].Transaction.Vout[j].value),
					block.Data[i].Transaction.Vout[j].PubKeyHash,
				},
				[]byte{},
			)
		}
	}
	var tmp []byte
	tmp = bytes.Join(
		[][]byte{
			header,
			IntToBytes(block.Nonce),
		},
		[]byte{},
	)
	//header = append(IntToBytes(block.Nonce))
	Hash := sha256.Sum256(tmp)
	hash := Hash[:]
	return hash
}

func (b *BlockTree) CompHashTree() []byte {
	timestamp := []byte(strconv.FormatInt(b.blockBroad.BlockHeader.TimeStamp, 10))
	var data []byte
	data = bytes.Join([][]byte{data, b.blockBroad.Root}, []byte{})
	headers := bytes.Join([][]byte{b.blockBroad.BlockHeader.PrevHash, data, timestamp}, []byte{})
	headers = bytes.Join(
		[][]byte{
			headers,
			IntToBytes(b.blockBroad.Nonce),
		},
		[]byte{},
	)
	Hash := sha256.Sum256(headers)
	hash := Hash[:]
	return hash
}

func (chain *BlockChainNew) AddBlock(block *BlockNew) {
	block.blockBroad.Hash = block.CompHashNew()
	for i := 0; i < chain.numTail; i++ {
		if bytes.Equal(chain.tails[i].blockBroad.Hash, block.blockBroad.BlockHeader.PrevHash) {
			block.prevPtr = chain.tails[i]
			//block.prevPtr = prev
			block.num = chain.tails[i].num + 1
			chain.tails[i] = block
			return
		}
	}
	for i := 0; i < chain.numTail; i++ {
		tmp := chain.tails[i].prevPtr
		for j := 0; j < tmp.num; j++ {
			if bytes.Equal(tmp.blockBroad.Hash, block.blockBroad.BlockHeader.PrevHash) {
				block.prevPtr = tmp
				block.num = tmp.num + 1
				chain.numTail += 1
				chain.tails = append(chain.tails, block)
				return
			} else {
				tmp = tmp.prevPtr
			}
		}
	}
}
func (chain *BlockChainArray) AddBlockArray(block *BlockArray) int {
	block.blockBroad.Hash = block.CompHashArray()
	for i := 0; i < chain.numTail; i++ {
		if bytes.Equal(chain.tails[i].blockBroad.Hash, block.blockBroad.BlockHeader.PrevHash) {
			block.prevPtr = chain.tails[i]
			block.num = chain.tails[i].num + 1
			chain.tails[i] = block
			return i
		}
	}
	for i := 0; i < chain.numTail; i++ {
		tmp := chain.tails[i].prevPtr
		if tmp == nil {
			continue
		}
		num := tmp.num
		for j := 0; j < num; j++ {
			if bytes.Equal(tmp.blockBroad.Hash, block.blockBroad.BlockHeader.PrevHash) {
				block.prevPtr = tmp
				block.num = tmp.num + 1
				chain.numTail += 1
				chain.tails = append(chain.tails, block)
				return chain.numTail - 1
			} else {
				tmp = tmp.prevPtr
			}
		}
	}
	return 0
}

func (chain *BlockChainArray) Delete(tailId int) {
	if chain.tails[tailId].num == 0 {
		chain.numTail--
		chain.tails = chain.tails[:len(chain.tails)-1]
	} else {
		chain.tails[tailId] = chain.tails[tailId].prevPtr
	}
}
func (chain *BlockChainTree) DeleteTree(tailId int) {
	if chain.tails[tailId].num == 0 {
		chain.numTail--
		chain.tails = chain.tails[:len(chain.tails)-1]
	} else {
		chain.tails[tailId] = chain.tails[tailId].prevPtr
	}
}
func (chain *BlockChainTree) AddBlockTree(block *BlockTree) int {
	//if block.blockBroad.tree == nil {
	//	fmt.Println("AddBlockTree err")
	//}
	tree, root := GenTree(block.blockBroad.Data)
	if !bytes.Equal(root, block.blockBroad.Root) {
		fmt.Println("AddBlockTree verify root error")
	}
	block.blockBroad.tree = tree
	block.blockBroad.Hash = block.CompHashTree()
	for i := 0; i < chain.numTail; i++ {
		if bytes.Equal(chain.tails[i].blockBroad.Hash, block.blockBroad.BlockHeader.PrevHash) {
			block.prevPtr = chain.tails[i]
			block.num = chain.tails[i].num + 1
			chain.tails[i] = block
			return i
		}
	}
	for i := 0; i < chain.numTail; i++ {
		tmp := chain.tails[i].prevPtr
		if tmp == nil {
			continue
		}
		num := tmp.num
		for j := 0; j < num; j++ {
			if bytes.Equal(tmp.blockBroad.Hash, block.blockBroad.BlockHeader.PrevHash) {
				block.prevPtr = tmp
				block.num = tmp.num + 1
				chain.numTail += 1
				chain.tails = append(chain.tails, block)
				return chain.numTail - 1
			} else {
				tmp = tmp.prevPtr
			}
		}
	}
	return 0
}

//func StartChainNew() *BlockChainNew {
//	startTrans := transaction{0, "start"}
//	Genesis := &BlockNew{BlockBroad{BlockHeader{0, []byte{}}, startTrans, []byte{}, 0},
//		0, nil} //这里的nonce是我随便写的
//	var tails = []*BlockNew{Genesis}
//	return &BlockChainNew{Genesis, tails, 1}
//}
//func StartChainArray() *BlockChainArray {
//	var startTrans []transaction
//	startTrans = append(startTrans, transaction{0, "start"})
//	Genesis := &BlockArray{BlockBroadArray{BlockHeader{0, []byte{}}, startTrans, []byte{}, 0},
//		0, nil} //这里的nonce是我随便写的
//	var tails = []*BlockArray{Genesis}
//	return &BlockChainArray{Genesis, tails, 1}
//}
//func StartChainTree() *BlockChainTree {
//	var startTrans []transaction
//	startTrans = append(startTrans, transaction{0, "start"})
//	Genesis := &BlockTree{BlockBroadTree{BlockHeader{0, []byte{}}, startTrans, []byte{}, 0, []byte{}},
//		0, nil} //这里的nonce是我随便写的
//	var tails = []*BlockTree{Genesis}
//	return &BlockChainTree{Genesis, tails, 1}
//}

func (chain *BlockChainNew) CheckNew() bool {
	for i := 0; i < chain.numTail; i++ {
		block := chain.tails[i]
		var prev *BlockNew
		for j := 0; j < block.num; j++ {
			prev = block.prevPtr
			if !bytes.Equal(prev.CompHashNew(), block.blockBroad.BlockHeader.PrevHash) {
				return false
			}
			block = prev
		}
	}
	return true
}
func (chain *BlockChainArray) CheckArray() bool {
	for i := 0; i < chain.numTail; i++ {
		block := chain.tails[i]
		var prev *BlockArray
		num := block.num
		for j := 0; j < num-1; j++ {
			prev = block.prevPtr
			if !bytes.Equal(prev.CompHashArray(), block.blockBroad.BlockHeader.PrevHash) {
				return false
			}
			block = prev
		}
		if num != 0 {
			prev = block.prevPtr
			if !bytes.Equal(prev.blockBroad.Hash, block.blockBroad.BlockHeader.PrevHash) {
				return false
			}
		}
	}
	return true
}
func (chain *BlockChainTree) CheckTree() bool {
	for i := 0; i < chain.numTail; i++ {
		block := chain.tails[i]
		var prev *BlockTree
		num := block.num
		for j := 0; j < num; j++ {
			prev = block.prevPtr
			if !bytes.Equal(prev.CompHashTree(), block.blockBroad.BlockHeader.PrevHash) {
				return false
			}
			//debug
			if block.blockBroad.tree == nil {
				fmt.Println("no tree", block.num, node)
			}
			//end
			block = prev
		}
	}
	return true
}
func (chain *BlockChainNew) LongestChain() []byte {
	// return the block in the longest chain
	longest := 0
	var idx int
	for i := 0; i < chain.numTail; i++ {
		if chain.tails[i].num > longest {
			longest = chain.tails[i].num
			idx = i
		}
	}
	return chain.tails[idx].blockBroad.Hash
}
func (chain *BlockChainArray) LongestChainArray() []byte {
	// return the block in the longest chain
	longest := 0
	idx := 0
	for i := 0; i < chain.numTail; i++ {
		if chain.tails[i].num > longest {
			longest = chain.tails[i].num
			idx = i
		}
	}
	return chain.tails[idx].blockBroad.Hash
}
func (chain *BlockChainTree) LongestChainTree() []byte {
	// return the block in the longest chain
	longest := 0
	var idx int
	for i := 0; i < chain.numTail; i++ {
		if chain.tails[i].num > longest {
			longest = chain.tails[i].num
			idx = i
		}
	}
	return chain.tails[idx].blockBroad.Hash
}
func (chain *BlockChainNew) LongestNum() int {
	longest := 0
	//var idx int
	for i := 0; i < chain.numTail; i++ {
		if chain.tails[i].num > longest {
			//idx = i
			longest = chain.tails[i].num
		}
	}
	return longest
}
func (chain *BlockChainArray) LongestNumArray() int {
	longest := 0
	//var idx int
	for i := 0; i < chain.numTail; i++ {
		if chain.tails[i].num > longest {
			//idx = i
			longest = chain.tails[i].num
		}
	}
	return longest
}
func (chain *BlockChainTree) LongestNumTree() int {
	longest := 0
	//var idx int
	for i := 0; i < chain.numTail; i++ {
		if chain.tails[i].num > longest {
			//idx = i
			longest = chain.tails[i].num
		}
	}
	return longest
}

// Acknowledge: https://jeiwan.net/posts/building-blockchain-in-go-part-1/
