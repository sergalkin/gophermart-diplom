package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/sergalkin/gophermart-diplom/internal/gophermart/storage"
	"github.com/sergalkin/gophermart-diplom/internal/gophermart/utils"
)

func AuthConcern() gin.HandlerFunc {
	return func(context *gin.Context) {
		login, err := utils.GetUserFromCookie(context)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusUnauthorized, utils.JSONErrorResponse("unauthorized"))
			return
		}

		if _, err = storage.NewDatabase().GetUserByName(login); err != nil {
			context.AbortWithStatusJSON(http.StatusUnauthorized, utils.JSONErrorResponse("unauthorized"))
			return
		}

		context.Next()
	}
}
