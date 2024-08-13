package framework

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func HTTPGet(url string, data interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read from response body: %s", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("url (%s) status code was %d, body: %s", url, resp.StatusCode, string(body))
	}

	return json.Unmarshal(body, data)
}
