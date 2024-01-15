package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// These structs I let ChatGPT create for me

// Note for self:
// The timeseries is hourly the coming 48 hours, after that the forecast is per 6 hours

type Details struct {
	AirPressureAtSeaLevel float64 `json:"air_pressure_at_sea_level"`
	AirTemperature        float64 `json:"air_temperature"`
	CloudAreaFraction     float64 `json:"cloud_area_fraction"`
	RelativeHumidity      float64 `json:"relative_humidity"`
	WindFromDirection     float64 `json:"wind_from_direction"`
	WindSpeed             float64 `json:"wind_speed"`
}

type Instant struct {
	Details Details `json:"details"`
}

type Next1Hours struct {
	Summary struct {
		SymbolCode string `json:"symbol_code"`
	} `json:"summary"`
	Details struct {
		PrecipitationAmount float64 `json:"precipitation_amount"`
	} `json:"details"`
}

type Timeseries struct {
	Time string `json:"time"`
	Data struct {
		Instant     Instant    `json:"instant"`
		Next12Hours struct{}   `json:"next_12_hours"`
		Next1Hours  Next1Hours `json:"next_1_hours"`
		Next6Hours  Next1Hours `json:"next_6_hours"`
	} `json:"data"`
}

type Properties struct {
	Meta struct {
		UpdatedAt string `json:"updated_at"`
		Units     struct {
			AirPressureAtSeaLevel string `json:"air_pressure_at_sea_level"`
			AirTemperature        string `json:"air_temperature"`
			CloudAreaFraction     string `json:"cloud_area_fraction"`
			PrecipitationAmount   string `json:"precipitation_amount"`
			RelativeHumidity      string `json:"relative_humidity"`
			WindFromDirection     string `json:"wind_from_direction"`
			WindSpeed             string `json:"wind_speed"`
		} `json:"units"`
	} `json:"meta"`
	Timeseries []Timeseries `json:"timeseries"`
}

type Coordinates []float64

type Geometry struct {
	Type        string      `json:"type"`
	Coordinates Coordinates `json:"coordinates"`
}

type Feature struct {
	Type       string     `json:"type"`
	Geometry   Geometry   `json:"geometry"`
	Properties Properties `json:"properties"`
}

func getWeatherData(lat string, long string) Feature {

	var fileContent []byte
	var err error

	var url strings.Builder
	fmt.Fprintf(&url, "https://api.met.no/weatherapi/locationforecast/2.0/compact?lat=%v&lon=%v", lat, long)

	fileContent = doGET(url.String())

	var featureData Feature

	err = json.Unmarshal(fileContent, &featureData)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		os.Exit(1)
	}

	return featureData
}

func prettifyWeather(weather string) string {
	switch {
	case strings.Contains(weather, "cloud"):
		return "☁️ "
	case strings.Contains(weather, "snow") || strings.Contains(weather, "sleet"):
		return "❄️ "
	case strings.Contains(weather, "clearsky_night") || strings.Contains(weather, "fair_night"):
		return "🌙"
	case strings.Contains(weather, "clearsky_day") || strings.Contains(weather, "fair_day"):
		return "☀️ "
	case strings.Contains(weather, "thunder"):
		return "🌩️ "
	case strings.Contains(weather, "rain"):
		return "🌧️ "
	case strings.Contains(weather, "fog"):
		return "🌫️ "
	default:
		{
			fmt.Println("Unknown code:", weather)
			return "❓"
		}
	}

}
