package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
)

type Weather struct {
	CurrentCondition []struct {
		Humidity    string `json:"humidity"`
		TempC       string `json:"temp_C"`
		TempF       string `json:"temp_F"`
		WeatherDesc []struct {
			Desc string `json:"value"`
		} `json:"weatherDesc"`
	} `json:"current_condition"`
	HourlyWeather []struct {
		Hourly []struct {
			ChanceOfRain    string `json:"chanceofrain"`
			ChanceOfThunder string `json:"chanceofthunder"`
			TempC           string `json:"tempC"`
			TempF           string `json:"tempF"`
			Time            string `json:"time"`
			WeatherDesc     []struct {
				Desc string `json:"value"`
			}
		} `json:"hourly"`
	} `json:"weather"`
}

type WeatherClient struct {
	HTTPClient *http.Client
}

func (wc *WeatherClient) getWeather(location string) (*Weather, error) {
	weatherUrl := "https://wttr.in/" + location + "?format=j1"
	weatherResponse, err := wc.HTTPClient.Get(weatherUrl)
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

func PadToFour(s string) string {
	for len(s) < 4 {
		s = "0" + s
	}
	return s
}

func ConvertMilitaryToStandard(militaryTime string) (string, error) {
	if len(militaryTime) != 4 {
		militaryTime = PadToFour(militaryTime)
	}
	formattedTime := militaryTime[:2] + ":" + militaryTime[2:]

	t, err := time.Parse("15:04", formattedTime)
	if err != nil {
		return "", err
	}

	return t.Format("03:04 PM"), nil
}

func main() {

	wc := &WeatherClient{HTTPClient: http.DefaultClient}

	location := ""
	args := os.Args
	if len(args) > 1 {
		location = args[1]
	}

	JsonData, err := wc.getWeather(location)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if len(JsonData.CurrentCondition) > 0 {
		fmt.Println("--Current Conditions--")
		color.Red("Temp in C: %s | Temp in F: %s",
			JsonData.CurrentCondition[0].TempC,
			JsonData.CurrentCondition[0].TempF)
		color.Cyan("Humidity: %s", JsonData.CurrentCondition[0].Humidity)
		color.Green("Weather Description: %s\n\n", JsonData.CurrentCondition[0].WeatherDesc[0].Desc)
	} else {
		fmt.Println("No current weather data available")
	}

	if len(JsonData.HourlyWeather) > 0 {
		fmt.Println("--Timely Updates--")

		for _, value := range JsonData.HourlyWeather[0].Hourly {
			standardtime, err := ConvertMilitaryToStandard(value.Time)
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}

			color.Magenta(standardtime)
			color.Red("C: %v | F: %v", value.TempC, value.TempF)
			color.Blue("Chance of rain: %s | Chance of thunderstorm: %v",
				value.ChanceOfRain,
				value.ChanceOfThunder)
			color.HiCyan("Description : %s", value.WeatherDesc[0].Desc)
			fmt.Println("_________________________________________________")
			fmt.Println("")
		}
	} else {
		fmt.Println("Hourly update not available")
	}

}
