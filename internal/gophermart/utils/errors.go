package utils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var (
	ErrUserAlreadyExists              = errors.New("user already exists")
	ErrUserNotFound                   = errors.New("user not found")
	ErrUserPasswordMissMatch          = errors.New("wrong username or password")
	ErrLuhnValidation                 = errors.New("invalid order number")
	ErrOrderAlreadyCreatedByOtherUser = errors.New("order already created by another user")
	ErrOrderUniqueViolation           = errors.New("order already created")
	ErrOrdersNotFound                 = errors.New("orders not found")
	ErrNegativeBalance                = errors.New("not enough balance to proceed")
	ErrEmptyRow                       = errors.New("no record found")
)

type responseError struct {
	Message       string         `json:"message"`
	BindingErrors []bindingError `json:"errors"`
}

type bindingError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func FormResponseBindingError(msg string, casted validator.ValidationErrors) *responseError {
	return &responseError{
		Message:       msg,
		BindingErrors: prepareBindingErrors(casted),
	}
}

func prepareBindingErrors(casted validator.ValidationErrors) []bindingError {
	var err []bindingError

	for _, v := range casted {
		err = append(err, bindingError{
			Field:   strings.ToLower(v.Field()),
			Message: fmt.Sprintf("validation failed on the '%s' tag", v.Tag()),
		})
	}

	return err
}

func JSONErrorResponse(msg string) gin.H {
	return gin.H{
		"message": msg,
		"errors":  []interface{}{},
	}
}

func NoPayloadProvided() gin.H {
	return JSONErrorResponse("Wrong payload provided. Expected to be valid Json")
}
