package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type customClaims struct {
	UserId     string   `json:"id"`
	ClaimRoles []string `json:"roles"`
	jwt.StandardClaims
}

func GetUserId(r *http.Request) int64 {
	authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
	jwtToken := authHeader[1]
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil {
		log.Printf("malformed jwt in GetUserId: %v", jwtToken)
		return 0
	}

	claims := token.Claims.(jwt.MapClaims)
	userId, ok := claims["id"].(string)
	if !ok {
		return 0
	}
	toReturn, err := strconv.ParseUint(userId, 10, 64)
	if err != nil {
		return 0
	}
	return int64(toReturn)
}

func GenerateToken(id string, roles []string, classCreator reflect.Type) (string, error) {
	te, _ := strconv.ParseInt(os.Getenv("TOKEN_EXP"), 10, 64)
	claims := &customClaims{
		UserId:     id,
		ClaimRoles: roles,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.NewTime(float64(time.Now().Add(time.Minute * time.Duration(te)).Unix())),
			Issuer:    classCreator.Name(),
		},
	}
	return GenerateTokenWithClaims(claims)
}

func GenerateTokenWithClaims(claims *customClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("creating a token %v", err))
	}
	return tokenString, err
}

func ValidateToken(jwtToken string) (jwt.Claims, error) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		for index, value := range claims {
			if index == "roles" {
				roles := value.([]interface{})
				for _, v := range roles {
					if v != "USER" && v != "SYSTEM" {
						return nil, errors.New("invalid role")
					}
				}
			}
		}
		return claims, nil
	}

	return nil, errors.New("not valid token")
}
