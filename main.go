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

type DataJson struct {
	Messages []string
}

type Replacements struct {
	Replacement string `json:"replacement"`
	Source      string `json:"source"`
}

type Blackhole struct {
	Replacements []Replacements
	Data         DataJson
}

func main() {
	var bh Blackhole

	// Read the replacement.json file
	err := replacementData(&bh)
	if err != nil {
		log.Fatal(err)
	}

	// Get and Read the data.json file from URL
	if err := getWrongDataFromURL(URL, &bh.Data.Messages); err != nil {
		log.Fatal(err)
	}

	repairData, err := repairWrongData(&bh)
	if err != nil {
		log.Fatal(err)
	}

	// Write the modified messages to result.json
	err = writeRepairDataToFile(&repairData, resultFile)
	if err != nil {
		log.Fatal(err)
	}

}

func repairWrongData(bh *Blackhole) ([]string, error) {
	var modifiedData []string

	// Apply the replacement rules to the messages
	for _, message := range bh.Data.Messages {
		for i := len(bh.Replacements) - 1; i >= 0; i-- {
			replacement := bh.Replacements[i].Replacement
			source := bh.Replacements[i].Source

			if source == "" && strings.Contains(message, replacement) {
				message = ""
				break // Remove the message
			}
			message = strings.ReplaceAll(message, replacement, source)

		}
		if message != "" {
			modifiedData = append(modifiedData, message)
		}
	}
	return modifiedData, nil

}

func writeRepairDataToFile(modifiedData *[]string, resultFile string) error {
	resultData, err := json.Marshal(modifiedData)
	if err != nil {
		return err
	}
	if err := os.WriteFile(resultFile, resultData, 0644); err != nil {
		return err
	}
	return nil
}

func replacementData(bh *Blackhole) error {
	// Read the replacement.json file
	replacementData, err := os.ReadFile(replacementJSONFile)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(replacementData, &bh.Replacements); err != nil {
		return err
	}
	return nil
}

func getWrongDataFromURL(url string, target interface{}) error {
	var myClient = &http.Client{Timeout: 10 * time.Second}
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
