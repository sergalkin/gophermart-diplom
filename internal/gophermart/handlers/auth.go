package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/sergalkin/gophermart-diplom/internal/gophermart/service"
	"github.com/sergalkin/gophermart-diplom/internal/gophermart/utils"
)

type authHandler struct {
	service service.Auth
}

type authData struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func NewAuthHandler(s service.Auth) *authHandler {
	return &authHandler{service: s}
}

func (a *authHandler) Register() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		data := authData{}
		if errB := ctx.ShouldBindJSON(&data); errB != nil {
			casted, ok := errB.(validator.ValidationErrors)
			if !ok {
				ctx.JSON(http.StatusBadRequest, utils.NoPayloadProvided())
				return
			}

			ctx.JSON(http.StatusBadRequest, utils.FormResponseBindingError("error in provided data", casted))
			return
		}

		err := a.service.Register(data.Login, data.Password)
		if err != nil {
			if errors.Is(err, utils.ErrUserAlreadyExists) {
				ctx.JSON(http.StatusConflict, utils.JSONErrorResponse(err.Error()))
				return
			}

			ctx.JSON(http.StatusInternalServerError, utils.JSONErrorResponse(err.Error()))
		}

		utils.SetUserCookie(ctx, data.Login)

		ctx.JSON(http.StatusOK, gin.H{"message": "user successfully created and authenticated"})
	}
}

func (a *authHandler) Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		data := authData{}
		if errB := ctx.ShouldBindJSON(&data); errB != nil {
			casted, ok := errB.(validator.ValidationErrors)
			if !ok {
				ctx.JSON(http.StatusBadRequest, utils.NoPayloadProvided())
				return
			}

			ctx.JSON(http.StatusBadRequest, utils.FormResponseBindingError("error in provided data", casted))
			return
		}

		_, err := a.service.Login(data.Login, data.Password)
		if err != nil {
			if errors.Is(err, utils.ErrUserNotFound) || errors.Is(err, utils.ErrUserPasswordMissMatch) {
				ctx.JSON(http.StatusUnauthorized, utils.JSONErrorResponse("wrong username or password"))
				return
			}

			ctx.JSON(http.StatusInternalServerError, utils.JSONErrorResponse(err.Error()))
			return
		}

		utils.SetUserCookie(ctx, data.Login)
		ctx.JSON(http.StatusOK, gin.H{"message": "user successfully authenticated"})
	}
}
