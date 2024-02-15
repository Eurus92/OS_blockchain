package main

type transaction struct {
	TimeStamp   int64
	Transaction Transaction
}

type Request struct {
	TimeStamp int64
	From      string
	To        string
	Amount    int
}

type transactionNode struct {
	Transaction Request
	prev        *transactionNode
	post        *transactionNode
}

type transactionPool struct {
	head *transactionNode
	tail *transactionNode
	len  int
}

func initPool() *transactionPool {
	head := new(transactionNode)
	tail := new(transactionNode)
	return &transactionPool{head, tail, 0}
}

func (pool *transactionPool) insert(transaction Request) {
	newNode := new(transactionNode)
	newNode.Transaction = transaction

	if pool.len == 0 {
		pool.head = newNode
		pool.tail = newNode
		pool.len++
		return
	} else {
		idx := pool.len
		ptr := new(transactionNode)
		ptr = pool.tail
		//fmt.Println(ptr)
		for idx > 0 {
			if ptr.Transaction.TimeStamp <= transaction.TimeStamp {
				if idx == pool.len {
					ptr.post = newNode
					newNode.prev = ptr
					newNode.post = nil
					pool.tail = newNode
					pool.len++
					return
				} else {
					//post := new(transactionNode)
					//post = ptr.post
					//newNode.prev = ptr
					//newNode.post = post
					//post.prev = newNode
					//ptr.post = newNode

					newNode.prev = ptr
					newNode.post = ptr.post
					ptr.post = newNode
					ptr.post.prev = newNode
					pool.len++
					return
				}
			} else {
				idx--
				if idx != 0 {
					ptr = ptr.prev
				}
			}
		}
		if idx == 0 {
			head := new(transactionNode)
			head = pool.head
			head.prev = newNode
			newNode.post = head
			pool.head = newNode
			pool.len++
			return
		}
	}
}

func (pool *transactionPool) deleteArray(list []Request) {
	listlen := len(list)
	for idx := 0; idx < listlen; idx++ {
		pool.delete(list[idx])
	}
}

func (pool *transactionPool) delete(transaction Request) {
	idx := pool.len
	ptr := pool.head
	for {
		if idx == 0 {
			break
		}
		if ptr.Transaction.TimeStamp == transaction.TimeStamp && ptr.Transaction.From == transaction.From && ptr.Transaction.To == transaction.To && ptr.Transaction.Amount == transaction.Amount {
			if pool.len == 1 {
				pool.head = nil
				pool.tail = nil
				pool.len--
				return
			} else if idx == pool.len {
				//newHead := pool.head.post
				pool.head.post.prev = nil
				pool.head = pool.head.post
				pool.len--
				return
			} else if idx == 1 {
				//newTail := pool.tail.prev
				pool.tail.prev.post = nil
				pool.tail = pool.tail.prev
				pool.len--
				return
			} else {
				ptr.prev.post = ptr.post
				ptr.post.prev = ptr.prev
				pool.len--
				return
			}
		} else {
			idx--
			ptr = ptr.post
		}
	}
}

func (pool *transactionPool) getHead(num int) []Request {
	var array []Request
	ptr := pool.head
	for i := 0; i < num; i++ {
		if i == pool.len {
			break
		} else {
			array = append(array, ptr.Transaction)
			ptr = ptr.post
		}
	}
	return array
}
