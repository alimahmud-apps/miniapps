package controllers

import (
	"fmt"
	"miniapps/services"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type TransactionController struct {
	ewalletService services.EWalletService
}

func NewTransactionController(ewalletService services.EWalletService) *TransactionController {
	return &TransactionController{ewalletService: ewalletService}
}

type TransactionRequest struct {
	UserID int     `json:"user_id" validate:"required"`
	Amount float64 `json:"amount" validate:"min=1"`
}

type TransactionSuccessResponse struct {
	Status        string        `json:"status"`
	TransactionId int           `json:"transaction_id"`
	NewBalance    CustomFloat64 `json:"new_balance"`
}
type CustomFloat64 float64

// MarshalJSON overrides the default JSON marshaling for CustomFloat64
func (f CustomFloat64) MarshalJSON() ([]byte, error) {
	// Format the float64 value with two decimal places
	return []byte(fmt.Sprintf("%.2f", f)), nil
}

type TransactionErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// customErrorMessage custom func to mapping error validate needed
func customErrorMessage(e validator.FieldError, FieldName string) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", FieldName)
	default:
		return fmt.Sprintf("Invalid %s", FieldName)
	}
}

func (t *TransactionController) Credit(c echo.Context) error {
	request := TransactionRequest{}
	if err := c.Bind(&request); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, TransactionErrorResponse{
			Status:  "error",
			Message: "invalid body request",
		})
	}

	// validate field froms truct
	if err := c.Validate(&request); err != nil {
		c.Logger().Error(err.Error())
		ValidationErrors := err.(validator.ValidationErrors)
		var errorMessages []TransactionErrorResponse
		for _, ve := range ValidationErrors {
			message := customErrorMessage(ve, ve.Field())
			errorMessages = append(errorMessages, TransactionErrorResponse{
				Status:  "error",
				Message: message,
			})
		}
		return c.JSON(http.StatusBadRequest, errorMessages)

	}

	transactionID, balance, err := t.ewalletService.Credit(request.UserID, request.Amount)
	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusInternalServerError, TransactionErrorResponse{
			Status:  "error",
			Message: err.Error(),
		})
	}

	c.Logger().Info(TransactionSuccessResponse{
		Status:        "success",
		TransactionId: int(transactionID),
		NewBalance:    CustomFloat64(balance),
	})

	return c.JSON(http.StatusOK, TransactionSuccessResponse{
		Status:        "success",
		TransactionId: int(transactionID),
		NewBalance:    CustomFloat64(balance),
	})
}

func (t *TransactionController) Debit(c echo.Context) error {
	request := TransactionRequest{}
	if err := c.Bind(&request); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, TransactionErrorResponse{
			Status:  "error",
			Message: "invalid body request",
		})
	}

	// validate field froms truct
	if err := c.Validate(&request); err != nil {
		c.Logger().Error(err.Error())
		ValidationErrors := err.(validator.ValidationErrors)
		var errorMessages []TransactionErrorResponse
		for _, ve := range ValidationErrors {
			message := customErrorMessage(ve, ve.Field())
			errorMessages = append(errorMessages, TransactionErrorResponse{
				Status:  "error",
				Message: message,
			})
		}
		return c.JSON(http.StatusBadRequest, errorMessages)

	}

	transactionID, balance, err := t.ewalletService.Debit(request.UserID, request.Amount)

	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusInternalServerError, TransactionErrorResponse{
			Status:  "error",
			Message: err.Error(),
		})
	}

	c.Logger().Info(TransactionSuccessResponse{
		Status:        "success",
		TransactionId: int(transactionID),
		NewBalance:    CustomFloat64(balance),
	})

	return c.JSON(http.StatusOK, TransactionSuccessResponse{
		Status:        "success",
		TransactionId: int(transactionID),
		NewBalance:    CustomFloat64(balance),
	})
}
