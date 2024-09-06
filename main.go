package main

import (
	"config"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func main() {

	name, err := config.GetConfig[int]("service")
	if err != nil {
		fmt.Print(err)
	}
	fmt.Print(name)
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Listen(":3000")

}
