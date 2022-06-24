package utils

import (
	"github.com/gin-gonic/gin"
)

const UserCookie = "user"

func SetUserCookie(ctx *gin.Context, login string) {
	sha, _ := Encode(login)

	ctx.SetCookie(UserCookie, sha, 36000, "/", "", false, false)
}

func GetUserFromCookie(ctx *gin.Context) (string, error) {
	var login string

	v, err := ctx.Cookie(UserCookie)
	if err != nil {
		return login, err
	}

	err = Decode(v, &login)
	if err != nil {
		return login, err
	}

	return login, err
}
