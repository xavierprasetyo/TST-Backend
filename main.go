package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Message struct {
	Msg string `json:"msg"`
}
type Auth struct {
	Token string `json:"token"`
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/api", test).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", router))

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
