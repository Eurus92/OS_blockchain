package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"time"
)

var subsidy int = 50

type txInput struct {
	id     []byte
	Vout   int
	Sign   []byte
	PubKey []byte
}

type txOutput struct {
	value      int
	PubKeyHash []byte
}

type Transaction struct {
	id   []byte
	Vin  []txInput
	Vout []txOutput
}

func (tx *Transaction) Hash() []byte {
	//tx := tran.Transaction
	var encoder bytes.Buffer
	enc := gob.NewEncoder(&encoder)
	err := enc.Encode(tx)
	if err != nil {
		log.Panicln(err)
	}
	var hash [32]byte
	hash = sha256.Sum256(encoder.Bytes())
	return hash[:]
}

// 奖励给挖矿者（to）的credit
func NewCoinbaseTx(to, data string) *transaction {
	if data == "" {
		data = fmt.Sprintf("credit to miner %s", to)
	}
	txIn := txInput{[]byte{}, -1, nil, []byte(data)}
	txOut := txOutput{subsidy, nil}
	pubKeyHash := Base58Decode([]byte(to))
	if len(pubKeyHash) == 0 {
		fmt.Println(to)
		fmt.Println(pubKeyHash)
		return &transaction{time.Now().UnixNano(), Transaction{nil, []txInput{txIn}, []txOutput{txOut}}}
	}
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	txOut.PubKeyHash = pubKeyHash
	tx := Transaction{nil, []txInput{txIn}, []txOutput{txOut}}
	tx.id = tx.Hash()
	TX := transaction{time.Now().UnixNano(), tx}
	return &TX
}

func (trans *transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]transaction) {
	// Coinbase or not
	tx := trans.Transaction
	if trans.IsCoinbase() {
		return
	}
	for _, vin := range tx.Vin {
		if prevTXs[hex.EncodeToString(vin.id)].Transaction.id == nil {
			log.Panic("Err: Previous tx nor correct")
		}
	}
	txCopy := trans.TrimmedCopy().Transaction //修剪后的副本
	for inID, vin := range txCopy.Vin {
		prevTx := prevTXs[hex.EncodeToString(vin.id)]
		txCopy.Vin[inID].Sign = nil
		txCopy.Vin[inID].PubKey = prevTx.Transaction.Vout[vin.Vout].PubKeyHash
		txCopy.id = txCopy.Hash()
		txCopy.Vin[inID].PubKey = nil

		r, s, err := ecdsa.Sign(rand.Reader, &privKey, txCopy.id)
		if err != nil {
			log.Panic(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)

		tx.Vin[inID].Sign = signature
	}

}

func (trans *transaction) Verify(prevTXs map[string]transaction) bool {
	tx := trans.Transaction
	if trans.IsCoinbase() {
		return true
	}

	for _, vin := range tx.Vin {
		//遍历输入交易，如果发现输入交易引用的上一交易的ID不存在，则Panic
		if prevTXs[hex.EncodeToString(vin.id)].Transaction.id == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}
	txCopy := trans.TrimmedCopy().Transaction //修剪后的副本
	curve := elliptic.P256()                  //椭圆曲线实例

	for inID, vin := range tx.Vin {
		prevTX := prevTXs[hex.EncodeToString(vin.id)]
		txCopy.Vin[inID].Sign = nil //双重验证
		txCopy.Vin[inID].PubKey = prevTX.Transaction.Vout[vin.Vout].PubKeyHash
		txCopy.id = txCopy.Hash()
		txCopy.Vin[inID].PubKey = nil

		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Sign)
		r.SetBytes(vin.Sign[:(sigLen / 2)])
		s.SetBytes(vin.Sign[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.PubKey)
		x.SetBytes(vin.PubKey[:(keyLen / 2)])
		y.SetBytes(vin.PubKey[(keyLen / 2):])

		rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
		if ecdsa.Verify(&rawPubKey, txCopy.id, &r, &s) == false {
			return false
		}
	}
	return true
}
func (tx *transaction) TrimmedCopy() transaction {
	var inputs []txInput
	var outputs []txOutput

	for _, vin := range tx.Transaction.Vin {
		inputs = append(inputs, txInput{vin.id, vin.Vout, nil, nil})
	}

	for _, vout := range tx.Transaction.Vout {
		outputs = append(outputs, txOutput{vout.value, vout.PubKeyHash})
	}

	txCopy := transaction{tx.TimeStamp, Transaction{tx.Transaction.id, inputs, outputs}}
	return txCopy
}
func (out *txOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}

func (out *txOutput) Lock(addr []byte) {
	pubKeyHash := Base58Decode(addr)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	out.PubKeyHash = pubKeyHash
}
func (trans *transaction) IsCoinbase() bool {
	tx := trans.Transaction
	return len(tx.Vin) == 1 && len(tx.Vin[0].id) == 0 && tx.Vin[0].Vout == -1
}
func (in *txInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := HashPubKey(in.PubKey)
	return bytes.Compare(lockingHash, pubKeyHash) == 0
}
func NewTXOutput(value int, address string) *txOutput {
	tx := &txOutput{value, nil}
	tx.Lock([]byte(address))
	return tx
}
