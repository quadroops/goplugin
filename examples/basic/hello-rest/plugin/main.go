package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

var (
	commands = []string{"hello.world", "hello.msg"}
	port     = flag.String("port", "8080", "Set custom port")
)

// Response main response payload
type Response struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data"`
	Message *string     `json:"message,omitempty"`
}

// Output used as single response output
type Output struct {
	Response string `json:"response"`
}

// Request json payload
type Request struct {
	Command string `json:"command"`
	Payload string `json:"payload"`
}

func init() {
	flag.Parse()
}

func main() {
	app := fiber.New()
	app.Use(logger.New())

	app.Get("/ping", func(c *fiber.Ctx) error {
		out := Response{
			Status: "success",
			Data: Output{
				Response: "pong",
			},
		}

		return c.Status(http.StatusOK).JSON(out)
	})

	app.Post("/exec", func(c *fiber.Ctx) error {
		payload := new(Request)
		if err := c.BodyParser(payload); err != nil {
			return err
		}

		var command string
		for _, c := range commands {
			if payload.Command == c {
				command = c
			}
		}

		if command == "" {
			msg := "Unknown command"
			out := Response{
				Status:  "error",
				Message: &msg,
			}

			return c.Status(http.StatusBadRequest).JSON(out)
		}

		payloadHex, err := hex.DecodeString(payload.Payload)
		if err != nil {
			msg := err.Error()
			out := Response{
				Status:  "error",
				Message: &msg,
			}

			return c.Status(http.StatusBadRequest).JSON(out)
		}

		respHex := hex.EncodeToString(payloadHex)
		out := Response{
			Status: "success",
			Data: Output{
				Response: respHex,
			},
		}

		return c.Status(http.StatusAccepted).JSON(out)
	})

	app.Listen(fmt.Sprintf(":%s", *port))
}
