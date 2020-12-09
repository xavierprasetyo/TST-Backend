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

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/api", test).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", router))
}

func test(w http.ResponseWriter, r *http.Request) {
	msg := Message{Msg: "Halo Halo"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msg)
}
