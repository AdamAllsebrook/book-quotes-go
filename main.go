package main

import (
	"html/template"
	"log"
	"net/http"
)

type Quote struct {
	Quote         string
	Book          string
	DateSubmitted string
}

var quotes = []Quote{
	{Quote: "Hello World", Book: "Some Book", DateSubmitted: "now"},
}

func main() {
	// Hello world, the web server

	println("Staring server on port 8080")
	http.HandleFunc("/", index)
	http.HandleFunc("/add-quote/", add_quote)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func index(w http.ResponseWriter, r *http.Request) {
	quotes := map[string][]Quote{
		"Quotes": quotes,
	}
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, quotes)
}

func add_quote(w http.ResponseWriter, r *http.Request) {
	quote := r.PostFormValue("quote")
	book := r.PostFormValue("book")
	date_submitted := "now"
	quote_record := Quote{Quote: quote, Book: book, DateSubmitted: date_submitted}
	quotes = append(quotes, quote_record)

	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.ExecuteTemplate(w, "quote", quote_record)
}
