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

	"github.com/diegorezm/DBlockchain/internals/utils"
	webutils "github.com/diegorezm/DBlockchain/internals/web_utils"
)

type Blockchain struct {
	Chain               []Block       `json:"chain"`
	TransactionsMempool []Transaction `json:"transactions_mempool"` // the mempool: pending txs
	Difficulty          uint32        `json:"difficulty"`
	ServerUrl           string        `json:"server_url"`
	CurrentNode         string        `json:"current_node"`
}

func NewBlockchain(currentNode string) *Blockchain {
	chain := make([]Block, 1)
	chain[0] = *generateGenesis()

	return &Blockchain{
		Chain:               chain,
		Difficulty:          2,
		ServerUrl:           "http://localhost:4040",
		CurrentNode:         currentNode,
		TransactionsMempool: make([]Transaction, 0),
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
		// TODO: Maybe i should add a way for the user to choose the transactions he wants to add
		Transactions: b.TransactionsMempool,
	}

	blockToMine := NewBlock(newBlockInsert)
	blockToMine.Timestamp = time.Now().Unix()

	newBlock := b.mine(blockToMine)

	err := isBlockPairValid(lastBlock, newBlock)

	if err != nil {
		return err
	}

	b.Chain = append(b.Chain, *newBlock)
	b.TransactionsMempool = make([]Transaction, 0)
	return nil
}

func (b *Blockchain) GetLastBlock() *Block {
	return &b.Chain[len(b.Chain)-1]
}

func (b *Blockchain) GetChain() []Block {
	return b.Chain
}

func (b *Blockchain) AppendTransaction(tx *Transaction) error {
	if err := b.ValidateTransaction(tx); err != nil {
		return err
	}

	b.TransactionsMempool = append(b.TransactionsMempool, *tx)
	return nil
}

// Get all unspent Transactions
func (bc *Blockchain) GetUTXOPool() []UTXO {
	utxos := make(map[string]UTXO)

	for _, block := range bc.Chain {
		for _, tx := range block.Transactions {
			for i, txOut := range tx.TxOuts {
				key := fmt.Sprintf("%s_%d", tx.Id, i)
				utxos[key] = UTXO{
					TxId:   tx.Id,
					Index:  int64(i),
					Output: txOut,
				}
			}

			for _, txIn := range tx.TxIns {
				key := fmt.Sprintf("%s_%d", txIn.TxOutId, txIn.TxOutIndex)
				delete(utxos, key)
			}
		}
	}

	result := make([]UTXO, 0, len(utxos))

	for _, u := range utxos {
		result = append(result, u)
	}

	return result
}

// Get unspent transactions by address
func (b *Blockchain) GetUTXPoolByAddress(address string) []UTXO {
	utxos := b.GetUTXOPool()
	result := make([]UTXO, 0)

	for _, u := range utxos {
		if u.Output.Address == address {
			result = append(result, u)
		}
	}

	return result
}

func (b *Blockchain) ValidateTransaction(tx *Transaction) error {
	utxos := b.GetUTXOPool()

	totalInput := float64(0)
	totalOutput := float64(0)

	for _, txIn := range tx.TxIns {
		// 1. Find matching UTXO
		utxoKey := fmt.Sprintf("%s_%d", txIn.TxOutId, txIn.TxOutIndex)

		var utxo *UTXO
		for _, u := range utxos {
			if u.TxId == txIn.TxOutId && u.Index == txIn.TxOutIndex {
				utxo = &u
				break
			}
		}
		if utxo == nil {
			return fmt.Errorf("invalid TxIn: no matching UTXO for %s", utxoKey)
		}

		// 2. Verify the signature
		pubKey, err := utils.DecodePublicKey(utxo.Output.Address)
		if err != nil {
			return fmt.Errorf("invalid public key for address %s", utxo.Output.Address)
		}

		if !VerifyTransactionSignature(tx.Id, txIn.Signature, pubKey) {
			return fmt.Errorf("invalid signature for input %s", utxoKey)
		}

		totalInput += utxo.Output.Amount
	}

	// 3. Validate outputs
	for _, txOut := range tx.TxOuts {
		totalOutput += txOut.Amount
	}

	// 4. Inputs must be ≥ outputs
	if !tx.IsSystem && totalInput < totalOutput {
		return fmt.Errorf("input (%.2f) < output (%.2f)", totalInput, totalOutput)
	}
	return nil
}

func (b *Blockchain) ReplaceChain() (bool, error) {
	replacementChain := []Block{}
	maxChainLen := len(b.Chain)
	nodes, err := getConnectedNodes(b.ServerUrl)

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

	if len(replacementChain) > 0 && IsChainValid(replacementChain) {
		b.Chain = replacementChain
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
	reqUrl := fmt.Sprintf("%s/api/chain", address)
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

		if strings.HasPrefix(computedHash, strings.Repeat("0", int(b.Difficulty))) {
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

func IsChainValid(chain []Block) bool {
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
	}
	block := NewBlock(newBlockInsert)
	hash := hashBlock(block)
	block.Hash = hash
	return block
}

type blockHeader struct {
	Index        uint64
	Timestamp    int64
	Transactions []Transaction
	PrevHash     string
	Nonce        uint64
}

func hashBlock(b *Block) string {
	header := blockHeader{
		Index:        b.Index,
		Timestamp:    b.Timestamp,
		Transactions: b.Transactions,
		PrevHash:     b.PrevHash,
		Nonce:        b.Nonce,
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
