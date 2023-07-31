package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {

	counter := 0
	url := os.Getenv("SERVICE_URL")
	cancelIntervalMs, err := strconv.Atoi(os.Getenv("CANCEL_INTERVAL_MS"))

	if err != nil || cancelIntervalMs <= 0 {
		log.Fatal("You need to specify a valid CANCEL_INTERVAL_MS in your environment variables")
	}

	frequency, err := strconv.Atoi(os.Getenv("FREQUENCY_S"))

	if err != nil || frequency <= 0 {
		log.Fatal("You need to specify a valid FREQUENCY_S in your environment variables")
	}

	badRequestFrequency, err := strconv.Atoi(os.Getenv("BAD_REQUEST_FREQUENCY"))

	if err != nil || badRequestFrequency <= 0 {
		log.Fatal("You need to specify a valid BAD_REQUEST_FREQUENCY in your environment variables")
	}

	if len(url) == 0 {
		log.Fatal("You need to specify SERVICE_URL in your environment variables")
	}

	ticker := time.NewTicker(time.Duration(frequency) * time.Second)

	for range ticker.C {
		counter++
		client := &http.Client{}
		req, err := http.NewRequest("GET", url, nil)

		if err != nil {
			fmt.Println("Error creating request:", err)
			continue
		}

		if counter%badRequestFrequency == 0 {
			// Create a context that will cancel the request after 25 ms
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cancelIntervalMs)*time.Millisecond)
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
