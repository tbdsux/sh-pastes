package api

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/postgres/v3"
	gonanoid "github.com/matoous/go-nanoid/v2"

	_ "github.com/joho/godotenv/autoload"
)

// Run starts the API server
func Run() {
	app := fiber.New()
	store := postgres.New(postgres.Config{
		Host:     "postgres",
		Database: "shpastes",
		Table:    "pastes",
		Username: "postgres",
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Port:     5432,
	})

	// Get the root
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString(fmt.Sprintf(`SH Raw Pastes Hosting


cat file.txt | curl -X POST --data-binary @- %s
		`, os.Getenv("HOST")))

	})

	// Create a new paste
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

	// Get paste by id
	app.Get("/:id", func(c *fiber.Ctx) error {
		data, err := store.Get(c.Params("id"))
		if err != nil {
			return c.SendStatus(500)
		}

		if data == nil {
			return c.Status(404).SendString("Not Found")
		}

		return c.SendString(string(data))
	})

	app.Listen(":9999")
}
