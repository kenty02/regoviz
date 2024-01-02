package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestListSamplesReturnsCorrectSamples(t *testing.T) {
	testDataDir := "testdata/tools_test/samples"
	samples, err := listSamples(testDataDir)
	if err != nil {
		t.Fatal(err)
	}

	for _, sample := range samples {
		if !strings.Contains(sample.InputExamples.Default, fmt.Sprintf("%s default input", sample.FileName)) {
			t.Errorf("expected %s default input, got %s", sample.FileName, sample.InputExamples.Default)
		}
		for inputName, input := range sample.InputExamples.AdditionalProps {
			if !strings.Contains(input, fmt.Sprintf("%s %s input", sample.FileName, inputName)) {
				t.Errorf("expected %s %s input, got %s", sample.FileName, inputName, input)
			}
		}
		if !strings.Contains(sample.DataExamples.Default, fmt.Sprintf("%s default data", sample.FileName)) {
			t.Errorf("expected %s default data, got %s", sample.FileName, sample.DataExamples.Default)
		}
		for dataName, data := range sample.DataExamples.AdditionalProps {
			if !strings.Contains(data, fmt.Sprintf("%s %s data", sample.FileName, dataName)) {
				t.Errorf("expected %s %s data, got %s", sample.FileName, dataName, data)
			}
		}
		if !strings.Contains(sample.QueryExamples.Default, fmt.Sprintf("%s default query", sample.FileName)) {
			t.Errorf("expected %s default query, got %s", sample.FileName, sample.QueryExamples.Default)
		}
		for queryName, query := range sample.QueryExamples.AdditionalProps {
			if !strings.Contains(query, fmt.Sprintf("%s %s query", sample.FileName, queryName)) {
				t.Errorf("expected %s %s query, got %s", sample.FileName, queryName, query)
			}
		}
	}
}

func TestReadSample(t *testing.T) {
	name := "../samples/rbac.rego"
	_, err := readSample(name, "./testdata")
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
