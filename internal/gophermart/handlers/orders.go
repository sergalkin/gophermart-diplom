package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sergalkin/gophermart-diplom/internal/gophermart/config"
	"github.com/sergalkin/gophermart-diplom/internal/gophermart/models"
	"io"
	"net/http"
	"strconv"
	"time"

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
	status := "NEW"
	accrual := models.Accrual{}
	url := fmt.Sprintf("%s/api/orders/%s", config.AccrualAddress(), order)

	client := http.Client{}
	request, err := http.NewRequest(http.MethodGet, url, bytes.NewReader([]byte{}))

	if err != nil {
		config.Logger().Error("error in creating request to accrual")
		return
	}

	var wait bool
	for !wait {
		response, errRes := client.Do(request)
		if errRes != nil {
			config.Logger().Error("error in request to accrual")
			return
		}

		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			if response.StatusCode == http.StatusNoContent {
				config.Logger().Error("order not registered")
				return
			}
			if response.StatusCode == http.StatusTooManyRequests {
				retryAfter := response.Header.Get("Retry-After")
				t, errT := strconv.Atoi(retryAfter)
				if errT != nil {
					config.Logger().Error("string to int casting failed. Fallback to default value")
					t = 60
				}
				time.Sleep(time.Duration(t) * time.Second)
				continue
			} else {
				config.Logger().Error("retrying....")
				time.Sleep(1 * time.Second)
				continue
			}
		}

		b, errB := io.ReadAll(response.Body)
		if errB != nil {
			config.Logger().Error(errB.Error())
		}

		err = json.Unmarshal(b, &accrual)
		if err != nil {
			config.Logger().Error("error in unmarshaling")
		}
		if accrual.Status == "INVALID" || accrual.Status == "PROCESSED" {
			wait = true
		}

		if (accrual.Status == "REGISTERED" || accrual.Status == "PROCESSING") && accrual.Status != status {
			status = accrual.Status
			if err = o.service.Update(0, order, accrual.Status); err != nil {
				config.Logger().Error(err.Error())
			}

		}
		if !wait {
			time.Sleep(time.Second * 1)
		}
	}

	if accrual.Status == "INVALID" {
		accrual.Value = 0
	}

	if err = o.service.Update(accrual.Value, order, accrual.Status); err != nil {
		config.Logger().Error(err.Error())
		return
	}

	if err = o.service.UpdateBalance(login, accrual.Value); err != nil {
		config.Logger().Error(err.Error())
		return
	}
}
