package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"math/big"
)

type TxOut struct {
	Address string  `json:"address"`
	Amount  float64 `json:"amount"`
}

type TxIn struct {
	TxOutId    string `json:"tx_out_id"`
	TxOutIndex int64  `json:"tx_out_index"`
	Signature  string `json:"signature"`
}

type UTXO struct {
	TxId   string `json:"tx_id"`
	Index  int64  `json:"index"`
	Output TxOut  `json:"output"`
}

type Transaction struct {
	Id       string  `json:"id"`
	TxIns    []TxIn  `json:"tx_ins"`
	TxOuts   []TxOut `json:"tx_outs"`
	IsSystem bool    `json:"is_system"`
}

type TransactionInput struct {
	TxIns    []TxIn  `json:"tx_ins"`
	TxOuts   []TxOut `json:"tx_outs"`
	IsSystem bool    `json:"is_system"`
}

func generateTransactionId(txIns []TxIn, txOuts []TxOut) (string, error) {
	var txInContentBuf, txOutContentBuf bytes.Buffer

	txInEnc := gob.NewEncoder(&txInContentBuf)
	err := txInEnc.Encode(txIns)

	if err != nil {
		return "", err
	}

	txOutEnc := gob.NewEncoder(&txOutContentBuf)
	err = txOutEnc.Encode(txOuts)

	if err != nil {
		return "", err
	}

	var combined bytes.Buffer
	combined.Write(txInContentBuf.Bytes())
	combined.Write(txOutContentBuf.Bytes())

	sum := sha256.Sum256(combined.Bytes())
	return fmt.Sprintf("%x", sum[:]), nil
}

func SignTransactionId(txId string, privKey *ecdsa.PrivateKey) (string, error) {
	hash := sha256.Sum256([]byte(txId))
	r, s, err := ecdsa.Sign(rand.Reader, privKey, hash[:])

	if err != nil {
		return "", err
	}

	signature := append(r.Bytes(), s.Bytes()...)

	return base64.StdEncoding.EncodeToString(signature), nil
}

func VerifyTransactionSignature(txId string, signature string, pubKey *ecdsa.PublicKey) bool {
	hash := sha256.Sum256([]byte(txId))
	sigBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false
	}

	r := new(big.Int).SetBytes(sigBytes[:len(sigBytes)/2])
	s := new(big.Int).SetBytes(sigBytes[len(sigBytes)/2:])
	return ecdsa.Verify(pubKey, hash[:], r, s)
}

func NewTransaction(input TransactionInput) (*Transaction, error) {
	id, err := generateTransactionId(input.TxIns, input.TxOuts)

	if err != nil {
		return nil, err
	}

	return &Transaction{
		Id:       id,
		TxIns:    input.TxIns,
		TxOuts:   input.TxOuts,
		IsSystem: input.IsSystem,
	}, nil
}

func NewSignedTransaction(input TransactionInput, privKey *ecdsa.PrivateKey) (*Transaction, error) {
	id, err := generateTransactionId(input.TxIns, input.TxOuts)

	if err != nil {
		return nil, err
	}

	for i := range input.TxIns {
		sig, _ := SignTransactionId(id, privKey)
		input.TxIns[i].Signature = sig
	}

	return &Transaction{
		Id:       id,
		TxIns:    input.TxIns,
		TxOuts:   input.TxOuts,
		IsSystem: input.IsSystem,
	}, nil
}
