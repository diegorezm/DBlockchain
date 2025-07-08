package blockchain

type BlockInsert struct {
	Index        uint64        `json:"index"`        // The of this block in the blockchain
	PrevHash     string        `json:"prevHash"`     // The hash of the previous block
	Transactions []Transaction `json:"transactions"` // The transactions inside of the block
}

// A block in the chain
type Block struct {
	BlockInsert `json:"block_insert"`
	Hash        string `json:"hash"`      // The current hash
	Nonce       uint64 `json:"nonce"`     // The cryptographic challenge
	Timestamp   int64  `json:"timestamp"` // The time the block was added to the chain
}

func NewBlock(data BlockInsert) *Block {
	return &Block{
		BlockInsert: data,
	}
}
