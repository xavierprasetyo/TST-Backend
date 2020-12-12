package main

import (
	"net/http"
	"strings"

	"cloud.google.com/go/firestore"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var api string
		path := strings.ReplaceAll(req.URL.Path, "/api", "")
		switch path {
		case "/items":
			api = "Get_Items"
			break
		case "/items/add":
			api = "Add_Item"
			break
		case "/items/check":
			api = "Check_Item"
			break
		case "/items/cross":
			api = "Cross_Item"
			break
		case "/items/delete":
			api = "Delete_Item"
			break
		}
		client, ctx, err := createClient()
		if err != nil {
			sendError(w, http.StatusInternalServerError, "Error Creating Firebase Client")
			return
		}
		_, err = client.Collection("log").Doc(api).Update(ctx, []firestore.Update{
			{
				Path:  "count",
				Value: firestore.Increment(1),
			},
		})
		if err != nil {
			sendError(w, http.StatusInternalServerError, "Error in Logger")
			return
		}
		next.ServeHTTP(w, req)
	})
}
