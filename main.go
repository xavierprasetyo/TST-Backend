package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	googleOauthId = os.Getenv("GOOGLE_OAUTH_ID")

	err = resetLog()
	if err != nil {
		log.Fatalf("Error resetting database")
	}

	fmt.Println("Start")

	router := mux.NewRouter()
	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.Use(HeaderInit, AuthMiddleware, Logger)
	apiRouter.HandleFunc("/items", getItems).Methods(http.MethodGet, http.MethodOptions)
	apiRouter.HandleFunc("/items/add", addItem).Methods(http.MethodPost)
	apiRouter.HandleFunc("/items/check", checkItem).Methods(http.MethodPut)
	apiRouter.HandleFunc("/items/cross", crossItem).Methods(http.MethodPut)
	apiRouter.HandleFunc("/items/delete", deleteItem).Methods(http.MethodDelete)
	apiRouter.HandleFunc("/items/delete", deleteItem).Methods(http.MethodDelete)

	apiRouter.HandleFunc("/log/reset", resetLogHandler).Methods(http.MethodPost)
	apiRouter.HandleFunc("/log", getLog).Methods(http.MethodGet)

	log.Fatal(http.ListenAndServe(":8000", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS", "DELETE"}), handlers.AllowedOrigins([]string{"*"}))(router)))
}
