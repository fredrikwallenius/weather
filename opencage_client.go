package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
)

func getCoordinates(city string) (string, string) {
	apiKey := os.Getenv("OPENCAGE_APIKEY")
	if apiKey == "" {
		fmt.Println("Please provide an API key for Opencage through the env variable OPENCAGE_APIKEY")
		os.Exit(1)
	}
	baseURL := "https://api.opencagedata.com/geocode/v1/json"

	requestURL := fmt.Sprintf("%s?q=%s&key=%s", baseURL, url.QueryEscape(city), apiKey)
	body := doGET(requestURL)

	var data map[string]interface{}
	err := json.Unmarshal(body, &data)
	if err != nil {
		os.Exit(1)
	}

	results := data["results"].([]interface{})
	if len(results) > 0 {
		geometry := results[0].(map[string]interface{})["geometry"].(map[string]interface{})
		lat := fmt.Sprintf("%.2f", geometry["lat"].(float64))
		lng := fmt.Sprintf("%.2f", geometry["lng"].(float64))
		return lat, lng
	}

	fmt.Println("I could not find the coordinates for city: ", city)
	os.Exit(1)
	return "0", "0" // I shouldn't need this line, somebody will have to explain to me why the compiler
	// don't understand it will never be reached
}
