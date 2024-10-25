package main

import (
	"database/sql"
	"fmt"
	"os"
	"purchasing_service/database"
	"purchasing_service/middleware"
	"purchasing_service/routers"
	"strconv"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB  *sql.DB
	err error
)

func main() {

	// Load environment variables from .env file
	err = godotenv.Load(".env")
	if err != nil {
		fmt.Println("ENV Failed to load!")
		panic(err)
	} else {
		fmt.Println("ENV Success !")
	}

	// Get database configurations from environment variables
	dbHost := os.Getenv("MYSQL_HOST")
	dbPort, _ := strconv.Atoi(os.Getenv("MYSQL_PORT"))
	dbUser := os.Getenv("MYSQL_USER")
	dbPassword := os.Getenv("MYSQL_PASSWORD")
	dbName := os.Getenv("MYSQL_DATABASE")

	// Create MySQL connection string (DSN)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbPort, dbName)

	// Open connection to MySQL database using GORM
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("DB Connection Failed!")
		panic(err)
	} else {
		fmt.Println("DB Connection Success!")
	}

	// Run database migrations and initialization
	// database.DbMigrate(db)
	database.Initialize(db)

	// Set the database for middleware
	middleware.SetDatabase(db)

	// Start the server
	PORT := ":3000"
	routers.StartServer().Run(PORT)
}
