package main

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
)

func main() {
	mySigningKey := []byte("asfasfdafasdfdasfa.")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name": "huzb",
		"date": time.Now().Unix() + 5,
		"this": "tt",
	})
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		log.Println(err.Error())
		return
	}
	fmt.Println("加密后的token字符串", tokenString)

	tp, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("token:", tp)
	fmt.Println("token.Claims:", tp.Claims)
	fmt.Println(tp.Claims.(jwt.MapClaims)["name"])
}
