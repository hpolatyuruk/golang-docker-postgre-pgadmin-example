package data

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

var DB *gorm.DB
var err error

type Post struct {
	gorm.Model
	Author  string
	Message string
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file. Error: %v", err)
	}
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")

	connectionString := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, dbUser, dbName, dbPassword)

	for i := 0; i < 5; i++ {
		DB, err = gorm.Open("postgres", connectionString) // gorm checks Ping on Open
		if err == nil {
			break
		}
		fmt.Printf("An error occured while trying to connect to db. Error=%v", err)
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		log.Fatalf("Error opening the db connection. Error= %v", err)
	}

	// create table if it does not exist
	if !DB.HasTable(&Post{}) {
		DB.CreateTable(&Post{})
	}
}
