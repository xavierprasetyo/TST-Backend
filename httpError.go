package main

import (
	"encoding/json"
	"net/http"
)

type Message struct {
	Msg string `json:"msg"`
}

func sendError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Message{Msg: message})
}
