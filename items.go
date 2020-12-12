package main

import (
	"encoding/json"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type (
	Item struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Checked *bool  `json:"checked"`
		Crossed *bool  `json:"crossed"`
	}

	IDResponse struct {
		ID string `json:"id"`
	}

	CheckResponse struct {
		Success bool `json:"success"`
		Checked bool `json:"checked"`
	}

	CrossResponse struct {
		Success bool `json:"success"`
		Crossed bool `json:"checked"`
	}

	DeleteResponse struct {
		Success bool   `json:"success"`
		Message string `json:"msg"`
	}
)

func getItems(w http.ResponseWriter, r *http.Request) {
	client, ctx, err := createClient()
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Error Creating Firebase Client")
		return
	}
	var items []Item
	iter := client.Collection("items").Documents(ctx)
	for {
		var check, cross bool
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			sendError(w, http.StatusInternalServerError, "Error Iterating Data")
			return
		}
		check = doc.Data()["checked"].(bool)
		cross = doc.Data()["crossed"].(bool)
		items = append(items, Item{
			ID:      doc.Ref.ID,
			Name:    doc.Data()["name"].(string),
			Checked: &check,
			Crossed: &cross,
		})
	}
	defer client.Close()
	json.NewEncoder(w).Encode(items)
}

func addItem(w http.ResponseWriter, r *http.Request) {
	var item Item
	_ = json.NewDecoder(r.Body).Decode(&item)
	if item.Name == "" {
		sendError(w, http.StatusBadRequest, "Required Body Not Found")
		return
	}
	client, ctx, err := createClient()
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Error Creating Firebase Client")
		return
	}
	result, _, err := client.Collection("items").Add(ctx, map[string]interface{}{
		"name":    item.Name,
		"checked": false,
		"crossed": false,
	})
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Error Adding Data")
		return
	}
	defer client.Close()
	json.NewEncoder(w).Encode(IDResponse{ID: result.ID})
}

func checkItem(w http.ResponseWriter, r *http.Request) {
	var item Item
	_ = json.NewDecoder(r.Body).Decode(&item)
	if item.Checked == nil {
		sendError(w, http.StatusBadRequest, "Required Body Not Found")
		return
	}
	client, ctx, err := createClient()
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Error Creating Firebase Client")
		return
	}
	_, err = client.Collection("items").Doc(item.ID).Update(ctx, []firestore.Update{
		{
			Path:  "checked",
			Value: *item.Checked,
		},
	})
	if err != nil {
		if grpc.Code(err) == codes.NotFound {
			sendError(w, http.StatusBadRequest, "ID Not Found")
			return
		} else {
			sendError(w, http.StatusInternalServerError, "Error Updating Data")
			return
		}
	}
	defer client.Close()
	json.NewEncoder(w).Encode(CrossResponse{Success: true, Crossed: *item.Checked})
}

func crossItem(w http.ResponseWriter, r *http.Request) {
	var item Item
	_ = json.NewDecoder(r.Body).Decode(&item)
	if item.Crossed == nil {
		sendError(w, http.StatusBadRequest, "Required Body Not Found")
		return
	}
	client, ctx, err := createClient()
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Error Creating Firebase Client")
		return
	}
	_, err = client.Collection("items").Doc(item.ID).Update(ctx, []firestore.Update{
		{
			Path:  "crossed",
			Value: *item.Crossed,
		},
	})
	if err != nil {
		if grpc.Code(err) == codes.NotFound {
			sendError(w, http.StatusBadRequest, "ID Not Found")
			return
		} else {
			sendError(w, http.StatusInternalServerError, "Error Updating Data")
			return
		}
	}
	defer client.Close()
	json.NewEncoder(w).Encode(CheckResponse{Success: true, Checked: *item.Crossed})
}

func deleteItem(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Query()["id"]) == 0 {
		sendError(w, http.StatusBadRequest, "No ID Provided")
		return
	}
	deleteID := r.URL.Query()["id"][0]
	log.Println(deleteID)
	client, ctx, err := createClient()
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Error Creating Firebase Client")
		return
	}
	docRef := client.Collection("items").Doc(deleteID)
	_, err = docRef.Get(ctx)
	if err != nil {
		if grpc.Code(err) == codes.NotFound {
			sendError(w, http.StatusBadRequest, "ID Not Found")
			return
		}
		sendError(w, http.StatusInternalServerError, "Error Searching Data")
		return
	}
	_, err = docRef.Delete(ctx)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Error Deleting Data")
		return
	}
	defer client.Close()
	json.NewEncoder(w).Encode(DeleteResponse{Success: true, Message: "Item deleted!"})
}
