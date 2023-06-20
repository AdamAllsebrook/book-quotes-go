package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"sort"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
)

// Boilerplate from https://echo.labstack.com/guide/templates/
type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

type Quote struct {
	Text          string
	Book          string
	DateSubmitted string
}

var app *pocketbase.PocketBase

func main() {
	app = pocketbase.New()

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
	records, err := app.Dao().FindRecordsByExpr("quotes")
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", err)
	}

	quotes := []Quote{}
	for _, record := range records {
		quotes = append(quotes, new_quote(record))
	}
	sort.Slice(quotes, func(i, j int) bool {
		return quotes[i].DateSubmitted > quotes[j].DateSubmitted
	})

	tmpl_args := map[string][]Quote{
		"Quotes": quotes,
	}
	return c.Render(http.StatusOK, "index.html", tmpl_args)
}

func add_quote(c echo.Context) error {
	text := c.FormValue("text")
	book := c.FormValue("book")

	collection, err := app.Dao().FindCollectionByNameOrId("quotes")
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", err)
	}

	record := models.NewRecord(collection)
	record.Set("text", text)
	record.Set("book", book)

	if err := app.Dao().SaveRecord(record); err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", err)
	}

	return c.Render(http.StatusOK, "quote", new_quote(record))
}

func new_quote(record *models.Record) Quote {
	return Quote{
		Text:          record.GetString("text"),
		Book:          record.GetString("book"),
		DateSubmitted: record.GetString("created"),
	}
}
