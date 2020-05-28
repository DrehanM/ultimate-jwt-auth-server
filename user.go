package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

type User struct {
	ID 				string `json:"id"`
	HashedPassword 	string `json:"password"`
	Username 		string `json:"username"`
	Email 			string `json:"email"`
	RefreshToken 	string `json:"refresh_token"`
	RefreshCount	uint64 `json:"refresh_count"`
}

func (env *Env) VerifyAndRegisterNewUser(c *gin.Context) *User {
	var credentials Credentials
	var user User

	err := c.ShouldBindJSON(&credentials)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, "invalid json provided")
		return nil
	}

	usernameExists := !env.db.
		Where("username = ?", credentials.Username).
		First(&User{}).
		RecordNotFound()

	if usernameExists {
		c.JSON(http.StatusUnprocessableEntity, "username taken")
		return nil
	}

	emailExists := !env.db.
		Where("email = ?", credentials.Username).
		First(&User{}).
		RecordNotFound()

	if emailExists {
		c.JSON(http.StatusUnprocessableEntity, "email already used by another account")
		return nil
	}

	hashedPassword := hashAndSalt(credentials.Password)
	userID, err := generateBase64ID(USERID_LENGTH)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return nil
	}

	user = User {
		Email: credentials.Email,
		Username: credentials.Username,
		HashedPassword: hashedPassword,
		ID: userID,
	}

	err = env.db.Create(&user).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return nil
	}

	return &user
}

func (env *Env) VerifyUser(c *gin.Context) *User {
	var credentials Credentials
	var user User

	err := c.ShouldBindJSON(&credentials)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, "invalid json provided")
		return nil
	}

	err = env.db.
		Where("username = ? OR email = ?", credentials.Username, credentials.Email).
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
