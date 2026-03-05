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

	ctxWithTx := context.WithValue(ctx, TxContextKey, tx)

	err := fn(ctxWithTx)

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func GetDB(ctx context.Context, defaultDB *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(TxContextKey).(*gorm.DB); ok {
		return tx
	}
	return defaultDB
}
