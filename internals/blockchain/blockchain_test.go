package blockchain

import (
	"testing"
)

func Test_NewBlockchain(t *testing.T) {
	blockchain := NewBlockchain()
	chain := blockchain.GetChain()

	if len(chain) != 1 {
		t.Errorf("Something went wrong while creating the genesis block.")
	}

	if chain[0].Index != 0 {
		t.Errorf("The index of the genesis block is wrong.")
	}

	if chain[0].PrevHash != nil {
		t.Errorf("The PrevHash of the genesis should be nil.")
	}
}

func Test_AppendBlock(t *testing.T) {
	blockchain := NewBlockchain()
	chain := blockchain.GetChain()

	if len(chain) != 1 {
		t.Errorf("Something went wrong while creating the genesis block.")
	}

	transactionInsert := TransactionInsert{
		From:   "a",
		To:     "b",
		Amount: 1,
	}
	blockchain.AppendTransaction(transactionInsert)
	blockchain.AppendTransaction(transactionInsert)

	transactions := blockchain.GetTransactions()
	if len(transactions) != 2 {
		t.Errorf("Something went wrong while adding transactions.\n")
	}

	blockchain.AppendBlock()
	chain = blockchain.GetChain()

	if len(chain) != 2 {
		t.Errorf("Something went wrong while mining the block.\n")
	}
}

func Test_AppendNode(t *testing.T) {
	blockchain := NewBlockchain()
	chain := blockchain.GetChain()

	if len(chain) != 1 {
		t.Errorf("Something went wrong while creating the genesis block.")
	}

	blockchain.AppendNode("localhost:3000")
	nodes := blockchain.getNodes()

	if len(nodes) != 1 {
		t.Errorf("Something went wrong while adding a node.\n")
	}
}

func Test_IsChainValid(t *testing.T) {
	blockchain := NewBlockchain()
	chain := blockchain.GetChain()

	if len(chain) != 1 {
		t.Errorf("Something went wrong while creating the genesis block.")
	}

	transactionInsert := TransactionInsert{
		From:   "a",
		To:     "b",
		Amount: 1,
	}
	blockchain.AppendTransaction(transactionInsert)
	blockchain.AppendTransaction(transactionInsert)

	if !isChainValid(blockchain.GetChain()) {
		t.Error("Something went wrong while appending blocks to the chain.")
	}
}
