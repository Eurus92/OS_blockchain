package main

import (
	"bytes"
	"crypto/sha256"
	"github.com/cbergoon/merkletree"
	"log"
	"strconv"
)

// TransData implements the Content interface provided by merkletree and represents the content stored in the tree.

// CalculateHash hashes the values of a TransData
func (t transaction) CalculateHash() ([]byte, error) {
	h := sha256.New()
	var tx []byte
	time := t.TimeStamp
	hash := t.Transaction.Hash()
	tx = bytes.Join([][]byte{tx, []byte(strconv.FormatInt(time, 10)), hash}, []byte{})

	if _, err := h.Write(tx); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

// Equals tests for equality of two Contents
func (t transaction) Equals(other merkletree.Content) (bool, error) {
	if t.TimeStamp != other.(transaction).TimeStamp {
		return false, nil
	}
	//flag := bytes.Equal(t.Transaction.Hash(), other.(transaction).Transaction.Hash())
	tx := other.(transaction).Transaction
	if bytes.Equal(t.Transaction.Hash(), tx.Hash()) {
		return false, nil
	}
	return true, nil
	//return t.trans == other.(TransData).trans, nil
}

func GenTree(txs []transaction) (*merkletree.MerkleTree, []byte) {
	var list []merkletree.Content
	for i := 0; i < len(txs); i++ {
		list = append(list, txs[i])
	}
	tree, err := merkletree.NewTree(list)
	if err != nil {
		log.Fatal(err)
	}
	mr := tree.MerkleRoot()
	//fmt.Println("root", mr)
	return tree, mr
}

func CheckMerkleRoot(tree *merkletree.MerkleTree) bool {
	vt, err := tree.VerifyTree()
	if err != nil {
		log.Fatal(err)
	}
	return vt
}
