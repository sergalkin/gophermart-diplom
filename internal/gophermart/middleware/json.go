package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func JSONConcern() gin.HandlerFunc {
	return func(context *gin.Context) {
		expected := "application/json"

		if context.Request.Header.Get("Content-Type") != expected {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Wrong Content-Type header. Expected to be " + expected,
				"errors":  []interface{}{},
			})
			return
		}

		context.Next()
	}
}
