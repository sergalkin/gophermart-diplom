package handlers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sergalkin/gophermart-diplom/internal/gophermart/utils"
	"net/http"

	"github.com/sergalkin/gophermart-diplom/internal/gophermart/service"
)

type balanceHandler struct {
	service service.Balance
}

type withdrawData struct {
	Order string  `json:"order" binding:"required"`
	Sum   float32 `json:"sum" binding:"required"`
}

func NewBalanceHandler(s service.Balance) *balanceHandler {
	return &balanceHandler{service: s}
}

func (b *balanceHandler) Check() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		login, err := utils.GetUserFromCookie(ctx)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, utils.JSONErrorResponse("unauthorized"))
			return
		}

		balance, errB := b.service.Get(login)
		if errB != nil {
			fmt.Println(errB)
			ctx.JSON(http.StatusInternalServerError, utils.JSONErrorResponse("Internal server error"))
			return
		}

		ctx.JSON(http.StatusOK, balance)
	}
}

func (b *balanceHandler) Withdraw() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		login, err := utils.GetUserFromCookie(ctx)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, utils.JSONErrorResponse("unauthorized"))
			return
		}

		data := &withdrawData{}
		if errB := ctx.ShouldBindJSON(&data); errB != nil {
			casted, ok := errB.(validator.ValidationErrors)
			if !ok {
				ctx.JSON(http.StatusBadRequest, utils.NoPayloadProvided())
				return
			}

			ctx.JSON(http.StatusBadRequest, utils.FormResponseBindingError("error in provided data", casted))
			return
		}

		if err = b.service.Withdraw(data.Sum, login, data.Order); err != nil {
			if errors.Is(err, utils.ErrLuhnValidation) {
				ctx.JSON(http.StatusUnprocessableEntity, utils.JSONErrorResponse(utils.ErrLuhnValidation.Error()))
				return
			}
			if errors.Is(err, utils.ErrNegativeBalance) {
				ctx.JSON(http.StatusPaymentRequired, utils.JSONErrorResponse("not enough funds"))
				return
			}

			ctx.JSON(http.StatusInternalServerError, utils.JSONErrorResponse("internal server error"))
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Successfully withdrawn",
		})
	}
}

func (b *balanceHandler) Withdrawals() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		login, err := utils.GetUserFromCookie(ctx)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, utils.JSONErrorResponse("unauthorized"))
			return
		}

		withdrawals, errW := b.service.Withdrawals(login)
		if errW != nil {
			if errors.Is(errW, utils.ErrEmptyRow) {
				ctx.JSON(http.StatusNoContent, utils.JSONErrorResponse(errW.Error()))
				return
			}

			ctx.JSON(http.StatusInternalServerError, utils.JSONErrorResponse("internal server error"))
			return
		}

		ctx.JSON(http.StatusOK, withdrawals)
	}
}
