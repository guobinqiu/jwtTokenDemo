package native

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"jwtTokenDemo/model"
	"net/http"
	"strings"
	"time"
)

func mainHandler() http.Handler {
	r := mux.NewRouter()
	r.Use(jwtAuth)
	r.HandleFunc("/v1/user/login", loginHandler).Methods(http.MethodGet)
	r.HandleFunc("/v1/user/info", infoHandler).Methods(http.MethodGet)
	return r
}

func jwtAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/user/login" {
			next.ServeHTTP(w, r)
			return
		}

		//获取token
		tokenString := ""
		bearToken := r.Header.Get("Authorization") //Authorization: Bearer <token>
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
			w.Write([]byte("Invalid auth token: " + err.Error()))
			return
		}
		if !token.Valid {
			w.Write([]byte("MyCustomClaims is not valid"))
			return
		}

		//验证通过
		ctx := context.WithValue(r.Context(), "user", claims.UserId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	//创建token
	claims := model.MyCustomClaims{
		"123456",
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(3 * time.Second).Unix(), //设置token过期时间
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(model.MySigningKey)

	w.Write([]byte(tokenString))
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(r.Context().Value("user").(string)))
}
