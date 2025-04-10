package main

import (
	"miniapps/config"
	"miniapps/controllers"
	"miniapps/repositories"
	"miniapps/services"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func main() {
	// Initialize database
	config.InitDB()
	// Create instances of repositories and services
	userRepo := repositories.NewUserRepository()
	transactionRepo := repositories.NewTransactionRepository()
	ewalletService := services.NewEWalletService(userRepo, transactionRepo)

	// Create instances of controllers
	transactionController := controllers.NewTransactionController(ewalletService)
	usersController := controllers.NewUsersController(ewalletService)

	// Setup Echo
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	e.Logger.SetLevel(2)
	// Routes
	e.POST("/api/transactions/credit", transactionController.Credit)
	e.POST("/api/transactions/debit", transactionController.Debit)
	e.POST("/api/users", usersController.Create)
	e.GET("/api/users/:id", usersController.Retrieve)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

// Custom validator for Echo
type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
