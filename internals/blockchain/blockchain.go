package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"
	"time"
)

type Blockchain struct {
	chain        []Block
	transactions []Transaction
	nodes        []Node
	difficulty   uint32
}

func NewBlockchain() *Blockchain {
	chain := make([]Block, 1)
	chain[0] = *generateGenesis()

	return &Blockchain{
		chain:        chain,
		transactions: make([]Transaction, 0),
		nodes:        make([]Node, 0),
		difficulty:   2,
	}
}

func (b *Blockchain) AppendBlock() error {
	lastBlock := b.GetLastBlock()
	if lastBlock == nil {
		return fmt.Errorf("Something went wrong while getting the last block.")
	}

	newBlockInsert := BlockInsert{
		PrevHash: &lastBlock.Hash,
		Index:    lastBlock.Index + 1,
		Data:     b.transactions,
	}

	blockToMine := NewBlock(newBlockInsert)
	blockToMine.Timestamp = time.Now().Unix()

	newBlock := b.mine(blockToMine)

	err := isBlockPairValid(lastBlock, newBlock)

	if err != nil {
		return err
	}

	b.chain = append(b.chain, *newBlock)
	b.transactions = make([]Transaction, 0)
	return nil
}

func (b *Blockchain) AppendTransaction(transactionInsert TransactionInsert) {
	transaction := NewTransaction(transactionInsert)
	b.transactions = append(b.transactions, *transaction)
}

func (b *Blockchain) AppendNode(addr string) {
	address, err := url.Parse(addr)

	if err != nil {
		fmt.Printf("ERROR: %s", err)
		return
	}

	b.nodes = append(b.nodes, Node{address: address})
}

func (b *Blockchain) GetLastBlock() *Block {
	return &b.chain[len(b.chain)-1]
}

func (b *Blockchain) GetChain() []Block {
	return b.chain
}

func (b *Blockchain) GetTransactions() []Transaction {
	return b.transactions
}

func (b *Blockchain) getNodes() []Node {
	return b.nodes
}

// This function mines the chain untils it finds a valid block, when this block is found
// the mining stops and the valid block is returned.
func (b *Blockchain) mine(blockToMine *Block) *Block {
	var nonce uint64 = 0

	for {
		blockToMine.Nonce = nonce
		computedHash := hashBlock(blockToMine)

		if strings.HasPrefix(computedHash, strings.Repeat("0", int(b.difficulty))) {
			blockToMine.Hash = computedHash
			break
		} else {
			nonce++
		}
		if nonce > 1_000_000_000 {
			panic("Mining failed to find a block within 1 billion attempts (difficulty too high or logic error)\n")
		}
	}
	return blockToMine
}

func isChainValid(chain []Block) bool {
	prevBlock := chain[0]
	for i := 1; i < len(chain)-1; i++ {
		currentBlock := chain[i]
		err := isBlockPairValid(&prevBlock, &currentBlock)
		if err != nil {
			fmt.Print(err)
			return false
		}
	}
	return true
}

func isBlockPairValid(prevBlock, nextBlock *Block) error {
	nextBlockHash := hashBlock(nextBlock)

	if nextBlockHash != nextBlock.Hash {
		return fmt.Errorf("ERORR: Hashes does not match. %s != %s\n", nextBlock.Hash, nextBlockHash)
	}

	if prevBlock.Hash != *nextBlock.PrevHash {
		return fmt.Errorf("ERORR: NextBlock hash does not match PrevBlock hash\n .")
	}

	if prevBlock.Index != nextBlock.Index-1 {
		return fmt.Errorf("ERORR: Indexes don't match.\n")
	}

	return nil
}

func generateGenesis() *Block {
	newBlockInsert := BlockInsert{
		PrevHash: nil,
		Index:    0,
		Data:     make([]Transaction, 0),
	}
	block := NewBlock(newBlockInsert)
	hash := hashBlock(block)
	block.Hash = hash
	return block
}

type blockHeader struct {
	Index     uint64
	Timestamp int64
	Data      []Transaction
	PrevHash  *string
	Nonce     uint64
}

func hashBlock(b *Block) string {
	header := blockHeader{
		Index:     b.Index,
		Timestamp: b.Timestamp,
		Data:      b.Data,
		PrevHash:  b.PrevHash,
		Nonce:     b.Nonce,
	}

	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)
	err := enc.Encode(header)
	if err != nil {
		panic(fmt.Sprintf("Failed to gob encode block header for hashing: %v", err))
	}
	sum := sha256.Sum256(buf.Bytes())
	hashSlice := sum[:]
	return hex.EncodeToString(hashSlice)
}
