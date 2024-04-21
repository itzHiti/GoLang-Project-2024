package main

import (
	"time"

	"OCM/pkg/OCM/model"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("Z2YgoOQifWNCwACjvkrlj34TdqjhH/redyS7i+d52p4=") // <- New secret key [Generated using test.go]

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateToken(user model.UserModel) (string, error) {
	expirationTime := jwt.NewNumericDate(time.Now().Add(24 * time.Hour))

	claims := &Claims{
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: expirationTime,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tknStr string) (*Claims, error) {
	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, err
		}
		return nil, err
	}

	if !tkn.Valid {
		return nil, err
	}

	return claims, nil
}
