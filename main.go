package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/robfig/cron/v3"
)

var (
	once         sync.Once
	wg           sync.WaitGroup
	firstRequest bool
)

func myScheduledTask() {
	defer wg.Done() // Decrement the wait group counter when the task is done
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
		once.Do(func() {
			// Increment the wait group counter on the first request
			wg.Add(1)

			// Schedule the task to run every minute after the first request
			go func() {
				// Wait for the main function to finish before starting the cron job
				wg.Wait()

				// Create a new cron scheduler
				cronScheduler := cron.New()

				// Schedule the task to run every minute
				_, err := cronScheduler.AddFunc("* * * * *", myScheduledTask)
				if err != nil {
					log.Fatal(err)
				}

				// Start the cron scheduler
				cronScheduler.Start()
			}()
		})

		firstRequest = true

		return c.SendString("Hello, this is your web service!")
	})

	log.Fatal(app.Listen(":3000"))
}

func main() {
	// Run the web server
	startWebServer()
}
