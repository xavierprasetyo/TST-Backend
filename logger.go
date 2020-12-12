package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"cloud.google.com/go/firestore"
)

var (
	APIs = [...]string{"Get_Items", "Add_Item", "Check_Item", "Cross_Item", "Delete_Item"}
)

type Log struct {
	Get_Items   int64 `json:"get_items"`
	Add_Item    int64 `json:"add_items"`
	Check_Item  int64 `json:"check_items"`
	Cross_Item  int64 `json:"cross_items"`
	Delete_Item int64 `json:"delete_items"`
}

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
		default:
			next.ServeHTTP(w, req)
			return
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
		defer client.Close()
		next.ServeHTTP(w, req)
	})
}

func resetLog() error {
	client, ctx, err := createClient()
	if err != nil {
		return err
	}
	batch := client.Batch()
	logRef := client.Collection("log")
	for _, value := range APIs {
		docRef := logRef.Doc(value)
		batch.Set(docRef, map[string]interface{}{
			"count": 0,
		})
	}
	_, err = batch.Commit(ctx)
	if err != nil {
		return err
	}
	defer client.Close()
	return nil
}

func resetLogHandler(w http.ResponseWriter, r *http.Request) {
	err := resetLog()
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Error resetting log")
	}
	json.NewEncoder(w).Encode(DeleteResponse{Success: true, Message: "Log Resetted!"})
}

func getLog(w http.ResponseWriter, r *http.Request) {
	client, ctx, err := createClient()
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Error Creating Firebase Client")
		return
	}
	var countLog = make(map[string]int64)
	logRef := client.Collection("log")
	for _, value := range APIs {
		doc, err := logRef.Doc(value).Get(ctx)
		if err != nil {
			sendError(w, http.StatusInternalServerError, "Error retrieving log")
			return
		}
		countLog[value] = doc.Data()["count"].(int64)
	}
	log := Log{
		Get_Items:   countLog["Get_Items"],
		Add_Item:    countLog["Add_Item"],
		Check_Item:  countLog["Check_Item"],
		Cross_Item:  countLog["Cross_Item"],
		Delete_Item: countLog["Delete_Item"],
	}
	defer client.Close()
	json.NewEncoder(w).Encode(log)
}
