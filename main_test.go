package main

import (
	"encoding/json"
	"os"
	"reflect"
	"sort"
	"testing"
)

const testResultFile = "expectedData.json"

func TestReplacementData(t *testing.T) {
	var bh Blackhole
	var expected []Replacements

	expectedData, err := os.ReadFile(replacementJSONFile)
	if err != nil {
		t.Errorf("read expected data from file %s error: %v", testResultFile, err)
	}
	if err := json.Unmarshal(expectedData, &expected); err != nil {
		t.Errorf("parse expected data from file %s error: %v", testResultFile, err)
	}

	if err := replacementData(&bh); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	sort.Slice(bh.Replacements, func(i, j int) bool { return i < j })
	sort.Slice(expected, func(i, j int) bool { return i < j })

	if !reflect.DeepEqual(bh.Replacements, expected) {
		t.Errorf("expected \n%v, but got \n%v", expected, bh.Replacements)
	}
}

func TestGetWrongDataFromURL(t *testing.T) {
	var target []string
	err := getWrongDataFromURL(URL, &target)
	if err != nil {
		t.Errorf("Error getting wrong data from URL: %v", err)
	}
	if len(target) == 0 {
		t.Errorf("Expected non-empty target slice, got empty slice")
	}
}

func TestWriteRepairDataToFile(t *testing.T) {
	data := []string{"hello", "world"}

	tmpfile, err := os.CreateTemp("", "result.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if err := writeRepairDataToFile(&data, tmpfile.Name()); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	contents, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("failed to read temp file: %v", err)
	}

	expected := `["hello","world"]`
	if string(contents) != expected {
		t.Errorf("expected %q, but got %q", expected, string(contents))
	}
}

func TestRepairWrongData(t *testing.T) {
	var expectedData []string

	testRawData, err := os.ReadFile(testResultFile)
	if err != nil {
		t.Errorf("read test data from file %s error: %v", testResultFile, err)
	}
	if err := json.Unmarshal(testRawData, &expectedData); err != nil {
		t.Errorf("get expected data from file %s error: %v", testResultFile, err)
	}

	var bh Blackhole
	err = replacementData(&bh)
	if err != nil {
		t.Errorf("read data from file %s error: %v", replacementJSONFile, err)
	}

	if err := getWrongDataFromURL(URL, &bh.Data.Messages); err != nil {
		t.Errorf("get wrong data from URL %s error: %v", URL, err)
	}

	actual, err := repairWrongData(&bh)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(actual, expectedData) {
		t.Errorf("expected %v, but got %v", expectedData, actual)
	}
}
