package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

var (
	db     *sql.DB
	dbHost = os.Getenv("DB_HOST")
	dbPort = os.Getenv("DB_PORT")
	dbUser = os.Getenv("DB_USER")
	dbPass = os.Getenv("DB_PASS")
	dbName = os.Getenv("DB_NAME")
)

func main() {
	// Database
	var err error

	db, err = sql.Open("postgres", fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v?sslmode=disable",
		dbUser,
		dbPass,
		dbHost,
		dbPort,
		dbName,
	))
	if err != nil {
		fmt.Println(err)
	}

	defer db.Close()

	// Router
	r := gin.Default()

	r.POST("sync/product", ProductWebhook)
	r.POST("sync/customer", CustomerWebhook)

	r.Run()
}
