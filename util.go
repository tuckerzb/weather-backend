package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

const nwsErrorMessage = "We are currently experiencing an issue communicating with the National Weather Service servers at this location. Please try again later."
const messageSendErrorMessage = "There was an issue sending your message."

func handleNwsError(err error, c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{"message": nwsErrorMessage})
	log.Println(err)
}

func handleMessageError(err error, c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{"message": messageSendErrorMessage})
	log.Println(err)
}

func makeGetRequest(url string, c *gin.Context) (body []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		handleNwsError(err, c)
		return nil, err
	}

	if resp.StatusCode != 200 {
		handleNwsError(err, c)
		return nil, err
	}

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		handleNwsError(err, c)
		return nil, err
	}
	return result, nil
}

func handleForecastResponse(body []byte, err error, c *gin.Context) {

	if err != nil {
		handleNwsError(err, c)
		return
	}

	var forecastData ForecastResponse
	jsonErr := json.Unmarshal(body, &forecastData)

	if jsonErr != nil {
		handleNwsError(jsonErr, c)
		return
	}

	c.IndentedJSON(http.StatusOK, forecastData.Properties.Periods)
}
