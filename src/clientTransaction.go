package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net/rpc"
	"strconv"
	"time"
)

var addrList = []string{"10.6.0.16", "10.6.0.17", "10.6.0.18", "10.6.0.19", "10.6.0.20"}
var numTX = 5

func trans(addr string) {
	//file, err := os.Open(filename)
	//if err != nil {
	//	panic(err)
	//}

	// 创建 Reader
	//r := bufio.NewReader(file)
	//var f = filename[7:9]
	//id, _ := strconv.Atoi(i)
	for i := 0; i < numTX; i++ {
		time.Sleep(time.Second * 5)
		//transactionBytes, err := r.ReadBytes('\n')
		//transaction := strings.TrimSpace(string(transactionBytes))
		//if err != nil && err != io.EOF {
		//	panic(err)
		//}
		//if err == io.EOF {
		//	break
		//}

		create(addr).broadcast()
	}
}

func create(addr string) *Request {
	timestamp := time.Now().UnixNano()
	//from := wallets.GetAddresses()[id]
	toId, _ := rand.Int(rand.Reader, big.NewInt(5))
	Id := toId.String()
	ID, _ := strconv.Atoi(Id)
	to := wallets.GetAddresses()[ID]
	//rand.Seed(time.Now().Unix())
	amount, _ := rand.Int(rand.Reader, big.NewInt(50))
	num := amount.String()
	amt, _ := strconv.Atoi(num)
	node := &Request{timestamp, addr, to, amt}
	return node
}

func (transaction *Request) broadcast() {
	for i := 0; i < len(addrList); i++ {
		go transaction.call(addrList[i])
	}
}

func (transaction *Request) call(addr string) {
	mar, err := json.Marshal(&transaction)
	if err != nil {
		fmt.Println("生成json错误")
	} else {
		cli, err := rpc.Dial("tcp", addr+":6010") //这里call不同的服务器
		if err != nil {
			panic(err)
		}

		var reply string
		err = cli.Call("ListenTrans.Listen", mar, &reply)
		if err != nil {
			panic(err)
		}
		fmt.Println(reply)
	}
}
