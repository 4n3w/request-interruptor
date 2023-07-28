package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	ticker := time.NewTicker(1 * time.Minute)
	counter := 0
	url := os.Getenv("SERVICE_URL")

	if len(url) == 0 {
		log.Fatal("You need to specify SERVICE_URL in your environment variables")
	}

	for range ticker.C {
		counter++
		client := &http.Client{}
		req, err := http.NewRequest("GET", url, nil)

		if err != nil {
			fmt.Println("Error creating request:", err)
			continue
		}

		if counter%7 == 0 {
			// Create a context that will cancel the request after 1 nanosecond
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
			defer cancel()

			// Attach the context to our request
			req = req.WithContext(ctx)
		}

		resp, err := client.Do(req)

		if err != nil {

			fmt.Printf("[%s] Request error [count: %d] : %s\n", time.Now().Format(time.RFC3339), counter, err)
		} else {
			defer resp.Body.Close()
			fmt.Printf("[%s] Request successful. code: [%d] [count: %d]\n", time.Now().Format(time.RFC3339), resp.StatusCode, counter)

		}
	}
}
