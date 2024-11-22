package database

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func DBConn() (*sqlx.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	environment := os.Getenv("ENVIRONMENT")
	var driverName string
	var driverSource string

	if environment == "dev" {
		driverName = "postgres"
		driverSource = "user=mellob dbname=roc8 sslmode=disable password=mellob1989 host=192.168.29.73 port=5432"
	} else {
		driverName = "postgres"
		driverSource = "user=mellob dbname=roc8 sslmode=disable password=mellob1989 host=192.168.29.73"
	}

	db, err := sqlx.Connect(driverName, driverSource)
	// db, err := sqlx.Connect(
	// 	"postgres",
	// 	"user=admin dbname=ev sslmode=disable password=Mellob198978SadcDWFewd host=postgres.default.svc.cluster.local",
	// )
	// db, err := sqlx.Connect(
	// 	"postgres",
	// 	"user=admin dbname=ev sslmode=disable password=Mellob198978SadcDWFewd host=localhost port=5433",
	// )
	if err != nil {
		log.Fatalln(err)
		return nil, err // Return nil slice and error
	}

	// Test the connection to the database
	if err := db.Ping(); err != nil {
		log.Fatal(err)
		return nil, err // Return nil slice and error
	} else {
		log.Println("Successfully Connected")
		return db, nil
	}
}
