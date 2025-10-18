package blockchain

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/diegorezm/DBlockchain/internals/utils"
)

func writeMetric(file_name, name string, iteration int, diff uint32, value time.Duration) {
	f := fmt.Sprintf("%s.csv", file_name)
	file, err := os.OpenFile(f,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	info, _ := file.Stat()
	if info.Size() == 0 {
		_ = writer.Write([]string{"metric_name", "iteration", "difficulty", "duration_ns"})
	}

	durationMs := float64(value.Microseconds()) / 1000.0

	record := []string{
		name,
		strconv.Itoa(iteration),
		strconv.Itoa(int(diff)),
		fmt.Sprintf("%.3f", durationMs),
	}
	_ = writer.Write(record)
}

// 💥 Mede o tempo médio de validação de uma transação
func BenchmarkMetric_TransactionValidation(b *testing.B) {
	bc := NewBlockchain("")

	priv, _ := utils.GenerateKeyPair()
	kp, _ := utils.EncodeKeyPair(priv)

	// Prepara um UTXO inicial
	tx := &Transaction{
		Id:     "funding-tx",
		TxOuts: []TxOut{{Address: kp.PublicKey, Amount: 10.0}},
	}
	block := NewBlock(BlockInsert{
		Index:    1,
		PrevHash: bc.Chain[len(bc.Chain)-1].Hash,
	})
	block.Transactions = []Transaction{*tx}
	bc.Chain = append(bc.Chain, *block)

	for i := 0; b.Loop(); i++ {
		txInput := TransactionInput{
			TxIns: []TxIn{{TxOutId: "funding-tx", TxOutIndex: 0}},
			TxOuts: []TxOut{
				{Address: "bob", Amount: 5.0},
				{Address: kp.PublicKey, Amount: 5.0},
			},
		}

		newTx, _ := NewSignedTransaction(txInput, priv)
		start := time.Now()
		_ = bc.ValidateTransaction(newTx)
		elapsed := time.Since(start)

		writeMetric("ts_validation_metrics", "TransactionValidation", i, bc.Difficulty, elapsed)
		b.Logf("Validação da transação levou: %v", elapsed)
	}
}

func BenchmarkMetric_MiningSpeedByDifficulty(b *testing.B) {
	for d := uint32(1); d <= 6; d++ {
		bc := NewBlockchain("")
		bc.Difficulty = d

		for i := range 10 {
			start := time.Now()
			_ = bc.AppendBlock()
			elapsed := time.Since(start)

			writeMetric("mining_speed_metrics", "MiningSpeed", i, d, elapsed)
			b.Logf("Dificuldade: %d — Tempo: %v", d, elapsed)
		}
	}
}
