package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type LibraryBook struct {
	BookAuthor      string `json:"bookAuthor"`
	Heading         string `json:"heading"`
	Identifier      int    `json:"identifier"`
	ISBN            string `json:"isbn"`
	PublicationYear int    `json:"publicationYear"`
}

var libraryBooks []LibraryBook
var latestBookIdentifier int

// Handler to get all books
func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(libraryBooks)
}

// Handler to get a specific book by ID
func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}
	for _, book := range libraryBooks {
		if book.Identifier == id {
			json.NewEncoder(w).Encode(book)
			return
		}
	}
	http.Error(w, "Book not found", http.StatusNotFound)
}

// Handler to create a new book
func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newBook LibraryBook
	err := json.NewDecoder(r.Body).Decode(&newBook)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	latestBookIdentifier++
	newBook.Identifier = latestBookIdentifier
	libraryBooks = append(libraryBooks, newBook)
	json.NewEncoder(w).Encode(newBook)
}

// Handler to update an existing book by ID
func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	var updatedBook LibraryBook
	err = json.NewDecoder(r.Body).Decode(&updatedBook)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	for i, book := range libraryBooks {
		if book.Identifier == id {
			updatedBook.Identifier = id // Ensuring Identifier is unchanged
			libraryBooks[i] = updatedBook
			json.NewEncoder(w).Encode(updatedBook)
			return
		}
	}
	http.Error(w, "Book not found", http.StatusNotFound)
}

// Handler to delete a book by ID
func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	for i, book := range libraryBooks {
		if book.Identifier == id {
			libraryBooks = append(libraryBooks[:i], libraryBooks[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.Error(w, "Book not found", http.StatusNotFound)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/libraryBooks", getBooks).Methods(http.MethodGet)
	router.HandleFunc("/libraryBooks/{id}", getBook).Methods(http.MethodGet)
	router.HandleFunc("/libraryBooks", createBook).Methods(http.MethodPost)
	router.HandleFunc("/libraryBooks/{id}", updateBook).Methods(http.MethodPut)
	router.HandleFunc("/libraryBooks/{id}", deleteBook).Methods(http.MethodDelete)

	http.ListenAndServe(":8000", router)
}
