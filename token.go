package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

const REFRESH_TOKEN_LENGTH = 64

// Generate JWT for the given user
// Return the generated signed token string with expiry (default = 15 min -- should be very short)
func (env *Env) CreateAccessToken(c *gin.Context, user *UserAuthInfo) (string, time.Time, error) {
	expiry :=  time.Now().Add(time.Minute * 15)

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["token_type"] = "access"
	claims["user_id"] = user.ID
	claims["exp"] = expiry.Unix() // 15 minutes

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return "", time.Time{}, err
	}
	return tokenString, expiry, nil
}

// Generate a new refresh_token as a 100 character Base64 ID for the given user
// Does not need JWX libraries since we store this in the auth database and verify via existence check
// Return the generated token string and token expiry (default = 28 days)
func (env *Env) CreateRefreshToken(c *gin.Context, user *UserAuthInfo) (string, time.Time, error) {
	expiry :=  time.Now().Add(time.Minute * 60 * 24 * 28)

	tokenString, err := generateBase64ID(REFRESH_TOKEN_LENGTH)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return "", time.Time{}, err
	}

	user.RefreshToken = tokenString
	user.RefreshCount += 1
	user.RefreshExpiry = expiry

	if err := env.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return "", time.Time{}, err
	}

	return tokenString, expiry, nil
}





