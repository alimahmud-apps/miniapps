package services

import (
	"errors"
	"miniapps/models"
	"miniapps/repositories"
	"sync"
	"time"
)

type EWalletService interface {
	Credit(userID int, amount float64) (int64, float64, error)
	Debit(userID int, amount float64) (int64, float64, error)
	UserCreate(uername string) (int, error)
	GetUsersByID(id int) (models.User, error)
}

type eWalletService struct {
	userRepo        repositories.UserRepository
	transactionRepo repositories.TransactionRepository
	mu              sync.Mutex
}

func NewEWalletService(userRepo repositories.UserRepository, transactionRepo repositories.TransactionRepository) EWalletService {
	return &eWalletService{
		userRepo:        userRepo,
		transactionRepo: transactionRepo,
	}
}

func (s *eWalletService) Credit(userID int, amount float64) (int64, float64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	tx, err := s.userRepo.BeginTransaction()
	if err != nil {
		return 0, 0, err
	}

	balance, err := s.userRepo.GetBalance(userID)
	if err != nil {
		return 0, 0, err
	}

	// Update balance
	err = s.userRepo.UpdateBalance(userID, amount, tx)
	if err != nil {
		errRb := s.userRepo.RollbackTransaction(tx)
		if errRb != nil {
			return 0, 0, errRb
		}
		return 0, 0, err
	}

	// Create transaction
	transaction := &models.Transaction{
		UserID:    userID,
		Amount:    amount,
		Type:      "credit",
		CreatedAt: time.Now(),
	}
	err = s.transactionRepo.CreateTransaction(transaction, tx)
	if err != nil {
		errRb := s.userRepo.RollbackTransaction(tx)
		if errRb != nil {
			return 0, 0, errRb
		}
		return 0, 0, err
	}
	err = s.userRepo.CommitTransaction(tx)
	if err != nil {
		return 0, 0, err
	}
	balance = balance + amount
	return transaction.ID, balance, err
}

func (s *eWalletService) Debit(userID int, amount float64) (int64, float64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	balance, err := s.userRepo.GetBalance(userID)
	if err != nil {
		return 0, 0, err
	}

	if balance < amount {
		return 0, 0, errors.New("insufficient funds")
	}

	tx, err := s.userRepo.BeginTransaction()
	if err != nil {
		return 0, 0, err
	}

	// Update balance
	err = s.userRepo.UpdateBalance(userID, -amount, tx)
	if err != nil {
		errRb := s.userRepo.RollbackTransaction(tx)
		if errRb != nil {
			return 0, 0, errRb
		}
		return 0, 0, err
	}

	// Create transaction
	transaction := &models.Transaction{
		UserID:    userID,
		Amount:    amount,
		Type:      "debit",
		CreatedAt: time.Now(),
	}

	err = s.transactionRepo.CreateTransaction(transaction, tx)
	if err != nil {
		errRb := s.userRepo.RollbackTransaction(tx)
		if errRb != nil {
			return 0, 0, errRb
		}
		return 0, 0, err
	}

	err = s.userRepo.CommitTransaction(tx)
	if err != nil {
		return 0, 0, err
	}

	balance = balance - amount
	return transaction.ID, balance, err
}

func (s *eWalletService) UserCreate(username string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	tx, err := s.userRepo.BeginTransaction()
	if err != nil {
		return 0, err
	}
	userId, err := s.userRepo.CreateUser(username, tx)
	if err != nil {
		errRb := s.userRepo.RollbackTransaction(tx)
		if errRb != nil {
			return 0, errRb
		}
		return 0, err
	}
	err = s.userRepo.CommitTransaction(tx)
	if err != nil {
		return 0, err
	}
	return userId, nil
}

func (s *eWalletService) GetUsersByID(id int) (models.User, error) {
	user, err := s.userRepo.GetUsersByID(id)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
