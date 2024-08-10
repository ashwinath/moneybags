package framework

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func HTTPGet(url string, data interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return json.Unmarshal(body, data)
}
