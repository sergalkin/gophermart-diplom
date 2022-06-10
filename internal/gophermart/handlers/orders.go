package handlers

import (
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/sergalkin/gophermart-diplom/internal/gophermart/service"
	"github.com/sergalkin/gophermart-diplom/internal/gophermart/utils"
)

type ordersHandler struct {
	service service.Orders
}

func NewOrdersHandler(s service.Orders) *ordersHandler {
	return &ordersHandler{service: s}
}

func (o *ordersHandler) Post() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		login, err := utils.GetUserFromCookie(ctx)
		if err != nil {
			ctx.Data(http.StatusUnauthorized, "text/plain", []byte("unauthorized"))
			return
		}

		body, errBody := io.ReadAll(ctx.Request.Body)
		if errBody != nil {
			ctx.Data(http.StatusInternalServerError, "text/plain", []byte("error in reading body"))
			return
		}
		order := string(body)

		err = o.service.Create(login, order)
		if err != nil {
			if errors.Is(err, utils.ErrLuhnValidation) {
				ctx.Data(http.StatusUnprocessableEntity, "text/plain", []byte(utils.ErrLuhnValidation.Error()))
				return
			}
			if errors.Is(err, utils.ErrOrderUniqueViolation) {
				ctx.Data(http.StatusOK, "text/plain", []byte(utils.ErrOrderUniqueViolation.Error()))
				return
			}
			if errors.Is(err, utils.ErrOrderAlreadyCreatedByOtherUser) {
				ctx.Data(http.StatusConflict, "text/plain", []byte(utils.ErrOrderAlreadyCreatedByOtherUser.Error()))
				return
			}

			ctx.Data(http.StatusInternalServerError, "text/plain", []byte("Internal server error"))
			return
		}

		ctx.Data(http.StatusAccepted, "text/plain", []byte("order accepted"))

		go o.accrual(login, order)
	}
}

func (o *ordersHandler) Get() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		login, err := utils.GetUserFromCookie(ctx)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, utils.JSONErrorResponse("unauthorized"))
			return
		}

		orders, errOrd := o.service.Get(login)
		if errOrd != nil {
			if errors.Is(errOrd, utils.ErrOrdersNotFound) {
				ctx.JSON(http.StatusNoContent, utils.JSONErrorResponse(errOrd.Error()))
				return
			}

			ctx.JSON(http.StatusInternalServerError, utils.JSONErrorResponse(errOrd.Error()))
			return
		}

		ctx.JSON(http.StatusOK, orders)
	}
}

func (o *ordersHandler) accrual(login, order string) {
	//todo
}
