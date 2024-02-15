package main

import "github.com/cbergoon/merkletree"

type BlockNew struct {
	blockBroad BlockBroad
	num        int
	prevPtr    *BlockNew
}

type BlockArray struct {
	blockBroad BlockBroadArray
	num        int
	prevPtr    *BlockArray
}
type BlockTree struct {
	blockBroad BlockBroadTree
	num        int
	prevPtr    *BlockTree
}

type BlockHeader struct {
	TimeStamp int64
	PrevHash  []byte
}

type BlockBroad struct {
	BlockHeader BlockHeader
	Data        transaction
	Hash        []byte
	Nonce       int64
}

type BlockBroadArray struct {
	BlockHeader BlockHeader
	Data        []transaction
	Hash        []byte
	Nonce       int64
}
type BlockBroadTree struct {
	BlockHeader BlockHeader
	Data        []transaction
	Hash        []byte
	Nonce       int64
	Root        []byte
	tree        *merkletree.MerkleTree
}
