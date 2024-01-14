package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/robfig/cron/v3"
)

func myScheduledTask() {
	fmt.Println("Executing scheduled task: ", time.Now())
	makeHTTPRequest()
}

func makeHTTPRequest() {
	// Make an HTTP GET request to the local endpoint
	resp, err := http.Get("http://localhost:3000")
	if err != nil {
		log.Println("HTTP request error:", err)
		return
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		log.Println("Unexpected response status:", resp.Status)
		return
	}

	fmt.Println("response body: ", resp.Body)
}

func startWebServer() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, this is your web service!")
	})

	log.Fatal(app.Listen(":3000"))
}

func main() {
	// Create a new cron scheduler
	c := cron.New()

	// Schedule the task to run every minute
	_, err := c.AddFunc("* * * * *", myScheduledTask)
	if err != nil {
		log.Fatal(err)
	}

	// Start the cron scheduler in a goroutine
	go c.Start()

	// Run the web server
	startWebServer()

	// Keep the main goroutine running
	select {}
}
