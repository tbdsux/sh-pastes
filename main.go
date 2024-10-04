package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/postgres/v3"
	gonanoid "github.com/matoous/go-nanoid/v2"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	app := fiber.New()
	store := postgres.New(postgres.Config{
		Database: "shpastes",
		Table:    "pastes",
		Username: "postgres",
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Port:     5432,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString(`SH Raw Pastes Hosting


cat file.txt | curl -X POST --data-binary @- {{host}}'
		`)

	})

	app.Post("/", func(c *fiber.Ctx) error {
		id, err := gonanoid.New()
		if err != nil {
			return c.SendStatus(500)
		}

		if err := store.Set(id, c.BodyRaw(), 0); err != nil {
			return c.SendStatus(500)
		}

		return c.SendString(id)
	})

	app.Get("/:id", func(c *fiber.Ctx) error {
		data, err := store.Get(c.Params("id"))
		if err != nil {
			return c.SendStatus(404)
		}

		return c.SendString(string(data))
	})

	app.Listen(":9999")
}
