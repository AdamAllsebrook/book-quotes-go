package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

type Quote struct {
	Text          string
	Book          string
	DateSubmitted string
}

func main() {
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Get("/", index)
	app.Post("/add-quote/", addQuote)

	log.Fatal(app.Listen(":3000"))
}

func index(c *fiber.Ctx) error {
	// quotes := []Quote{}
	// for _, record := range records {
	// 	quotes = append(quotes, new_quote(record))
	// }
	// sort.Slice(quotes, func(i, j int) bool {
	// 	return quotes[i].DateSubmitted > quotes[j].DateSubmitted
	// })

	quotes := []Quote{
		{
			Text:          "This is a quote",
			Book:          "This is a book",
			DateSubmitted: "This is a date",
		},
	}
	return c.Render("index", fiber.Map{
		"Quotes": quotes,
	})
}

func addQuote(c *fiber.Ctx) error {
	text := c.FormValue("text")
	book := c.FormValue("book")
	quote := Quote{
		Text:          text,
		Book:          book,
		DateSubmitted: "Today",
	}
	log.Println("Added quote: ", quote)

	return c.Render("quote", fiber.Map{
		"Quote": quote,
	})
}
