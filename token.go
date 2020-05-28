package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"time"
)

func (env *Env) CreateAccessToken(c *gin.Context, user *User) (string, time.Time, error) {
	expiry :=  time.Now().Add(time.Minute * 15)

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["token_type"] = "access"
	claims["user_id"] = user.ID
	claims["exp"] = expiry.Unix() // 15 minutes

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return "", time.Time{}, err
	}
	return tokenString, expiry, nil
}

func (env *Env) CreateRefreshToken(c *gin.Context, user *User) (string, time.Time, error) {
	expiry :=  time.Now().Add(time.Minute * 60 * 24 * 28)

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["token_type"] = "refresh"
	claims["user_id"] = user.ID
	claims["exp"] = expiry.Unix()  // 28 days

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return "", time.Time{}, err
	}

	user.RefreshToken = tokenString
	user.RefreshCount += 1

	if err := env.db.Model(user).Update("value", "refresh_count").Error; err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return "", time.Time{}, err
	}

	return tokenString, expiry, nil
}




