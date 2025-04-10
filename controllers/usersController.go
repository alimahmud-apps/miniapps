package controllers

import (
	"miniapps/services"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type UsersController struct {
	ewalletService services.EWalletService
}

func NewUsersController(ewalletService services.EWalletService) *UsersController {
	return &UsersController{ewalletService: ewalletService}
}

type UsersRequest struct {
	Username string `json:"username" validate:"required"`
}

type UsersSuccessResponse struct {
	Status   string `json:"status"`
	UsersId  int    `json:"user_id"`
	Username string `json:"username"`
}

type UsersErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (t *UsersController) Create(c echo.Context) error {
	request := UsersRequest{}
	if err := c.Bind(&request); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, UsersErrorResponse{
			Status:  "error",
			Message: "invalid body request",
		})
	}

	// validate field froms truct
	if err := c.Validate(&request); err != nil {
		c.Logger().Error(err.Error())
		ValidationErrors := err.(validator.ValidationErrors)
		var errorMessages []UsersErrorResponse
		for _, ve := range ValidationErrors {
			message := customErrorMessage(ve, ve.Field())
			errorMessages = append(errorMessages, UsersErrorResponse{
				Status:  "error",
				Message: message,
			})
		}
		return c.JSON(http.StatusBadRequest, errorMessages)

	}

	userId, errCreate := t.ewalletService.UserCreate(request.Username)
	if errCreate != nil {
		c.Logger().Error(errCreate.Error())
		return c.JSON(http.StatusInternalServerError, UsersErrorResponse{
			Status:  "error",
			Message: errCreate.Error(),
		})
	}

	c.Logger().Info(UsersSuccessResponse{
		Status:   "success",
		UsersId:  userId,
		Username: request.Username,
	})

	return c.JSON(http.StatusOK, UsersSuccessResponse{
		Status:   "success",
		UsersId:  userId,
		Username: request.Username,
	})
}

func (t *UsersController) Retrieve(c echo.Context) error {
	id := c.Param("id")
	idx, _ := strconv.Atoi(id)
	user, err := t.ewalletService.GetUsersByID(idx)
	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusInternalServerError, UsersErrorResponse{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, user)
}
