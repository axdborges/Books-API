package main

import (
	"database/sql"
	"gobooks/internal/service"
	"gobooks/internal/web"
	"net/http"

	"log"
	// "os"

	_ "github.com/mattn/go-sqlite3"
	_ "github.com/lib/pq"
)

var connStr string = "user=postgres password=123456 dbname=gobooks host=localhost port=5432 sslmode=disable"

func main() {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		// panic(err)
		log.Fatalf("failed to connect to the database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Erro ao conectar ao banco de dados:", err)
	}
	defer db.Close()
	
	// Inicializando o serviço
	bookService := service.NewBookService(db)

	// Inicializando os handlers
	bookHandlers := web.NewBookHandlers(bookService)

		// Verifica se o CLI foi chamado
		// if len(os.Args) > 1 && os.Args[1] == "search" {
		// 	bookCLI := cli.NewBookCLI(bookService)
		// 	bookCLI.Run()
		// 	return
		// }

	// Inicia Servidor HTTP
	router := http.NewServeMux()

	// Configurando as rotas RESTful
	router.HandleFunc("GET /books", bookHandlers.GetBooks)
	router.HandleFunc("POST /books", bookHandlers.CreateBook)
	router.HandleFunc("GET /books/{id}", bookHandlers.GetBookByID)
	router.HandleFunc("PUT /books/{id}", bookHandlers.UpdateBook)
	router.HandleFunc("DELETE /books/{id}", bookHandlers.DeleteBook)

	// Iniciando o servidor
	log.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
	// http.ListenAndServe(":8080", router)
}