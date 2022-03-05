package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	sendinblue "github.com/sendinblue/APIv3-go-library/lib"
)

func init() {
	// Only load an .env file if we are in dev environment
	if os.Getenv("GO_ENV") == "development" {
		err := godotenv.Load(".env")

		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
}

func main() {
	// Create Gin router
	router := gin.Default()

	// Add CORS middleware
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	router.Use(cors.New(config))

	// Create status route
	router.GET("/api/status", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, gin.H{"message:": "Alive"})
	})

	router.GET("/api/getForecastFromLocation", getForecastFromLocation)
	router.GET("/api/getForecastFromLandmark", getForecastFromLandmark)
	router.POST("/api/sendMessage", sendMessage)

	// Start the server
	if os.Getenv("GO_ENV") == "development" {
		router.Run("0.0.0.0:8081")
	} else {
		prodDomain := os.Getenv("PROD_DOMAIN")
		port := os.Getenv("PORT")
		router.Run(prodDomain + ":" + port)
	}
}

func getForecastFromLocation(c *gin.Context) {
	lat := c.Query("lat")
	long := c.Query("long")

	// Using strings.Builder is the best way to concatenate strings
	var requestUrl strings.Builder
	requestUrl.WriteString("https://api.weather.gov/points/")
	requestUrl.WriteString(lat)
	requestUrl.WriteString(",")
	requestUrl.WriteString(long)

	// Poll NWS API to get gridpoints / forecast URL from lat/long
	body, err := makeGetRequest(requestUrl.String(), c)

	if err != nil {
		handleNwsError(err, c)
		return
	}

	var pointData PointResponse
	jsonErr := json.Unmarshal(body, &pointData)

	if jsonErr != nil {
		handleNwsError(jsonErr, c)
		return
	}

	// Get actual forecast with our returned gridpoints / forecast URL
	getForecastFromGridpoints(pointData.Properties.Forecast, c)

}

func getForecastFromGridpoints(forecastLink string, c *gin.Context) {
	// Poll NWS API for forecast for grid points
	body, err := makeGetRequest(forecastLink, c)

	// Return to the client
	handleForecastResponse(body, err, c)
}

func getForecastFromLandmark(c *gin.Context) {
	gridX := c.Query("gridX")
	gridY := c.Query("gridY")
	gridID := c.Query("gridID")

	// Using strings.Builder is the best way to concatenate strings
	var requestUrl strings.Builder
	requestUrl.WriteString("https://api.weather.gov/gridpoints/")
	requestUrl.WriteString(gridID)
	requestUrl.WriteString("/")
	requestUrl.WriteString(gridX)
	requestUrl.WriteString(",")
	requestUrl.WriteString(gridY)
	requestUrl.WriteString("/forecast")

	// Poll NWS API for forecast for gridpoints
	body, err := makeGetRequest(requestUrl.String(), c)

	// Return to the client
	handleForecastResponse(body, err, c)
}

func sendMessage(c *gin.Context) {
	// Initialize SendinBlue API
	var ctx context.Context
	cfg := sendinblue.NewConfiguration()
	//Configure API key authorization: api-key
	cfg.AddDefaultHeader("api-key", os.Getenv("SENDINBLUE_API_KEY"))

	sib := sendinblue.NewAPIClient(cfg)

	// Bind POST request body data to MessageData struct
	var postData MessageData
	err := c.BindJSON(&postData)

	if err != nil {
		handleMessageError(err, c)
	}

	// Create email message body
	htmlContent := fmt.Sprintf(`A new message has been sent from CDT Weather:<br />
              Sender Name: %s<br />
              Sender Email: %s<br /><br />
              Message: %s<br />`, postData.Name, postData.Email, postData.Message)

	// Build email metadata
	body := sendinblue.SendSmtpEmail{
		HtmlContent: htmlContent,
		Subject:     "New Message from Weather Backend",
		Sender: &sendinblue.SendSmtpEmailSender{
			Name:  "Weather Backend",
			Email: "zbtucker@gmail.com",
		},
		To: []sendinblue.SendSmtpEmailTo{{
			Name:  "Weather Backend",
			Email: "zbtucker@gmail.com",
		}},
	}
	// Send email
	_, _, sendErr := sib.TransactionalEmailsApi.SendTransacEmail(ctx, body)

	if sendErr != nil {
		handleMessageError(err, c)
	}

	// Return to the client
	c.JSON(http.StatusOK, gin.H{"message": "Message Sent!"})

}
