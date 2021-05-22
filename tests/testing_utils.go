package tests

import (
	"fmt"
	"github.com/SherbazHashmi/goblog/api/controllers"
	"github.com/jinzhu/gorm"
	"log"
	"os"
)

func SetupDatabase(server *controllers.Server) {
	var err error
	testDbDriver := os.Getenv("TEST_DB_DRIVER")
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("TEST_DB_HOST"), os.Getenv("TEST_DB_USER"), os.Getenv("TEST_DB_PASSWORD"),
		os.Getenv("TEST_DB_NAME"), os.Getenv("TEST_DB_PORT"),
	)
	print(server)
	server.DB, err = gorm.Open(testDbDriver, dsn)

	if err != nil {
		fmt.Printf("Cannot connect to the database with the following properties\n%s\n", dsn)
		log.Fatal("[ERROR]", err)
		return
	}

	fmt.Printf("test database connection established")
}
