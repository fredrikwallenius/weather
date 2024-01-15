package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type DataPoint struct {
	DateTime      string
	Temp          float64
	Wind          float64
	Precipitation float64
	Weather       string
}

func main() {

	city := getPromptedCity()
	fmt.Println("Fetching weather forecast for", city)

	lat, long := getCoordinates(city)
	fmt.Printf("Using these coordinates: %v, %v\n\n", lat, long)

	for true {
		fmt.Print("\033[H\033[2J") // https://stackoverflow.com/a/22892171
		fmt.Printf("Current forecast for %v (times are in GMT)\n", city)

		fetchAndPrintWeatherData(lat, long)

		fmt.Printf("(this forecast is updated every 5 minutes, last update was done %v)", time.Now().Format("15:04"))
		time.Sleep(5 * time.Minute)
	}
}

func getPromptedCity() string {
	args := os.Args

	if len(args) < 2 {
		fmt.Println("Enter a city you want the forecast for as a command line argument.")
		os.Exit(1)
	}

	arg := args[1]

	if len(arg) <= 3 {
		fmt.Println("A city name need to be 4 characters or longer")
		os.Exit(1)
	}

	return arg
}

func fetchAndPrintWeatherData(lat string, long string) {

	wData := getWeatherData(lat, long)

	// Initially thought I'd get these from the API too. But it made formatting the output
	// unessecary complex as I'd have to measure the length of them..
	temp_unit := "C"
	wind_unit := "m/s"
	prec_unit := "mm"

	dataPoints := []DataPoint{}

	// Look at the forecast for every 3 hours, one day forward, and put into my own serie
	for i := 0; i <= 24; i += 3 {
		dataPoint := wData.Properties.Timeseries[i]
		dataPoints = append(dataPoints, DataPoint{
			DateTime:      parseTime(dataPoint.Time),
			Temp:          dataPoint.Data.Instant.Details.AirTemperature,
			Wind:          dataPoint.Data.Instant.Details.WindSpeed,
			Precipitation: dataPoint.Data.Next1Hours.Details.PrecipitationAmount,
			Weather:       dataPoint.Data.Next1Hours.Summary.SymbolCode,
		})
	}

	var lineBreak strings.Builder
	fmt.Fprint(&lineBreak, "|-----------------------------------------------------------------------------------------------------------")

	var lineTimes strings.Builder
	var lineTemps strings.Builder
	var lineWinds strings.Builder
	var linePrecs strings.Builder
	var lineWeathers strings.Builder
	allLines := []*strings.Builder{
		&lineBreak,
		&lineTimes,
		&lineBreak,
		&lineTemps,
		&lineWinds,
		&linePrecs,
		&lineWeathers,
		&lineBreak,
	}

	for _, dp := range dataPoints {
		fmt.Fprintf(&lineTimes, "|%+10s ", dp.DateTime)
		fmt.Fprintf(&lineTemps, "|%+10s ", strconv.FormatFloat(dp.Temp, 'f', -1, 64)+" "+temp_unit)
		fmt.Fprintf(&lineWinds, "|%+10s ", strconv.FormatFloat(dp.Wind, 'f', -1, 64)+" "+wind_unit)
		fmt.Fprintf(&linePrecs, "|%+10s ", strconv.FormatFloat(dp.Precipitation, 'f', -1, 64)+" "+prec_unit)
		fmt.Fprint(&lineWeathers, "|      "+prettifyWeather(dp.Weather)+"   ")
	}
	for _, line := range allLines {
		fmt.Print(line.String() + "|\n")
	}

}
