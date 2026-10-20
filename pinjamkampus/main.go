package main

import (
	"log"
	"net/http"
)

func main() {
	handler := NewHandler()
	log.Println("Server berjalan di port 8080...")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
