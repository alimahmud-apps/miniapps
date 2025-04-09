package repositories

import (
	"database/sql"
	"log"
	"miniapps/config"
)

type UserRepository interface {
	UpdateBalance(userID int, amount float64, tx *sql.Tx) error
	GetBalance(userID int) (float64, error)
	BeginTransaction() (*sql.Tx, error)
	CommitTransaction(*sql.Tx) error
	RollbackTransaction(*sql.Tx) error
}

type userRepo struct{}

func NewUserRepository() UserRepository {
	return &userRepo{}
}

func (r *userRepo) BeginTransaction() (*sql.Tx, error) {
	return config.DB.Begin()
}
func (r *userRepo) CommitTransaction(tx *sql.Tx) error {
	return tx.Commit()
}
func (r *userRepo) RollbackTransaction(tx *sql.Tx) error {
	return tx.Rollback()
}

func (r *userRepo) UpdateBalance(userID int, amount float64, tx *sql.Tx) error {
	query := `UPDATE users SET balance = balance + $1 WHERE id = $2`
	_, err := tx.Exec(query, amount, userID)
	if err != nil {
		log.Println("Failed to update balance:", err)
		return err
	}
	return nil
}

func (r *userRepo) GetBalance(userID int) (float64, error) {
	query := `select balance from users WHERE id = $1`
	row := config.DB.QueryRow(query, userID)
	balance := float64(0)
	err := row.Scan(&balance)
	if err != nil {
		log.Println("Failed to get balance:", err)
		return 0, err
	}

	return balance, nil
}
