package api

import (
	"fmt"
	"github.com/SherbazHashmi/goblog/api/controllers"
	"github.com/SherbazHashmi/goblog/api/seed"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var server = controllers.Server{}

func Run() {
	var err error
	err = godotenv.Load()

	if err != nil {
		log.Fatalf("Error getting envionrment, %v", err)
	} else {
		fmt.Println("environment variables loading")
	}

	server.Initialize(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PORT"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))
	seed.Load(server.DB)
	server.Run(":8080")
}
