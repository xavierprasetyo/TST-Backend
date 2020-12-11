package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
)

type (
	Item struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Checked bool   `json:"checked"`
		Crossed bool   `json:"crossed"`
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
)

var items []Item

func getItems(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(items)
}

func addItem(w http.ResponseWriter, r *http.Request) {
	var item Item
	_ = json.NewDecoder(r.Body).Decode(&item)
	item.ID = strconv.Itoa(rand.Intn(100000000))
	item.Checked = false
	item.Crossed = false
	items = append(items, item)
	json.NewEncoder(w).Encode(IDResponse{ID: item.ID})
}

func checkItem(w http.ResponseWriter, r *http.Request) {
	var item Item
	_ = json.NewDecoder(r.Body).Decode(&item)
	for i := range items {
		if items[i].ID == item.ID {
			items[i].Checked = item.Checked
		}
	}
	json.NewEncoder(w).Encode(CheckResponse{Success: true, Checked: item.Checked})
}

func crossItem(w http.ResponseWriter, r *http.Request) {
	var item Item
	_ = json.NewDecoder(r.Body).Decode(&item)
	for i := range items {
		if items[i].ID == item.ID {
			items[i].Crossed = item.Crossed
		}
	}
	json.NewEncoder(w).Encode(CrossResponse{Success: true, Crossed: item.Crossed})
}

func deleteItem(w http.ResponseWriter, r *http.Request) {
	var index int
	deleteID := r.URL.Query()["id"][0]
	for i := range items {
		if items[i].ID == deleteID {
			index = i
		}
	}
	items[index] = items[len(items)-1] // Copy last element to index i.
	items = items[:len(items)-1]
	json.NewEncoder(w).Encode(items)
}
