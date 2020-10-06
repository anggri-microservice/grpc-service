package main

import "github.com/gofiber/fiber/v2"

func main() {
	app := fiber.New()

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendFile("./view/site/index.html")
	})

	app.Listen(":5000")
}
