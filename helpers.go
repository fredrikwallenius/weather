package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Make a GET request, return the body as a byte slice
func doGET(apiURL string) []byte {
	client := &http.Client{}

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		os.Exit(1)
	}

	// The met.no API recommended using a unique user-agent
	req.Header.Set("User-Agent", "FW-App/0.1")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("API request failed. Status:", resp.Status)
		os.Exit(1)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		os.Exit(1)
	}

	return responseBody
}

// Parses a datetime string in the format "2024-01-15T08:00:00Z"
// and return just the time in the format "15:04".
// If it fails it will exit the application.
func parseTime(ts string) string {
	parsedTime, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		fmt.Println("Error parsing time:", err)
		os.Exit(1)
	}

	return parsedTime.Format("15:04")
}
