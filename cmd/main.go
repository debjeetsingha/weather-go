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
		TempC       string `json:"temp_C"`
		TempF       string `json:"temp_F"`
		WeatherDesc []struct {
			Desc string `json:"value"`
		}
	} `json:"current_condition"`
}

func main() {

	color.Blue("Fetching Weather details...")
	client := &http.Client{}

	location := ""

	args := os.Args
	if len(args) > 1 {
		location = args[1]
	}

	weatherUrl := "https://wttr.in/" + location + "?format=j1"
	// fmt.Println(weatherUrl)

	weatherResponse, err := client.Get(weatherUrl)
	if err != nil {
		fmt.Println("Error getting weather data: ", err)
		return
	}

	defer weatherResponse.Body.Close()

	if weatherResponse.StatusCode != http.StatusOK {
		fmt.Println("Error: ", weatherResponse.Status)
		return
	}

	body, err := io.ReadAll(weatherResponse.Body)

	if err != nil {
		fmt.Println("Error reading response Body: ", err)
		return
	}

	var weatherData Weather
	err = json.Unmarshal(body, &weatherData)
	if err != nil {
		fmt.Println("Unable to parse json data: ", err)
		// panic(err)
	}

	color.Red("Temp in C: %s\n",
		weatherData.CurrentCondition[0].TempC)
	color.Red("Temp in F: %s\n",
		weatherData.CurrentCondition[0].TempF)
	color.Green("Weather Description: %s\n",
		weatherData.CurrentCondition[0].WeatherDesc[0].Desc)
}
