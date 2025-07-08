package blockchain

import (
	"testing"

	"github.com/diegorezm/DBlockchain/internals/utils"
)

func TestBlockchain_Genesis(t *testing.T) {
	blockchain := NewBlockchain("")
	if len(blockchain.Chain) == 0 {
		t.Error("The genesis block was not created.\n")
	}
}

func TestBlockchain_AppendBlock(t *testing.T) {
	blockchain := NewBlockchain("")
	if len(blockchain.Chain) == 0 {
		t.Error("The genesis block was not created.\n")
	}
	err := blockchain.AppendBlock()
	if err != nil {
		t.Errorf("The block was not appended to the chain.\n%v\n", err)
	}
}
func TestBlockchain_Transaction(t *testing.T) {
	blockchain := NewBlockchain("")

	// 1. Generate a key pair (sender)
	priv, err := utils.GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate keypair: %v", err)
	}

	// 2. Encode public key as address
	keypair, err := utils.EncodeKeyPair(priv)

	if err != nil {
		t.Fatalf("Failed to encode public key: %v", err)
	}

	// 3. Manually create a UTXO by adding a funding transaction in a new block
	fundingTx := &Transaction{
		Id:     "funding-tx-1",
		TxIns:  []TxIn{}, // coinbase or genesis, no input
		TxOuts: []TxOut{{Address: keypair.PublicKey, Amount: 5.0}},
	}
	block := NewBlock(BlockInsert{
		Index:    1,
		PrevHash: blockchain.Chain[len(blockchain.Chain)-1].Hash,
	})
	block.Transactions = []Transaction{*fundingTx}
	block.Hash = "mockedhash"
	blockchain.Chain = append(blockchain.Chain, *block)

	// 4. Build transaction input using UTXO from fundingTx
	txInput := TransactionInput{
		TxIns: []TxIn{{
			TxOutId:    "funding-tx-1",
			TxOutIndex: 0,
			Signature:  "", // to be signed
		}},
		TxOuts: []TxOut{
			{Address: "bob-address", Amount: 3.0},
			{Address: keypair.PublicKey, Amount: 2.0},
		},
	}

	// 5. Sign and create transaction
	newTX, err := NewSignedTransaction(txInput, priv)
	if err != nil {
		t.Fatalf("Failed to create signed transaction: %v", err)
	}

	// 6. Validate transaction
	if err := blockchain.ValidateTransaction(newTX); err != nil {
		t.Fatalf("Transaction failed validation: %v", err)
	}

	// 7. Append transaction to mempool (or next block)
	blockchain.AppendTransaction(newTX)

	// 8. Mine the block
	err = blockchain.AppendBlock()
	if err != nil {
		t.Fatalf("Failed to mine block with transaction: %v", err)
	}

	// 9. Check if transaction was included
	latest := blockchain.GetLastBlock()
	found := false

	for _, tx := range latest.Transactions {
		if tx.Id == newTX.Id {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("Transaction was not included in the mined block")
	}
}
