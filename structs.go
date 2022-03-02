package main

type PointResponse struct {
	Context  []interface{} `json:"@context"`
	ID       string        `json:"id"`
	Type     string        `json:"type"`
	Geometry struct {
		Type        string    `json:"type"`
		Coordinates []float64 `json:"coordinates"`
	} `json:"geometry"`
	Properties struct {
		ID                  string `json:"@id"`
		Type                string `json:"@type"`
		Cwa                 string `json:"cwa"`
		ForecastOffice      string `json:"forecastOffice"`
		GridID              string `json:"gridId"`
		GridX               int    `json:"gridX"`
		GridY               int    `json:"gridY"`
		Forecast            string `json:"forecast"`
		ForecastHourly      string `json:"forecastHourly"`
		ForecastGridData    string `json:"forecastGridData"`
		ObservationStations string `json:"observationStations"`
		RelativeLocation    struct {
			Type     string `json:"type"`
			Geometry struct {
				Type        string    `json:"type"`
				Coordinates []float64 `json:"coordinates"`
			} `json:"geometry"`
			Properties struct {
				City     string `json:"city"`
				State    string `json:"state"`
				Distance struct {
					UnitCode string  `json:"unitCode"`
					Value    float64 `json:"value"`
				} `json:"distance"`
				Bearing struct {
					UnitCode string `json:"unitCode"`
					Value    int    `json:"value"`
				} `json:"bearing"`
			} `json:"properties"`
		} `json:"relativeLocation"`
		ForecastZone    string `json:"forecastZone"`
		County          string `json:"county"`
		FireWeatherZone string `json:"fireWeatherZone"`
		TimeZone        string `json:"timeZone"`
		RadarStation    string `json:"radarStation"`
	} `json:"properties"`
}

type ForecastResponse struct {
	Context  []interface{} `json:"@context"`
	Type     string        `json:"type"`
	Geometry struct {
		Type        string        `json:"type"`
		Coordinates [][][]float64 `json:"coordinates"`
	} `json:"geometry"`
	Properties struct {
		// Updated           time.Time `json:"updated"`
		Units             string `json:"units"`
		ForecastGenerator string `json:"forecastGenerator"`
		// GeneratedAt       time.Time `json:"generatedAt"`
		// UpdateTime        time.Time `json:"updateTime"`
		// ValidTimes        time.Time `json:"validTimes"`
		Elevation struct {
			UnitCode string  `json:"unitCode"`
			Value    float64 `json:"value"`
		} `json:"elevation"`
		Periods []struct {
			Number           int         `json:"number"`
			Name             string      `json:"name"`
			StartTime        string      `json:"startTime"`
			EndTime          string      `json:"endTime"`
			IsDaytime        bool        `json:"isDaytime"`
			Temperature      int         `json:"temperature"`
			TemperatureUnit  string      `json:"temperatureUnit"`
			TemperatureTrend interface{} `json:"temperatureTrend"`
			WindSpeed        string      `json:"windSpeed"`
			WindDirection    string      `json:"windDirection"`
			Icon             string      `json:"icon"`
			ShortForecast    string      `json:"shortForecast"`
			DetailedForecast string      `json:"detailedForecast"`
		} `json:"periods"`
	} `json:"properties"`
}

type MessageData struct {
	Name    string
	Email   string
	Message string
}
