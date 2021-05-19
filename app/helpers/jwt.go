package helpers

import (
	"log"
	"nodnotes-api/app/config"
	"strconv"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/jwt"
)

type UserClaims struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

func JWTSign(userID uint, username string) (string, error) {
	duration, _ := time.ParseDuration(strconv.Itoa(config.C.JWT.MaxAgeDay) + "d")
	token, err := jwt.Sign(jwt.HS256, []byte(config.C.JWT.Signature), UserClaims{userID, username}, jwt.MaxAge(duration))
	return string(token), err
}

func JWTVerify(token string) (*UserClaims, error) {
	verifiedToken, err := jwt.Verify(jwt.HS256, []byte(config.C.JWT.Signature), []byte(token))
	if err != nil {
		return nil, err
	}
	var userClaims UserClaims
	if err = verifiedToken.Claims(&userClaims); err != nil {
		return nil, err
	}
	return &userClaims, nil
}

func GetUser(ctx iris.Context) *UserClaims {
	cookie, err := ctx.Request().Cookie("AUTH_TOKEN")
	if err != nil {
		log.Println(err)
		return nil
	}
	userClaims, err := JWTVerify(cookie.Value)
	if err != nil {
		return nil
	}
	return userClaims
}
