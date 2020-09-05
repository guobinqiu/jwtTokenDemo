package gin

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"jwtTokenDemo/model"
	"strings"
	"time"
)

func mainHandler() *gin.Engine {
	r := gin.Default()
	r.GET("/v1/user/login", loginHandler)
	r.GET("/v1/user/info", jwtAuth, infoHandler)
	return r
}

func jwtAuth(c *gin.Context) {
	//获取token
	tokenString := ""
	bearToken := c.GetHeader("Authorization") //Authorization: Bearer <token>
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
		c.Writer.Write([]byte("Invalid auth token: " + err.Error()))
		return
	}
	if !token.Valid {
		c.Writer.Write([]byte("MyCustomClaims is not valid"))
		return
	}

	//验证通过
	c.Set("user", claims.UserId)
	c.Next()
}

func loginHandler(c *gin.Context) {
	//创建token
	claims := model.MyCustomClaims{
		"123456",
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(3 * time.Second).Unix(), //设置token过期时间
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(model.MySigningKey)

	c.Writer.Write([]byte(tokenString))
}

func infoHandler(c *gin.Context) {
	c.Writer.Write([]byte(c.Value("user").(string)))
}
