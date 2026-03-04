package repository

import (
	"belajar-go/models"
	"context"

	"gorm.io/gorm"
)

type txKey string

const TxContextKey txKey = "tx_db"

type txManager struct {
	db *gorm.DB
}

func NewTransactionManager(db *gorm.DB) models.TransactionManager {
	return &txManager{db: db}
}

func (t *txManager) WithTransaction(ctx context.Context, fn models.RunInTransactionFunc) error {
	// Memulai transaksi GORM
	tx := t.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Sisipkan object transaksi (tx) ke dalam context yang baru
	ctxWithTx := context.WithValue(ctx, TxContextKey, tx)

	// Jalankan fungsi (service logic) menggunakan context yang sudah ada transaksinya
	err := fn(ctxWithTx)

	// Jika terjadi error dari fungsi, Rollback
	if err != nil {
		tx.Rollback()
		return err
	}

	// Jika sukses, Commit
	return tx.Commit().Error
}

// --- HELPER UNTUK REPOSITORY ---
// Fungsi ini akan mengekstrak DB dari context jika ada transaksi berjalan.
// Jika tidak ada transaksi, ia akan menggunakan DB biasa.
func GetDB(ctx context.Context, defaultDB *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(TxContextKey).(*gorm.DB); ok {
		return tx
	}
	return defaultDB
}
