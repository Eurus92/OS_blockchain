package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"time"
)

var minedBlock int
var IsGen = true
var node string

// var wallets, _ = NewWallets()
var wallets Wallets

func main() {
	start := time.Now()
	node = os.Args[1]
	filename1 := os.Args[2]
	Tree := os.Args[3]
	isTree, _ := strconv.Atoi(Tree)
	wallets.Wallets = make(map[string]*Wallet)
	err := wallets.LoadFromFile()
	id, _ := strconv.Atoi(node)
	var address string
	err = wallets.LoadFromFile()
	fmt.Println("load", err)
	addr := wallets.GetAddresses()
	fmt.Println(addr)
	for i := 0; i < 5; i++ {
		for j := 0; j < 4-i; j++ {
			if compare(addr[j], addr[j+1]) {
				continue
			} else {
				tmp := addr[j]
				addr[j] = addr[j+1]
				addr[j+1] = tmp
			}

		}
	}
	address = addr[id-16]
	go trans(address)
	flag := false
	init := false
	var timeV int
	if isTree == 0 {
		go Listen()
		Gen := NewGenesisBlock()
		if Gen != nil {
			fmt.Println("Find Genesis")
			minedBlock++
			message := Message{true, *Gen, []Request{}}
			message.broadcast()
			breakPow = false
		}
		//fmt.Println(Gen)
		for {
			if IsGen == false {
				break
			}
		}
		//fmt.Println(chain.tails[0])
		for {
			if pool.len > 0 {
				init = true
				hash := chain.LongestChainArray()
				block, requests := NewBlockArray(hash, chain, false)
				//fmt.Println(block)
				pow := block.NewPoW()
				res := pow.proof(false)
				if res {
					//fmt.Println("mined")
					minedBlock++
					message := Message{false, *block, requests}
					message.broadcast()
					breakPow = false
				} else {
					//fmt.Println("not mined")
				}
				//fmt.Println(pool.len)
			} else if init == true {
				time.Sleep(time.Second * 5)
				if pool.len == 0 {
					flag = true
				}
			}
			if flag {
				break
			}
		}
		timeV = VerifyArray(chain)
		duration := time.Since(start).Milliseconds()
		err0 := os.MkdirAll("result", os.ModePerm)
		if err0 != nil {
			fmt.Println("mkdir err")
			return
		}
		file, err := os.OpenFile(filename1, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0644)
		if err != nil {
			panic(err)
		} else {
			_, err = file.WriteString(strconv.FormatFloat(float64(duration), 'E', -1, 64) + "\n")
			_, err = file.WriteString(strconv.Itoa(timeV) + "\n")
			_, err = file.WriteString(strconv.Itoa(minedBlock) + "\n")
			_, err = file.WriteString(strconv.Itoa(chain.LongestNumArray()) + "\n")
			err = file.Close()
			if err != nil {
				panic(err)
			}
			fmt.Println("node" + node + "finished")
		}
		fmt.Println("node", node, ":", chain)
		fmt.Println("node", node, "check", chain.CheckArray())
	} else {
		go ListenTree()
		Gen := NewGenesisBlockTree()
		if Gen != nil {
			fmt.Println("Find Genesis")
			minedBlock++
			message := MessageTree{true, *Gen, []Request{}}
			message.broadcast()
			breakPow = false
		}
		for {
			if IsGen == false {
				break
			}
		}
		for {
			if pool.len > 0 {
				init = true
				hash := chainTree.LongestChainTree()
				block, requests := NewBlockTree(hash, chainTree, false)
				//fmt.Println(block)
				pow := block.NewPoWTree()
				res := pow.proofTree(false)
				if res {
					//fmt.Println("mined")
					minedBlock++
					message := MessageTree{false, *block, requests}
					message.broadcast()
					breakPow = false
				} else {
					//fmt.Println("not mined")
				}
				//fmt.Println(pool.len)
			} else if init == true {
				time.Sleep(time.Second * 5)
				if pool.len == 0 {
					flag = true
				}
			}
			if flag {
				break
			}
		}
		timeV = VerifyTree(chainTree)
		duration := time.Since(start).Milliseconds()
		err0 := os.MkdirAll("result", os.ModePerm)
		if err0 != nil {
			fmt.Println("mkdir err")
			return
		}
		file, err := os.OpenFile(filename1, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0644)
		if err != nil {
			panic(err)
		} else {
			_, err = file.WriteString(strconv.FormatFloat(float64(duration), 'E', -1, 64) + "\n")
			_, err = file.WriteString(strconv.Itoa(timeV) + "\n")
			_, err = file.WriteString(strconv.Itoa(minedBlock) + "\n")
			_, err = file.WriteString(strconv.Itoa(chainTree.LongestNumTree()) + "\n")
			err = file.Close()
			if err != nil {
				panic(err)
			}
			fmt.Println("node" + node + "finished")
		}
		fmt.Println("node", node, ":", chainTree)
		fmt.Println("node", node, "check", chainTree.CheckTree())
	}
}

func compare(s1 string, s2 string) bool {
	d1, _ := strconv.Atoi(s1)
	d2, _ := strconv.Atoi(s2)
	if d1 < d2 {
		return true
	}
	return false
}
func VerifyArray(chain *BlockChainArray) int {
	longest := 0
	var idx = 0
	for i := 0; i < chain.numTail; i++ {
		if chain.tails[i].num > longest {
			longest = chain.tails[i].num
			idx = i
		}
	}
	block := chain.tails[idx]
	Start := time.Now()
	num := block.num
	var prev *BlockArray
	for i := 0; i < num; i++ {
		prev = block.prevPtr
		if !bytes.Equal(block.CompHashArray(), block.blockBroad.Hash) {
			fmt.Println("block hash error", block.num)
		}
		//fmt.Println(block.CompHashArray(), node, "line1compute")
		//fmt.Println(block.blockBroad.Hash, node, "line1hash")
		if i == num-1 {
			if !bytes.Equal(prev.blockBroad.Hash, block.blockBroad.BlockHeader.PrevHash) {
				fmt.Println("prev block hash error", block.num)
			}
		} else {
			if !bytes.Equal(prev.CompHashArray(), block.blockBroad.BlockHeader.PrevHash) {
				fmt.Println("prev block hash error", block.num)
				//fmt.Println(prev.CompHashArray(), node, "compute prev hash")
				//fmt.Println(prev.blockBroad.Hash, block.num, node, "prev hash")
				//fmt.Println(block.blockBroad.BlockHeader.PrevHash, node, "PrevHash")
			}
		}
		block = prev
	}
	duration := time.Since(Start).Nanoseconds()
	//fmt.Println(duration)
	return int(duration)
}
func VerifyTree(chain *BlockChainTree) int {
	longest := 0
	var idx = 0
	for i := 0; i < chain.numTail; i++ {
		if chain.tails[i].num > longest {
			longest = chain.tails[i].num
			idx = i
		}
	}
	block := chain.tails[idx]
	Start := time.Now()
	var prev *BlockTree
	num := block.num
	for i := 0; i < chain.tails[idx].num-1; i++ {
		//fmt.Println(block.blockBroad.tree)
		//vt, err := block.blockBroad.tree.VerifyTree()
		//if err != nil {
		//	log.Fatal(err)
		//	//fmt.Println("root", err)
		//	//fmt.Println(vt)
		//}
		//if !vt {
		//	fmt.Println("block root error")
		//}
		prev = block.prevPtr
		if !bytes.Equal(block.CompHashTree(), block.blockBroad.Hash) {
			fmt.Println("block hash error", block.num)
		}
		if i == num-1 {
			if !bytes.Equal(prev.blockBroad.Hash, block.blockBroad.BlockHeader.PrevHash) {
				fmt.Println("prev block hash error", block.num)
			}
		} else {
			if !bytes.Equal(prev.CompHashTree(), block.blockBroad.BlockHeader.PrevHash) {
				fmt.Println("prev block hash error", block.num)
			}
		}
		block = prev

	}
	duration := time.Since(Start).Nanoseconds()
	//fmt.Println(duration)
	return int(duration)
}
