package controllers

import (
	"fmt"
	"github.com/SherbazHashmi/goblog/api/models"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"net/http"
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

func (s *Server) Initialize(dbDriver, dbUser, dbPort, dbPassword, dbHost, dbName string) {
	var err error
	if dbDriver == "postgres" {
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			dbHost, dbUser, dbPassword, dbName, dbPort)

		s.DB, err = gorm.Open(dbDriver, dsn)

		if err != nil {
			fmt.Printf("Cannot connect to %s database\n", "")
			log.Fatal("[Error] ", err)
		} else {
			fmt.Printf("Database connection established.")
		}
	}

	s.DB.Debug().AutoMigrate(&models.User{}, &models.Ticket{}) // Database migration
	s.Router = mux.NewRouter()
	s.initializeRoutes()
}

func (s *Server) Run(addr string) {
	fmt.Println("Listening to port 8080")
	log.Fatal(http.ListenAndServe(addr, s.Router))
}
