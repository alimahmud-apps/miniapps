package repositories

import (
	"database/sql"
	"log"
	"miniapps/models"
)

type TransactionRepository interface {
	CreateTransaction(transaction *models.Transaction, tx *sql.Tx) error
}

type transactionRepo struct{}

func NewTransactionRepository() TransactionRepository {
	return &transactionRepo{}
}

func (r *transactionRepo) CreateTransaction(transaction *models.Transaction, tx *sql.Tx) error {
	query := `INSERT INTO transactions (user_id, amount, type, created_at) VALUES ($1, $2, $3, $4) RETURNING id`
	var lastInsertId int64
	err := tx.QueryRow(query, transaction.UserID, transaction.Amount, transaction.Type, transaction.CreatedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("Failed to create transaction:", err)
		return err
	}
	transaction.ID = lastInsertId
	return nil
}
