package main

import (
	"encoding/json"
	"net/http"
)

type Message struct {
	Msg string `json:"msg"`
}

func HeaderInit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, req)
	})
}

func sendError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Message{Msg: message})
}
