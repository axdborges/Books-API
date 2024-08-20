package web

import (
	"encoding/json"
	"fmt"
	"gobooks/internal/service"
	"log"
	"net/http"
	"strconv"
)

type BookHandlers struct {
	service *service.BookService
}

func NewBookHandlers(service *service.BookService) *BookHandlers {
	return &BookHandlers{service: service}
}

func (h *BookHandlers) GetBooks(w http.ResponseWriter, r *http.Request) {
	books, err := h.service.GetBooks()
	fmt.Println(books)
	// Erro 500 caso tenha algum problema com o servidor
	if err != nil {
		http.Error(w, "Failed to get Books", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func (h *BookHandlers) CreateBook(w http.ResponseWriter, r *http.Request) {
	var book service.Book
	err := json.NewDecoder(r.Body).Decode(&book)
	// Erro 400 porque o json está errado com o modelo de book
	if err != nil {
		log.Fatal("Esse erro =>", err)
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	err = h.service.CreateBook(&book)
	if err != nil {
		log.Fatal("Esse erro =>", err)
		http.Error(w, "failed to create book", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

func (h *BookHandlers) GetBookByID(w http.ResponseWriter, r *http.Request) {
	// Pega o ID da rota GET /books/{id}
	idStr := r.PathValue("id")
	// Converte string para int
	id, err := strconv.Atoi(idStr)
	// Caso problema com ID errado, erro 400
	if err != nil {
		log.Fatal("Esse erro =>", err)
		http.Error(w, "invalid book ID", http.StatusBadRequest)
		return
	}

	book, err := h.service.GetBookByID(id)
	// Caso tenha problema com gerar o book, erro 500
	if err != nil {
		log.Fatal("Esse erro =>", err)
		http.Error(w, "failed to get book", http.StatusInternalServerError)
		return
	}
	// Caso não exista book, erro 404
	if book == nil {
		log.Fatal("Esse erro =>", err)
		http.Error(w, "book not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func (h *BookHandlers) UpdateBook(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Fatal("Esse erro =>", err)
		http.Error(w, "invalid book ID", http.StatusBadRequest)
		return
	}

	var book service.Book
	// Dessa forma a variável é declarada e já verificada, economizando linhas
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		log.Fatal("Esse erro =>", err)
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}
	book.ID = id

	if err := h.service.UpdateBook(&book); err != nil {
		log.Fatal("Esse erro =>", err)
		http.Error(w, "failed to update book", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(book)
}

func (h *BookHandlers) DeleteBook(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Fatal("Esse erro =>", err)
		http.Error(w, "invalid book ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteBook(id); err != nil {
		log.Fatal("Esse erro =>", err)
		http.Error(w, "failed to delete book", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
