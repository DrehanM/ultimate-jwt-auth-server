package main

import (
	"alo-auth/database"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

const defaultPort = "3010"

var (
	router = gin.Default()
)

type Env struct {
	db *database.Db
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	db, err := database.New(
		database.ConnString(
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	env := &Env{db: db}

	router.POST("/login", Env.Login)
	router.POST("/register", Env.Register)
	router.POST("/refresh", Env.Refresh)

	log.Fatal(router.Run(":8080"))

}
