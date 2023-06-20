package main

import (
	"html/template"
	"log"
	"net/http"
)

func main() {
	// Hello world, the web server

	println("Staring server on port 8080")
	http.HandleFunc("/hello", helloHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func helloHandler(w http.ResponseWriter, req *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}
