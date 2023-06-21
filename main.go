package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"log"
	"net/http"
	"time"
)

type Quote struct {
	Text      string
	Book      string
	CreatedAt string
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", "book-quotes.db")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	println("Staring server on port 8080")
	http.HandleFunc("/", index)
	http.HandleFunc("/add-quote/", addQuote)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func index(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM quotes ORDER BY created_at DESC")

	if err != nil {
		log.Fatal(err)
	}

	var quotes []Quote

	for rows.Next() {
		var id int
		var text string
		var book string
		var created_at string
		err = rows.Scan(&id, &text, &book, &created_at)

		if err != nil {
			log.Fatal(err)
		}

		quote := Quote{Text: text, Book: book, CreatedAt: created_at}
		quotes = append(quotes, quote)
	}

	defer rows.Close()

	tmpl_args := map[string][]Quote{
		"Quotes": quotes,
	}
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, tmpl_args)
}

func addQuote(w http.ResponseWriter, r *http.Request) {
	text := r.PostFormValue("text")
	book := r.PostFormValue("book")
	created_at := time.Now().UTC().String()
	quote := Quote{Text: text, Book: book, CreatedAt: created_at}

	stmt, err := db.Prepare("insert into quotes(text, book, created_at) values(?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(text, book, created_at)
	if err != nil {
		log.Fatal(err)
	}

	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.ExecuteTemplate(w, "quote", quote)
}
