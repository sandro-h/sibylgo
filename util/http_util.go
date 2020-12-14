package util

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// FetchJSONAsModel sends a GET request to the url and maps the JSON response to the
// target model.
func FetchJSONAsModel(client *http.Client, url string, user string, password string, target interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	if user != "" {
		req.SetBasicAuth(user, password)
	}

	r, err := client.Do(req)
	if err != nil {
		return err
	}
	if r.StatusCode < 200 || r.StatusCode >= 300 {
		return fmt.Errorf("request returned HTTP %d", r.StatusCode)
	}

	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
