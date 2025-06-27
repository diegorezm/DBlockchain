package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	webutils "github.com/diegorezm/DBlockchain/internals/web_utils"
)

type Blockchain struct {
	chain        []Block
	transactions []Transaction
	difficulty   uint32
	serverUrl    string
	// Current Node indicates the IP of the current client in the blockchain
	currentNode string
}

func NewBlockchain(currentNode string) *Blockchain {
	chain := make([]Block, 1)
	chain[0] = *generateGenesis()

	return &Blockchain{
		chain:        chain,
		transactions: make([]Transaction, 0),
		difficulty:   2,
		serverUrl:    "http://localhost:4040",
		currentNode:  currentNode,
	}
}

func (b *Blockchain) AppendBlock() error {
	lastBlock := b.GetLastBlock()
	if lastBlock == nil {
		return fmt.Errorf("Something went wrong while getting the last block.")
	}

	newBlockInsert := BlockInsert{
		PrevHash: lastBlock.Hash,
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

func (b *Blockchain) GetLastBlock() *Block {
	return &b.chain[len(b.chain)-1]
}

func (b *Blockchain) GetChain() []Block {
	return b.chain
}

func (b *Blockchain) GetTransactions() []Transaction {
	return b.transactions
}

func (b *Blockchain) replaceChain() (bool, error) {
	replacementChain := []Block{}
	maxChainLen := len(b.chain)
	nodes, err := getConnectedNodes(b.serverUrl)

	if err != nil {
		return false, err
	}

	for _, address := range nodes {
		chain, err := getBlockchainFromNode(address)

		if err != nil {
			return false, fmt.Errorf("ERROR: something went wrong while connecting to node %s. %s\n", address, err.Error())
		}

		//fmt.Printf("Node: %s. Chain: %v\n", address, chain)

		if len(chain) > maxChainLen {
			maxChainLen = len(chain)
			replacementChain = chain
		}
	}

	if len(replacementChain) > 0 && isChainValid(replacementChain) {
		b.chain = replacementChain
		return true, nil
	}
	return false, nil
}

func getConnectedNodes(serverUrl string) ([]string, error) {
	reqUrl := fmt.Sprintf("%s/nodes", serverUrl)

	res, err := http.Get(reqUrl)

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		defer res.Body.Close()
		return nil, fmt.Errorf("received non-OK status %d from %s: %s", res.StatusCode, reqUrl, string(bodyBytes))
	}

	if err != nil {
		defer res.Body.Close()
		return []string{}, err
	}

	response, err := webutils.ParseJSON[webutils.JSONResponse[[]string]](res.Body)

	if err != nil {
		return []string{}, err
	}
	return response.Data, nil
}

func getBlockchainFromNode(address string) ([]Block, error) {
	reqUrl := fmt.Sprintf("%s/chain", address)
	fmt.Printf("Sending request to: %s\n", reqUrl)

	res, err := http.Get(reqUrl)

	if res.StatusCode != http.StatusOK {
		defer res.Body.Close()
		bodyBytes, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("received non-OK status %d from %s: %s", res.StatusCode, reqUrl, string(bodyBytes))
	}

	if err != nil {
		defer res.Body.Close()
		return []Block{}, err
	}

	response, err := webutils.ParseJSON[webutils.JSONResponse[[]Block]](res.Body)

	if err != nil {
		return []Block{}, err
	}
	return response.Data, nil
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
		prevBlock = currentBlock
	}
	return true
}

func isBlockPairValid(prevBlock, nextBlock *Block) error {
	nextBlockHash := hashBlock(nextBlock)
	if nextBlockHash != nextBlock.Hash {
		return fmt.Errorf(
			"ERROR: Block hash mismatch for Block #%d. Stored hash: %s, Re-computed hash: %s. Block details: %v\n",
			nextBlock.Index, nextBlock.Hash, nextBlockHash, nextBlock,
		)
	}

	if prevBlock.Hash != nextBlock.PrevHash {
		return fmt.Errorf(
			"ERROR: Previous hash mismatch for Block #%d. Expected PrevHash (from prev block #%d): %s, Got PrevHash (in current block): %s\n",
			nextBlock.Index, prevBlock.Index, prevBlock.Hash, nextBlock.PrevHash,
		)
	}

	if prevBlock.Index != nextBlock.Index-1 {
		return fmt.Errorf(
			"ERROR: Block index sequence mismatch. Previous block #%d, Next block #%d. Expected Next block to be #%d\n",
			prevBlock.Index, nextBlock.Index, prevBlock.Index+1,
		)
	}

	return nil
}

func generateGenesis() *Block {
	newBlockInsert := BlockInsert{
		PrevHash: "",
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
	PrevHash  string
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
		panic(fmt.Sprintf("Failed to gob encode block header for hashing: %v\n", err))
	}
	sum := sha256.Sum256(buf.Bytes())
	hashSlice := sum[:]
	return hex.EncodeToString(hashSlice)
}
