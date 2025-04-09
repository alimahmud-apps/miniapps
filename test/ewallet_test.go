package tests

import (
	"miniapps/config"
	"miniapps/logs"
	"miniapps/repositories"
	"miniapps/services"
	"sync"
	"testing"
)

func TestCreditTransactions(t *testing.T) {
	// Initialize database
	config.InitDB()

	userRepo := repositories.NewUserRepository()
	transactionRepo := repositories.NewTransactionRepository()
	ewalletService := services.NewEWalletService(userRepo, transactionRepo)
	logger := logs.NewLogger()
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(5)
		go func() {
			defer wg.Done()
			for j := 0; j < 2; j++ {
				transactionID, balance, err := ewalletService.Credit(123, 100.00) //userid , amount
				logger.Info.Println(transactionID, balance)
				if err != nil {
					logger.Error.Println(err)
				}
			}
		}()
	}

	wg.Wait()
}
