package blockchain

type BlockInsert struct {
	Index    uint64        `json:"index"`    // The of this block in the blockchain
	Data     []Transaction `json:"data"`     // The transactions inside of the block
	PrevHash *string       `json:"prevHash"` // The hash of the previous block
}

type Block struct {
	BlockInsert
	Hash      string `json:"hash"`      // The current hash
	Nonce     uint64 `json:"nonce"`     // The cryptographic challenge
	Timestamp int64  `json:"timestamp"` // The time the block was added to the chain
}

func NewBlock(data BlockInsert) *Block {
	return &Block{
		BlockInsert: data,
	}
}
