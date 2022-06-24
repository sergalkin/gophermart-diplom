package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func TextPlainConcern() gin.HandlerFunc {
	return func(context *gin.Context) {
		expected := "text/plain"

		if context.Request.Header.Get("Content-Type") != expected {
			context.Data(
				http.StatusBadRequest,
				expected,
				[]byte("Wrong Content-Type header. Expected to be "+expected),
			)
			return
		}

		context.Next()
	}
}
