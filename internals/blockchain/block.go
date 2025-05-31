package blockchain

type BlockInsert struct {
	Index    uint64        // The of this block in the blockchain
	Data     []Transaction // The transactions inside of the block
	PrevHash *string       // The hash of the previous block
}

type Block struct {
	BlockInsert
	Hash      string // The current hash
	Nonce     uint64 // The cryptographic challenge
	Timestamp int64  // The time the block was added to the chain
}

func NewBlock(data BlockInsert) *Block {
	return &Block{
		BlockInsert: data,
	}
}
