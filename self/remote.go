package self

// map json to a struct
type DataStruct struct {
	StationCode    string  `json:"station_code"`
	StreamName     string  `json:"stream_name"`
	RainIntensity  int     `json:"rain_intensity"`
	RainLevel      float64 `json:"rain_level"`
	CurrentWeather string  `json:"current_weather"`
	Timedate       string  `json:"timedate"`
}
