package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	createClient()
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	googleOauthId = os.Getenv("GOOGLE_OAUTH_ID")

	fmt.Println("Start")

	router := mux.NewRouter()
	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.Use(HandlerInit, AuthMiddleware, Logger)
	apiRouter.HandleFunc("/items", getItems).Methods(http.MethodGet)
	apiRouter.HandleFunc("/items/add", addItem).Methods(http.MethodPost)
	apiRouter.HandleFunc("/items/check", checkItem).Methods(http.MethodPut)
	apiRouter.HandleFunc("/items/cross", crossItem).Methods(http.MethodPut)
	apiRouter.HandleFunc("/items/delete", deleteItem).Methods(http.MethodDelete)

	log.Fatal(http.ListenAndServe(":8000", router))

}

func HandlerInit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, req)
	})
}
