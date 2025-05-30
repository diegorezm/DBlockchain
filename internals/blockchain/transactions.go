package blockchain

type TransactionInsert struct {
	from   string
	to     string
	amount float32
}

type Transaction struct {
	TransactionInsert
}

func NewTransaction(data TransactionInsert) *Transaction {
	return &Transaction{
		TransactionInsert: data,
	}
}
