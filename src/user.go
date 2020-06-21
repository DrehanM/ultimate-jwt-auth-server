package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"time"
)

const USERID_LENGTH = 20

// object that user authorization info will be stored in
type UserAuthInfo struct {
	ID 				string `json:"id"`
	HashedPassword 	string `json:"password"`
	Username 		string `json:"username"`
	Email 			string `json:"email"`
	RefreshToken 	string `json:"refresh_token"`
	RefreshCount	uint64 `json:"refresh_count"`
	RefreshExpiry   time.Time `json:"last_refresh"`
}

// object used to marshal request body
type Credentials struct {
	Password 	string `json:"password"`
	Email 		string `json:"email"`
	Username 	string `json:"username"`
	Credentials string `json:"credentials"`
}


// Caller: Register
// Creates and returns new user pointer
// Validates that username and email do not exist in database already
// Hashes/salts and stores password
func (env *Env) VerifyAndRegisterNewUser(c *gin.Context) *UserAuthInfo {
	var credentials Credentials
	var user UserAuthInfo

	err := c.ShouldBindJSON(&credentials)
	if err != nil || credentials.Username == "" || credentials.Password == "" || credentials.Email  == "" {
		c.JSON(http.StatusUnprocessableEntity, "invalid json provided")
		return nil
	}

	usernameExists := !env.db.
		Where("username = ?", credentials.Username).
		First(&UserAuthInfo{}).
		RecordNotFound()

	if usernameExists {
		c.JSON(http.StatusUnprocessableEntity, "username taken")
		return nil
	}

	emailExists := !env.db.
		Where("email = ?", credentials.Email).
		First(&UserAuthInfo{}).
		RecordNotFound()

	if emailExists {
		c.JSON(http.StatusUnprocessableEntity, "email already used by another account")
		return nil
	}

	hashedPassword, err := hashAndSalt(credentials.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return nil
	}

	userID, err := generateBase64ID(USERID_LENGTH)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return nil
	}

	user = UserAuthInfo{
		Email: credentials.Email,
		Username: credentials.Username,
		HashedPassword: hashedPassword,
		ID: userID,
		RefreshToken: "", //to be updated in CreateRefreshToken
		RefreshCount: 0,
	}

	err = env.db.Create(&user).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return nil
	}

	return &user
}

// Caller: Login
// Returns user if provided username/email + password is valid
// Precondition: assumes caller only sends EITHER username OR email, not both
func (env *Env) VerifyUser(c *gin.Context) *UserAuthInfo {
	var credentials Credentials
	var user UserAuthInfo

	err := c.ShouldBindJSON(&credentials)
	if err != nil || credentials.Credentials == "" || credentials.Password == "" {
		c.JSON(http.StatusUnprocessableEntity, "invalid json provided")
		return nil
	}

	err = env.db.
		Where("username = ? OR email = ?", credentials.Credentials, credentials.Credentials).
		First(&user).
		Error

	if gorm.IsRecordNotFoundError(err) {
		c.JSON(http.StatusUnauthorized, "no such user")
		return nil
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, "unexpected error")
		return nil
	}

	if !verifyPassword(user.HashedPassword, credentials.Password) {
		c.JSON(http.StatusUnauthorized, "password is incorrect")
		return nil
	}

	return &user
}

// Caller: Refresh
// Returns UserAuthInfo if refresh token is valid and not expired
func (env *Env) GetUserFromRefreshToken(c *gin.Context) *UserAuthInfo {
	refreshTokenCookie, err := c.Request.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return nil
	}

	var user UserAuthInfo

	err = env.db.Where("refresh_token = ?", refreshTokenCookie.Value).First(&user).Error
	if gorm.IsRecordNotFoundError(err) {
		c.JSON(http.StatusUnauthorized, "invalid refresh token")
		return nil
	}
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return nil
	}

	if time.Since(user.RefreshExpiry) > 0{
		c.JSON(http.StatusUnauthorized, "refresh token expired. need to login")
		return nil
	}

	return &user
}
