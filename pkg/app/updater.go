package app

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func CheckForNewUpdate(version string) bool {
	currentVersion := fmt.Sprintf("v%s", version)
	newVersion, err := getLatestVersionNumber()
	if err != nil {
		return false
	}

	if currentVersion == newVersion {
		return false
	}

	return true
}

func getLatestVersionNumber() (string, error) {
	req, err := http.NewRequest("GET", "https://github.com/metagunner/habheath/releases/latest", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	data := struct {
		TagName string `json:"tag_name"`
	}{}
	if err := dec.Decode(&data); err != nil {
		return "", err
	}

	return data.TagName, nil
}
