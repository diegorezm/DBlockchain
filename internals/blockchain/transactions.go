package blockchain

type TransactionInsert struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float32 `json:"amount"`
}

type TransactionBulkRequest struct {
	Transactions []TransactionInsert `json:"transactions"`
}

type Transaction struct {
	TransactionInsert
}

func NewTransaction(data TransactionInsert) *Transaction {
	return &Transaction{
		TransactionInsert: data,
	}
}
