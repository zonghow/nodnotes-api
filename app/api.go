package app

import (
	"fmt"
	"nodnotes-api/app/config"
	"nodnotes-api/app/handlers"
	"nodnotes-api/docs"

	"github.com/iris-contrib/swagger"
	"github.com/iris-contrib/swagger/v12/swaggerFiles"
	"github.com/kataras/iris/v12"
)

func Start() error {
	app := iris.Default()

	docs.SwaggerInfo.Title = "nodnotes"
	docs.SwaggerInfo.Description = "nodnotes api docs"
	docs.SwaggerInfo.BasePath = "/v1"
	swaggerUI := swagger.Handler(swaggerFiles.Handler, swagger.Config{
		URL:          "/swagger/doc.json",
		DeepLinking:  true,
		DocExpansion: "list",
		DomID:        "#swagger-ui",
		Prefix:       "/swagger",
	})
	app.Get("/swagger", swaggerUI)
	app.Get("/swagger/{any:path}", swaggerUI)

	v1 := app.Party("/v1")
	usersAPI := v1.Party("/user")
	{
		usersAPI.Get("/", handlers.GetCurrentUser)
		usersAPI.Post("/signin", handlers.UserSignin)
		usersAPI.Post("/login", handlers.UserLogin)
		usersAPI.Post("/logout", handlers.UserLogout)
	}
	err := app.Listen(fmt.Sprintf(":%d", config.C.Http.Port))
	if err != nil {
		return err
	}
	return nil
}
