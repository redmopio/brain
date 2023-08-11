package self

import (
	"bytes"
	"net/http"
)

// map json to a struct
type DataStruct struct {
	StationCode    string  `json:"station_code"`
	StreamName     string  `json:"stream_name"`
	RainIntensity  int     `json:"rain_intensity"`
	RainLevel      float64 `json:"rain_level"`
	CurrentWeather string  `json:"current_weather"`
	Datetime       string  `json:"date_time"`
}

func (b *BrainEngine) callHasuraEndpoint(bodyContent string) (bool, error) {
	// HTTP endpoint
	posturl := "https://redmop.practical-action.minsky.cc/api/rest/v1/instrument-records"

	// JSON body
	body := []byte(bodyContent)

	// Create a HTTP post request
	r, err := http.NewRequest("POST", posturl, bytes.NewBuffer(body))
	if err != nil {
		return false, err
	}

	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("x-hasura-admin-secret", b.HasuraToken)

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return false, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return false, nil
	}

	return true, nil
}
