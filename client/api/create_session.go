package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func CreateSession() (string, error) {
	request, err := http.NewRequest("POST", "http://127.0.0.1/v1/session/", nil)
	if err != nil {
		return "", err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var sessionResponse SessionResponse
	err = json.Unmarshal(body, &sessionResponse)
	if err != nil {
		return "", err
	}

	return sessionResponse.Data, nil
}
