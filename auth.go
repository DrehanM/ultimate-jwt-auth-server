package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

var jwtKey = []byte(os.Getenv("API_SECRET"))

type Credentials struct {
	Password 	string `json:"password"`
	Email 		string `json:"email"`
	Username 	string `json:"username"`

}

func (env *Env) Login(c *gin.Context) {
	user := env.VerifyUser(c)
	if user == nil {
		return
	}

	env.processTokens(c, user)
}

func (env *Env) Register(c *gin.Context) {
	user := env.VerifyAndRegisterNewUser(c)
	if user == (nil) {
		return
	}

	env.processTokens(c, user)
}

func (env *Env) Refresh(c *gin.Context) {


}

func (env *Env) processTokens(c *gin.Context, user *User) {
	accessToken, accessExpiry, err := env.CreateAccessToken(c, user)
	if err != nil {
		return
	}
	refreshToken, refreshExpiry, err := env.CreateRefreshToken(c, user)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{"access_token": accessToken, "access_token_expiry": accessExpiry.Unix()})
	http.SetCookie(
		c.Writer,
		&http.Cookie{
			Name:     "refresh_token",
			Value:    refreshToken,
			Expires:  refreshExpiry,
			Path:     "/refresh",
			Domain:   "", //need to change this
			Secure:   true,
			HttpOnly: true,
		},
	)
}

