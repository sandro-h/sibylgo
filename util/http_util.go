package util

import (
	"encoding/json"
	"net/http"
)

// FetchJSONAsModel sends a GET request to the url and maps the JSON response to the
// target model.
func FetchJSONAsModel(client *http.Client, url string, target interface{}) error {
	r, err := client.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
