package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/diegorezm/DBlockchain/internals/blockchain"
	"github.com/diegorezm/DBlockchain/internals/set"
)

type Blockchain struct {
	chain        []Block
	transactions []Transaction
	nodes        set.Set[Node]
	difficulty   uint32
}

func NewBlockchain() *Blockchain {
	chain := make([]Block, 1)
	chain[0] = *generateGenesis()

	return &Blockchain{
		chain:        chain,
		transactions: make([]Transaction, 0),
		nodes:        *set.NewSet[Node](),
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

	b.nodes.Add(Node{address: address})
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

func (b *Blockchain) replaceChain() {
	replacementChain := []Block{}
	maxChainLen := len(b.chain)

	b.nodes.ForEach(func(key Node) {
		url := key.address.String()
		chain, err := GetBlockchainFromNode(url)

		if err != nil {
			fmt.Printf(err.Error())
			return
		}

		if len(chain) > maxChainLen {
			maxChainLen = len(chain)
			replacementChain = chain
		}
	})

	if len(replacementChain) > 0 && isChainValid(replacementChain) {
		b.chain = replacementChain
	}
}

func GetBlockchainFromNode(reqUrl string) ([]Block, error) {
	res, err := http.Get(reqUrl)
	if err != nil {
		return []Block{}, err
	}

	defer res.Body.Close()

	var chain []Block

	dec := gob.NewDecoder(res.Body)
	err = dec.Decode(&chain)

	if err != nil {
		return []Block{}, err
	}

	return chain, nil
}

func (b *Blockchain) getNodes() []Node {
	return b.nodes.ToSlice()
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
