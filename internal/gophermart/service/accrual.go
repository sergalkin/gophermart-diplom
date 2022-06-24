package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/sergalkin/gophermart-diplom/internal/gophermart/config"
	"github.com/sergalkin/gophermart-diplom/internal/gophermart/models"
	"github.com/sergalkin/gophermart-diplom/internal/gophermart/storage"
)

type Accrual interface {
	CalcAccrual()
}

var (
	AccrualCh chan models.OrderToProcess
	service   Orders
)

func InitAccrualService(s storage.Storage) {
	AccrualCh = make(chan models.OrderToProcess)
	service = NewOrdersService(s)
}

func CalcAccrual() {
	for data := range AccrualCh {
		go accrual(data.Login, data.Order)
	}
}

func accrual(login, order string) {
	status := "NEW"
	accrualModel := models.Accrual{}
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

		err = json.Unmarshal(b, &accrualModel)
		if err != nil {
			config.Logger().Error("error in unmarshaling")
		}
		if accrualModel.Status == "INVALID" || accrualModel.Status == "PROCESSED" {
			wait = true
		}

		if (accrualModel.Status == "REGISTERED" || accrualModel.Status == "PROCESSING") && accrualModel.Status != status {
			status = accrualModel.Status
			if err = service.Update(0, order, accrualModel.Status); err != nil {
				config.Logger().Error(err.Error())
			}

		}
		if !wait {
			time.Sleep(time.Second * 1)
		}
	}

	if accrualModel.Status == "INVALID" {
		accrualModel.Value = 0
	}

	if err = service.Update(accrualModel.Value, order, accrualModel.Status); err != nil {
		config.Logger().Error(err.Error())
		return
	}

	if err = service.UpdateBalance(login, accrualModel.Value); err != nil {
		config.Logger().Error(err.Error())
		return
	}
}
