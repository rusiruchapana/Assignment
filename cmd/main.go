package main

import (
	"book-api/internal/controller"
	"book-api/internal/repository"
	"book-api/internal/service"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	repo := repository.NewFileBookRepository("books.json")
	bookService := service.NewBookService(repo)
	bookController := controller.NewBookController(bookService)

	router := mux.NewRouter()

	router.HandleFunc("/books", bookController.CreateBook).Methods("POST")

	log.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
