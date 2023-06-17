package transaction

import (
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"godatasource"
	"gorm.io/gorm"
)

func ExecuteTransaction(f func() error) error {
	return ExecuteTransactionWithDB(godatasource.DefaultConnection(), f)
}

func ExecuteTransactionWithDB(db *gorm.DB, f func() error) error {
	tx := db.Begin()
	txId, _ := uuid.NewUUID()
	log.Debug().Any("txId", txId).Msg("Start transaction")
	err := f()
	if err != nil {
		tx.Rollback()
		log.Warn().Any("txId", txId).Msg("Got error. Rollback transaction")
	} else {
		tx.Commit()
		log.Debug().Any("txId", txId).Msg("Transaction committed")
	}
	return err
}
