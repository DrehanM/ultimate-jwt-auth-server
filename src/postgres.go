package main
import (

	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

// postgres driver
_ "github.com/lib/pq"
)

type Db struct {
	*gorm.DB
}

// Db is our database struct used for interacting with the database

// New makes a new database using the connection string and
// returns it, otherwise returns the error
func InitDB(connString string) (*Db, error) {
	//db, err := sql.Open("postgres", connString)
	db, err := gorm.Open("postgres", connString)
	if err != nil {
		return nil, err
	}

	// Check that our connection is good
	/* err = db.Ping()
	if err != nil {
		return nil, err
	} */

	fmt.Println("Successfully connected to postgres DB")

	db = db.Debug() //remove when using outside of dev

	db.AutoMigrate(UserAuthInfo{}) //database migration


	return &Db{db}, nil
}

// ConnString returns a connection string based on the parameters it's given
// This would normally also contain the password, however we're not using one
func ConnString(host string, port string, user string, password string, dbName string) string {

	//if host is not given, init with default postgres settings
	if host == "" {
		host = "127.0.0.1"
		port = "5432"
		user = "postgres"
		dbName = "postgres"
		password = "password"

	}
	fmt.Printf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s\n",
		host,
		port,
		user,
		dbName,
		password,
	)
	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		host,
		port,
		user,
		dbName,
		password,
	)
}
