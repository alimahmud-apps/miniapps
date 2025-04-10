package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"miniapps/config"
	"miniapps/models"
	"strings"

	"github.com/lib/pq"
)

type UserRepository interface {
	UpdateBalance(userID int, amount float64, tx *sql.Tx) error
	GetBalance(userID int) (float64, error)
	GetUsersByID(userID int) (models.User, error)
	CreateUser(username string, tx *sql.Tx) (int, error)

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

		if err == sql.ErrNoRows {
			return 0, errors.New("user not found")
		}
		log.Println("Failed to get balance:", err)
		return 0, err
	}

	return balance, nil
}

func (r *userRepo) CreateUser(username string, tx *sql.Tx) (int, error) {
	query := `INSERT INTO users (username) VALUES ($1) RETURNING id`
	var lastInsertId int
	err := tx.QueryRow(query, username).Scan(&lastInsertId)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if strings.Contains(pqErr.Message, "users_username_key") {
				return 0, errors.New("username already exists")
			}
		}
		log.Println("Failed to create user:", err)
		return 0, err
	}
	fmt.Println("lastInsertId", lastInsertId)
	return lastInsertId, nil
}

func (r *userRepo) GetUsersByID(userID int) (models.User, error) {
	query := `select id,username,balance,created_at from users WHERE id = $1`

	var user models.User
	err := config.DB.Get(&user, query, userID)
	if err != nil {

		if err == sql.ErrNoRows {
			return models.User{}, errors.New("user not found")
		}
		log.Println("Failed to get balance:", err)
		return models.User{}, err
	}

	return user, nil
}
