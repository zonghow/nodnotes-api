package handlers

import (
	"errors"
	"fmt"
	"log"
	"nodnotes-api/app/db"
	"nodnotes-api/app/helpers"
	"nodnotes-api/app/models"

	"github.com/go-playground/validator"
	"github.com/kataras/iris/v12"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserReq struct {
	Username string `json:"username" validate:"required,min=3,max=40,username"`
	Password string `json:"password" validate:"required,min=6,max=20"`
}

type tokenResp struct {
	Token string `json:"token"`
}

// GetCurrentUser godoc
// @Accept  json
// @Produce  json
// @Success 200
// @Router /v1/user [get]
func GetCurrentUser(ctx iris.Context) {
	user := helpers.GetUser(ctx)
	if user == nil {
		ctx.StatusCode(iris.StatusForbidden)
		return
	}
}

// UserLogout godoc
// @Success 200
// @Router /v1/user/logout [post]
func UserLogout(ctx iris.Context) {
	ctx.RemoveCookie("AUTH_TOKEN")
	ctx.StatusCode(iris.StatusOK)
}

func UserLogin(ctx iris.Context) {
	var userReq UserReq
	err := ctx.ReadJSON(&userReq)
	if err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			validationErrors := helpers.WrapValidationErrors(errs)
			ctx.StopWithProblem(
				iris.StatusBadRequest,
				iris.NewProblem().
					Title("Validation error").
					Key("errors", validationErrors),
			)
		}
		ctx.StopWithStatus(iris.StatusInternalServerError)
		return
	}
	user := models.UserModel{
		Username: userReq.Username,
	}
	db.D.Client.Where("username = ?", user.Username).First(&user)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userReq.Password))
	if err != nil {
		ctx.StopWithProblem(
			iris.StatusBadRequest,
			iris.NewProblem().Title("username or password error"),
		)
		return
	}
	token, err := helpers.JWTSign(user.ID, user.Username)
	if err != nil {
		log.Fatalln(err)
		ctx.StopWithProblem(
			iris.StatusInternalServerError,
			iris.NewProblem().Title(err.Error()),
		)
		return
	}
	ctx.SetCookieKV("AUTH_TOKEN", token)
	ctx.JSON(
		tokenResp{
			Token: token,
		},
	)
	return
}

func UserSignin(ctx iris.Context) {
	var userReq UserReq
	err := ctx.ReadJSON(&userReq)
	if err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			validationErrors := helpers.WrapValidationErrors(errs)
			ctx.StopWithProblem(
				iris.StatusBadRequest,
				iris.NewProblem().
					Title("Validation error").
					Key("errors", validationErrors),
			)
		}
		ctx.StopWithStatus(iris.StatusInternalServerError)
		return
	}
	user := models.UserModel{
		Username: userReq.Username,
	}
	if err := db.D.Client.Where("username = ?", user.Username).First(&user).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.StopWithProblem(
				iris.StatusInternalServerError,
				iris.NewProblem().DetailErr(err),
			)
			return
		}
	} else {
		ctx.StopWithProblem(
			iris.StatusBadRequest,
			iris.NewProblem().Title(fmt.Sprintf("username: %s already exists", user.Username)),
		)
		return
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(userReq.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.StopWithProblem(
			iris.StatusInternalServerError,
			iris.NewProblem().Title(err.Error()),
		)
	}
	user.Password = string(hashPassword)
	if err := db.D.Client.Create(&user).Error; err != nil {
		ctx.StopWithProblem(
			iris.StatusInternalServerError,
			iris.NewProblem().Title(err.Error()),
		)
	}
	token, err := helpers.JWTSign(user.ID, user.Username)
	if err != nil {
		log.Fatalln(err)
		ctx.StopWithProblem(
			iris.StatusInternalServerError,
			iris.NewProblem().Title(err.Error()),
		)
		return
	}
	ctx.SetCookieKV("AUTH_TOKEN", token)
	ctx.JSON(
		tokenResp{
			Token: token,
		},
	)
	return
}
