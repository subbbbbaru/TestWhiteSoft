package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	URL                 = "https://raw.githubusercontent.com/thewhitesoft/student-2023-assignment/main/data.json"
	replacementJSONFile = "replacement.json"
	resultFile          = "result.json"
)

type DataJson []string

type Replacements struct {
	Replacement string `json:"replacement"`
	Source      string `json:"source"`
}

var myClient = &http.Client{Timeout: 10 * time.Second}

func main() {
	// Read the replacement.json file
	replacementData, err := os.ReadFile(replacementJSONFile)
	if err != nil {
		log.Fatal(err)
	}

	var replacements []Replacements

	if err := json.Unmarshal(replacementData, &replacements); err != nil {
		log.Fatal(err)
	}

	var dataData DataJson

	// Get and Read the data.json file
	err = getJson(URL, &dataData)
	if err != nil {
		log.Fatal(err)
	}

	// Apply the replacement rules to the messages
	var modifiedData []string
	for _, message := range dataData {
		for i := len(replacements) - 1; i >= 0; i-- {
			replacement := replacements[i].Replacement
			source := replacements[i].Source

			if source == "" && strings.Contains(message, replacement) {
				message = ""
				break // Remove the message
			}
			if message != "" {
				message = strings.ReplaceAll(message, replacement, source)
			}

		}
		if message != "" {
			modifiedData = append(modifiedData, message)
		}
	}

	// Write the modified messages to result.json
	err = writeMessagesToFile(&modifiedData)
	if err != nil {
		log.Fatal(err)
	}

}

func writeMessagesToFile(modifiedData *[]string) error {
	resultData, err := json.Marshal(modifiedData)
	if err != nil {
		return err
	}
	if err := os.WriteFile(resultFile, resultData, 0644); err != nil {
		return err
	}
	return nil
}

func getJson(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
