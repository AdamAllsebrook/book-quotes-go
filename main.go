package main

import (
	"html/template"
	"log"
	"net/http"
	"sort"
	"time"
)

type Quote struct {
	Quote         string
	Book          string
	DateSubmitted string
}

var quotes = []Quote{
	{Quote: "Hello World", Book: "Some Book", DateSubmitted: time.Now().UTC().String()},
}

func main() {
	println("Staring server on port 8080")
	http.HandleFunc("/", index)
	http.HandleFunc("/add-quote/", add_quote)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func index(w http.ResponseWriter, r *http.Request) {
	sorted_quotes := quotes
	sort.Slice(sorted_quotes, func(i, j int) bool {
		return sorted_quotes[i].DateSubmitted > sorted_quotes[j].DateSubmitted
	})

	tmpl_args := map[string][]Quote{
		"Quotes": sorted_quotes,
	}
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, tmpl_args)
}

func add_quote(w http.ResponseWriter, r *http.Request) {
	quote := r.PostFormValue("quote")
	book := r.PostFormValue("book")
	date_submitted := time.Now().UTC().String()
	quote_record := Quote{Quote: quote, Book: book, DateSubmitted: date_submitted}
	quotes = append(quotes, quote_record)

	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.ExecuteTemplate(w, "quote", quote_record)
}
