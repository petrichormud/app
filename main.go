package main

import (
	"database/sql"
	"embed"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	fiber "github.com/gofiber/fiber/v2"
	html "github.com/gofiber/template/html/v2"
)

//go:embed web/views/*
var viewsfs embed.FS

func main() {
	db, err := sql.Open("mysql", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping: %v", err)
	}

	var hello string
	if err := db.QueryRow("SELECT 'hello from MySQL!';").Scan(&hello); err != nil {
		log.Fatalf("query: %v", err)
	}

	engine := html.NewFileSystem(http.FS(viewsfs), ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("web/views/index", fiber.Map{
			"Title": "Hello, World!",
		}, "web/views/layouts/main")
	})

	app.Static("/", "./web/static")

	app.Listen(":8008")
}
