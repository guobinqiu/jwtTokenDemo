package iris

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris"
	"jwtTokenDemo/model"
	"net/http"
	"strings"
	"time"
)

func mainHandler() http.Handler {
	app := iris.New()
	app.Get("/v1/user/login", loginHandler)
	app.Get("/v1/user/info", jwtAuth, infoHandler)
	app.Build()
	return app
}

func jwtAuth(ctx iris.Context) {
	//获取token
	tokenString := ""
	bearToken := ctx.Request().Header.Get("Authorization") //Authorization: Bearer <token>
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		tokenString = strArr[1]
	}

	//验证token
	claims := &model.MyCustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return model.MySigningKey, nil
	})
	if err != nil {
		ctx.Write([]byte("Invalid auth token: " + err.Error()))
		return
	}
	if !token.Valid {
		ctx.Write([]byte("MyCustomClaims is not valid"))
		return
	}

	//验证通过
	ctx.Values().Set("user", claims.UserId)
	ctx.Next()
}

func loginHandler(ctx iris.Context) {
	//创建token
	claims := model.MyCustomClaims{
		"123456",
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(3 * time.Second).Unix(), //设置token过期时间
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(model.MySigningKey)

	ctx.Write([]byte(tokenString))
}

func infoHandler(ctx iris.Context) {
	ctx.Write([]byte(ctx.Values().Get("user").(string)))
}
