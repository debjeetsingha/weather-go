package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/fatih/color"
)

type Weather struct {
	CurrentCondition []struct {
		Humidity    string `json:"humidity"`
		TempC       string `json:"temp_C"`
		TempF       string `json:"temp_F"`
		WeatherDesc []struct {
			Desc string `json:"value"`
		}
	} `json:"current_condition"`
}

func getWeather(location string) (*Weather, error) {
	weatherUrl := "https://wttr.in/" + location + "?format=j1"
	weatherResponse, err := http.Get(weatherUrl)
	if err != nil {
		fmt.Println("Error getting weather data: ", err)
		return nil, err
	}

	defer weatherResponse.Body.Close()

	if weatherResponse.StatusCode != http.StatusOK {
		fmt.Println("Error: ", weatherResponse.Status)
		return nil, fmt.Errorf("recived non-OK HTTP Status %s", weatherResponse.Status)
	}

	body, err := io.ReadAll(weatherResponse.Body)

	if err != nil {
		fmt.Println("Error reading response Body: ", err)
		return nil, err
	}

	var weatherData Weather
	err = json.Unmarshal(body, &weatherData)
	if err != nil {
		fmt.Println("Unable to parse json data: ", err)
	}
	return &weatherData, nil

}

func main() {

	fmt.Println("--Current Conditions--")

	location := ""
	args := os.Args
	if len(args) > 1 {
		location = args[1]
	}

	JsonData, err := getWeather(location)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if len(JsonData.CurrentCondition) > 0 {
		color.Red("Temp in C: %s", JsonData.CurrentCondition[0].TempC)
		color.Red("Temp in F: %s", JsonData.CurrentCondition[0].TempF)
		color.Cyan("Humidity: %s", JsonData.CurrentCondition[0].Humidity)
		color.Green("Weather Description: %s", JsonData.CurrentCondition[0].WeatherDesc[0].Desc)
	} else {
		fmt.Println("No current weather data available")
	}

}
