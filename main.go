package main

import (
	"encoding/json"
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
	// items = append(items, Item{
	// 	ID:      "1",
	// 	Name:    "Telur",
	// 	Checked: false,
	// 	Crossed: true,
	// })
	// items = append(items, Item{
	// 	ID:      "2",
	// 	Name:    "Ayam",
	// 	Checked: true,
	// 	Crossed: true,
	// })

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

func Middleware(h http.Handler, middleware ...func(http.Handler) http.Handler) http.Handler {
	for _, mw := range middleware {
		h = mw(h)
	}
	return h
}

func HandlerInit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, req)
	})
}

func test(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var t Auth
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}
	isValid, error := verifyIdToken(t.Token)
	if error == nil {
		log.Println(isValid)
	} else {
		panic(error)
	}
}
