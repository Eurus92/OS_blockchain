package main

import (
	"encoding/json"
	"fmt"
	"net/rpc"
)

func (block *BlockBroad) broadcast() {
	for i := 0; i < len(addrList); i++ {
		go block.call(addrList[i])
	}
}

func (block *BlockBroad) call(addr string) {
	//fmt.Println(block)
	mar, err := json.Marshal(&block)
	//fmt.Println(string(mar))
	if err != nil {
		fmt.Println("生成json错误")
	} else {
		//首先是通过rpc.Dial拨号RPC服务，然后通过client.Call调用具体的RPC方法
		cli, err := rpc.Dial("tcp", addr+":6010") //这里call不同的服务器
		if err != nil {
			panic(err)
		}

		var reply string
		//在调用client.Call时，
		//第一个参数是用点号链接的RPC服务名字和方法名字，
		//第二和第三个参数分别我们定义RPC方法的两个参数。
		err = cli.Call("ListenBlock.ListenB", mar, &reply)
		if err != nil {
			panic(err)
		}
		fmt.Println(reply)
	}
}
