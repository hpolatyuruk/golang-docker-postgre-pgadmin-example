package data

import (
	"fmt"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var DB *gorm.DB
var err error

type Post struct {
	gorm.Model
	Author  string
	Message string
}

func init() {
	fmt.Println("init.go init")
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_Password")
	dbName := os.Getenv("DB_NAME")
	connectionString := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, dbUser, dbName, dbPassword)
	fmt.Println(connectionString)
	if err != nil {
		fmt.Println(err)
		//panic(err)
	}
	for i := 0; i < 5; i++ {
		DB, err = gorm.Open("postgres", connectionString) // gorm checks Ping on Open
		if err == nil {
			break
		}
		fmt.Printf("An error occured while trying to connect to db. Error=%v", err)
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		fmt.Println(err)
		//panic(err)
	}

	DB.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))

	// create table if it does not exist
	if !DB.HasTable(&Post{}) {
		DB.CreateTable(&Post{})
	}

	testPost := Post{Author: "Dorper", Message: "GoDoRP is Dope"}
	DB.Create(&testPost)
}
