package service

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Book struct {
	ID     int
	Title  string
	Author string
	Genre  string
}

func (b Book) GetFullBook() string {
	return b.Title + " by " + b.Author
}

type BookService struct {
	db *sql.DB
}

func NewBookService(db *sql.DB) *BookService {
  return &BookService{db: db}
}

func (s *BookService) CreateBook(book *Book) error {
  query := "Insert into books (title, author, genre) values($1, $2, $3) RETURNING id, title, author, genre"

  var newBook Book
  err := s.db.QueryRow(query, book.Title, book.Author, book.Genre).Scan(&newBook.ID, &newBook.Title, &newBook.Author, &newBook.Genre)
  if err != nil {
    log.Fatal("Esse erro =>", err)
    return err
  }

  return nil
}

func (s *BookService) GetBooks() ([]Book, error) { 
  query := "Select id, title, author, genre from books"

  rows, err := s.db.Query(query)
  if err != nil {
    return nil, err
  }

  var books []Book
  for rows.Next() {
    var book Book
    err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Genre)
    if err != nil {
      return nil, err
    }
    books = append(books, book)
  }

  return books, nil
}

func (s *BookService) GetBookByID(id int) (*Book, error) {
  query := "select id, title, author, genre from books where id = $1"
  row := s.db.QueryRow(query, id)

  var book Book
  err := row.Scan(&book.ID, &book.Title, &book.Author, &book.Genre)
  if err != nil {
    return nil, err
  }

  return &book, nil
}

func (s *BookService) UpdateBook(book *Book) error {
  query := "update books set title=$1, author=$2, genre=$3 where id=$4"
  _, err := s.db.Exec(query, book.Title, book.Author, book.Genre, book.ID)
  return err
}

func (s *BookService) DeleteBook(id int) error {
  query := "delete from books where id=$1"
  _, err := s.db.Exec(query, id)
  return err
}

// SearchBooksByName busca livros pelo nome (título) no banco de dados.
func (s *BookService) SearchBooksByName(name string) ([]Book, error) {
	query := "SELECT id, title, author, genre FROM books WHERE title LIKE $1"
	rows, err := s.db.Query(query, "%"+name+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Genre); err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	return books, nil
}

// SimulateReading simula a leitura de um livro com base em um tempo de leitura.
func (s *BookService) SimulateReading(bookID int, duration time.Duration, results chan<- string) {
	book, err := s.GetBookByID(bookID)
	if err != nil || book == nil {
		results <- fmt.Sprintf("Livro com ID %d não encontrado.", bookID)
		return
	}

	time.Sleep(duration) // Simula o tempo de leitura.
	results <- fmt.Sprintf("Leitura do livro '%s' concluída!", book.Title)
}

func (s *BookService) SimulateMultipleReadings(bookIDs []int, duration time.Duration) []string {
	results := make(chan string, len(bookIDs)) // Canal com buffer para evitar bloqueio

	// Lança as goroutines para simular a leitura.
	for _, id := range bookIDs {
		go func(bookID int) {
			s.SimulateReading(bookID, duration, results)
		}(id)
	}

	var responses []string
	for range bookIDs {
		responses = append(responses, <-results)
	}
	close(results) // Fechamento do canal após coleta de todos os resultados

	return responses
}