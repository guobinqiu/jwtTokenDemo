package model

import "github.com/dgrijalva/jwt-go"

var MySigningKey = []byte("AllYourBase")

/*
JWT claims struct
*/
type MyCustomClaims struct {
	UserId string
	jwt.StandardClaims
}
