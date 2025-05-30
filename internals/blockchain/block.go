package blockchain

type BlockInsert struct {
	Index    int64         // The of this block in the blockchain
	Data     []Transaction // The transactions inside of the block
	hash     string        // The current hash
	PrevHash *string       // The hash of the previous block
}

type Block struct {
	BlockInsert
	Nonce     float64
	Timestamp float64
}

func NewBlock(data BlockInsert) *Block {
	return &Block{
		BlockInsert: data,
	}
}
