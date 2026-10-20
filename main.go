package main

import (
	"encoding/json"
	"net/http"
)

type Item struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

func main() {
	http.HandleFunc("GET /products", getProducts)
	http.ListenAndServe(":8080", nil)
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	products := []Item{
		{Id: 1, Name: "Milk", Price: 10000},
		{Id: 2, Name: "Bread", Price: 5000},
		{Id: 3, Name: "Coffee", Price: 50000},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}
