package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type Product struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Price int64  `json:"price"`
}

var Products = []Product{
	{ID: 1, Name: "Milk", Price: 10000},
	{ID: 2, Name: "Bread", Price: 5000},
	{ID: 3, Name: "Coffee", Price: 50000},
}

var mu sync.RWMutex

func main() {
	http.HandleFunc("GET /products", getProducts)
	http.HandleFunc("GET /products/{id}", getProduct)
	http.HandleFunc("POST /products", addProduct)
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Fatalf("Сервер не смог запуститься: %v", err)
	}
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	products := listProducts()
	writeJSON(w, http.StatusOK, products)
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id < 1 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Валидация идентификатора провалилась"})
		return
	}

	product, ok := findProduct(id)
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Продукт не найден"})
		return
	}
	writeJSON(w, http.StatusOK, product)
}

func findProduct(id int64) (Product, bool) {
	for _, val := range listProducts() {
		if val.ID == id {
			return val, true
		}
	}
	return Product{}, false
}

func listProducts() []Product {
	mu.RLock()
	defer mu.RUnlock()
	return Products
}

func addProduct(w http.ResponseWriter, r *http.Request) {
	p := Product{}
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "некорректный json"})
		return
	}

	if p.Name == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Имя обязательно"})
		return
	}

	if p.Price < 1 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Цена обязательна"})
		return
	}

	var lastID int64 = 0
	mu.Lock()
	defer mu.Unlock()
	for _, product := range Products {
		if product.Name == p.Name {
			writeJSON(w, http.StatusConflict, map[string]string{"error": "Такой товар уже существует"})
			return
		}
		if lastID < product.ID {
			lastID = product.ID
		}
	}

	p.ID = lastID + 1
	Products = append(Products, p)
	writeJSON(w, http.StatusCreated, p)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
