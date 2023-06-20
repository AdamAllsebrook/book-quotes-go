package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

// Boilerplate from https://echo.labstack.com/guide/templates/
type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

type Quote struct {
	Quote         string
	Book          string
	DateSubmitted string
}

var quotes = []Quote{
	{Quote: "Hello World", Book: "Some Book", DateSubmitted: time.Now().UTC().String()},
}

func main() {
	app := pocketbase.New()

	// Pre-compile templates
	t := &Template{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.Renderer = t
		e.Router.GET("/", index)
		e.Router.POST("/add-quote/", add_quote)

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

func index(c echo.Context) error {
	sorted_quotes := quotes
	sort.Slice(sorted_quotes, func(i, j int) bool {
		return sorted_quotes[i].DateSubmitted > sorted_quotes[j].DateSubmitted
	})

	tmpl_args := map[string][]Quote{
		"Quotes": sorted_quotes,
	}
	return c.Render(http.StatusOK, "index.html", tmpl_args)
}

func add_quote(c echo.Context) error {
	quote := c.FormValue("quote")
	book := c.FormValue("book")
	date_submitted := time.Now().UTC().String()
	quote_record := Quote{Quote: quote, Book: book, DateSubmitted: date_submitted}
	quotes = append(quotes, quote_record)

	return c.Render(http.StatusOK, "quote", quote_record)
}
