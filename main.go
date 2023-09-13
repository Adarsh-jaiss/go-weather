package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"log"	
)

const myAPI = "https://api.weatherapi.com/v1/forecast.json?key=6827c1cd51ee44e3902192605230604&q=%s&days=1&aqi=no&alerts=no"

var weather Weather

func main()  {

	// Treating this as a CLI Application --> Just uncomment the line below and run.
	fmt.Println("This is the CLI Output of your code")
	GetWeather()

	http.HandleFunc("/weather",Report)
	fmt.Println("--------------------------------------------")
	fmt.Println("This is the GUI output of the tool")
	fmt.Println("Starting a new server at port 3000")
	fmt.Println("Open --> http://localhost:3000/weather")
	log.Fatal(http.ListenAndServe(":3000",nil))
	
}

func Report(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
		// Indent the JSON response for better readability
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ") // Use two spaces for indentation
		err := encoder.Encode(weather)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
	}
}



func GetWeather() {

	var locationInput string
	fmt.Print("Enter the location : ")
	fmt.Scanln(&locationInput)

	apiURL := fmt.Sprintf(myAPI, locationInput)
	response, err := http.Get(apiURL)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		panic("Weather API is not working")
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	// Using the weather variable from global scope.
	err = json.Unmarshal(body, &weather)
	if err != nil {
		panic(err)
	}

	location, current, hours := weather.Location, weather.Current, weather.Forecast.Forecastday[0].Hour

	fmt.Printf("%s, %s : %.0f°C, %s\n", location.Name, location.Country, current.TempC, current.Condition.Text)

	for _, hour := range hours {
		date := time.Unix(hour.TimeEpoch, 0)
		myReport := fmt.Sprintf("%s - %.0f°C, %.0f, %s\n", date.Format("15:04"), hour.TempC, hour.Chanceofrain, hour.Condition.Text)

		if date.Before(time.Now()) {
			continue
		}
		fmt.Print(myReport)
	}

	// You can add conditional formatting based on the chance of rain here if needed
	// if hour.Chanceofrain < 40 {
	// 	fmt.Print(myReport)
	// }else{
	// 	color.Red(myReport)
	// }
}

type Weather struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`
	Current struct {
		TempC     float64 `json:"temp_c"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`
	Forecast struct {
		Forecastday []struct {
			Hour []struct {
				TimeEpoch int64   `json:"time_epoch"`
				TempC     float64 `json:"temp_c"`
				Condition struct {
					Text string `json:"text"`
				} `json:"condition"`
				Chanceofrain float64 `json:"chance_of_rain"`
			} `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
}
