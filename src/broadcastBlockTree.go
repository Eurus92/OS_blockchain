package main

import (
	"encoding/json"
	"fmt"
	"net/rpc"
)

func (message *MessageTree) broadcast() {
	for i := 0; i < len(addrList); i++ {
		go message.callTree(addrList[i])
	}
}

func (message *MessageTree) callTree(addr string) {
	//fmt.Println(block)
	mar, err := json.Marshal(&message)
	//fmt.Println(string(mar))
	if err != nil {
		fmt.Println("生成json错误")
	} else {
		cli, err := rpc.Dial("tcp", addr+":6010") //这里call不同的服务器
		if err != nil {
			panic(err)
		}

		var reply string
		err = cli.Call("ListenBlockTree.ListenBT", mar, &reply)
		if err != nil {
			panic(err)
		}
		fmt.Println(reply)
	}
}
