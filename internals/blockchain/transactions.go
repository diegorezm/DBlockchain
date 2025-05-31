package blockchain

type TransactionInsert struct {
	From   string
	To     string
	Amount float32
}

type Transaction struct {
	TransactionInsert
}

func NewTransaction(data TransactionInsert) *Transaction {
	return &Transaction{
		TransactionInsert: data,
	}
}
