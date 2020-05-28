package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"

)

var jwtKey = []byte(os.Getenv("API_SECRET")) // secret to sign JWT

func (env *Env) Login(c *gin.Context) {
	user := env.VerifyUser(c)
	if user == nil {
		return
	}

	env.processTokens(c, user)
}

func (env *Env) Register(c *gin.Context) {
	user := env.VerifyAndRegisterNewUser(c)
	if user == nil {
		return
	}

	env.processTokens(c, user)
}

func (env *Env) Refresh(c *gin.Context) {
	user := env.GetUserFromRefreshToken(c)
	if user == nil {
		return
	}

	env.processTokens(c, user)
}


// Return access token and access token expiry as payload
// Return refresh token as HttpOnly cookie
func (env *Env) processTokens(c *gin.Context, user *UserAuthInfo) {
	accessToken, accessExpiry, err := env.CreateAccessToken(c, user)
	if err != nil {
		return
	}
	refreshToken, refreshExpiry, err := env.CreateRefreshToken(c, user)
	if err != nil {
		return
	}

	http.SetCookie(
		c.Writer,
		&http.Cookie{
			Name:     "refresh_token",
			Value:    refreshToken,
			Expires:  refreshExpiry,
			Path:     "/refresh",
			Domain:   domain,
			Secure:   protocol == "https",
			HttpOnly: true,
		},
	)

	c.JSON(http.StatusOK, gin.H{"access_token": accessToken, "access_token_expiry": accessExpiry.Unix()})

}

