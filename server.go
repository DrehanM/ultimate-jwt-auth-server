package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

const defaultProtocol = "http"
const defaultPort = "3010"
const defaultHost = "127.0.0.1"

var (
	router = gin.Default()
)

// wrapper struct to pass db into resolvers
type Env struct {
	db *Db
}

var protocol string
var host string
var domain string

func main() {

	// fetch environment variables in .env
	//err := godotenv.Load("../.env")

	protocol = os.Getenv("PROTOCOL")
	if protocol == "" {
		protocol = defaultProtocol
	}

	host = os.Getenv("HOST")
	if host == "" {
		host = defaultHost
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	domain = host
	db, err := InitDB(
		ConnString(
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

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, //need to change
		AllowMethods:     []string{"POST"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
	}))


	env := &Env{db: db}

	router.POST("/login", env.Login)
	router.POST("/register", env.Register)
	router.POST("/refresh", env.Refresh)

	log.Fatal(router.Run(":"+port))

}
