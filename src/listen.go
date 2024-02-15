package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/rpc"
)

var pool = initPool()

// var chain = StartChainNew()
var breakPow = false

type ListenTrans struct{}

func (l *ListenTrans) Listen(transactionJson []byte, reply *string) error {
	lisTrans := Request{}
	err := json.Unmarshal(transactionJson, &lisTrans)
	if err != nil {
		fmt.Println(err)
	}
	pool.insert(lisTrans)
	*reply = "ACKT"
	return nil
}

//type ListenBlock struct {
//}
//
//func (l *ListenBlock) ListenB(blockJson []byte, reply *string) error {
//	lisBlock := BlockBroad{}
//	err := json.Unmarshal(blockJson, &lisBlock)
//	if err != nil {
//		fmt.Println(err)
//	}
//	chain.AddBlock(&BlockNew{lisBlock, 0, nil})
//	pool.delete(lisBlock.Data)
//	//fmt.Println(lisBlock.Data)
//	//TODO
//	//add more judge to breakPow
//	breakPow = true
//	*reply = "ACKB"
//	return nil
//}

type ListenBlockArray struct {
}
type ListenBlockTree struct {
}
type Message struct {
	Genesis    bool
	BlockBroad BlockBroadArray
	Requests   []Request
}
type MessageTree struct {
	Genesis    bool
	BlockBroad BlockBroadTree
	Requests   []Request
}

func (l *ListenBlockArray) ListenBA(blockJson []byte, reply *string) error {
	lisBlock := Message{}
	err := json.Unmarshal(blockJson, &lisBlock)
	if err != nil {
		fmt.Println(err)
	}
	if lisBlock.Genesis == true {
		NewBlockChain(&lisBlock.BlockBroad)
		*reply = "Receive Genesis"
		IsGen = false
	} else {
		chain.AddBlockArray(&BlockArray{lisBlock.BlockBroad, 0, nil})
		pool.deleteArray(lisBlock.Requests)
		*reply = "ACKB"
	}
	//fmt.Println(lisBlock.Data)
	//TODO
	//add more judge to breakPow
	breakPow = true
	return nil
}
func (l *ListenBlockTree) ListenBT(blockJson []byte, reply *string) error {
	lisBlock := MessageTree{}
	err := json.Unmarshal(blockJson, &lisBlock)
	if err != nil {
		fmt.Println(err)
	}
	if lisBlock.Genesis == true {
		NewBlockChainTree(&lisBlock.BlockBroad)
		*reply = "Receive Genesis"
		IsGen = false
	} else {
		chainTree.AddBlockTree(&BlockTree{lisBlock.BlockBroad, 0, nil})
		pool.deleteArray(lisBlock.Requests)
		*reply = "ACKB"
	}
	breakPow = true
	return nil
}

func Listen() {
	err := rpc.RegisterName("ListenTrans", new(ListenTrans))
	if err != nil {
		panic(err)
	}
	err = rpc.RegisterName("ListenBlockArray", new(ListenBlockArray))
	lis, err := net.Listen("tcp", ":6010")
	if err != nil {
		panic(err)
	}
	for {
		if false {
			break
		}
		con, err := lis.Accept()
		if err != nil {
			panic(err)
		}

		go rpc.ServeConn(con)
	}
}

func ListenTree() {
	err := rpc.RegisterName("ListenTrans", new(ListenTrans))
	if err != nil {
		panic(err)
	}
	err = rpc.RegisterName("ListenBlockTree", new(ListenBlockTree))
	lis, err := net.Listen("tcp", ":6010")
	if err != nil {
		panic(err)
	}
	for {
		if false {
			break
		}
		con, err := lis.Accept()
		if err != nil {
			panic(err)
		}

		go rpc.ServeConn(con)
	}
}
