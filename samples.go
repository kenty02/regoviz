package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	api "vizrego-poc/vizrego"
)

func processFile(ext string, dataMap *map[string]map[string]string, path string) error {
	_, fileName := filepath.Split(path)
	policyName := strings.TrimSuffix(fileName, ext)
	dataName := "default"
	if len(policyName) > len(dataName)+len(ext) {
		dataName = policyName[:len(policyName)-len(dataName)-len(ext)-1]
	}
	dataContent, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if _, ok := (*dataMap)[policyName]; !ok {
		(*dataMap)[policyName] = make(map[string]string)
	}
	(*dataMap)[policyName][dataName] = string(dataContent)
	return nil
}

func listSamples(dir string) ([]api.Sample, error) {
	var samples []api.Sample
	inputs := map[string]map[string]string{}
	data := map[string]map[string]string{}
	queries := map[string]map[string]string{}
	regos := map[string]string{}

	// list files in samples/ recursively
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".data.json") {
			err := processFile(".data.json", &data, path)
			if err != nil {
				return err
			}
		} else if strings.HasSuffix(path, ".input.json") {
			err := processFile(".input.json", &inputs, path)
			if err != nil {
				return err
			}
		} else if strings.HasSuffix(path, ".query.txt") {
			err := processFile(".query.txt", &queries, path)
			if err != nil {
				return err
			}
		}
		if filepath.Ext(path) == ".rego" {
			sample, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			policyName := filepath.Base(strings.TrimSuffix(path, ".rego"))
			regos[policyName] = string(sample)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	for policyName, rego := range regos {
		defaultInput := ""
		if _, ok := inputs[policyName]; ok {
			defaultInput = inputs[policyName]["default"]
		} else {
			return nil, fmt.Errorf("no default input for %s", policyName)
		}
		// omit default input from inputs
		delete(inputs[policyName], "default")

		defaultData := ""
		if _, ok := data[policyName]; ok {
			defaultData = data[policyName]["default"]
		} else {
			return nil, fmt.Errorf("no default data for %s", policyName)
		}
		// omit default data from data
		delete(data[policyName], "default")

		defaultQueries := ""
		if _, ok := queries[policyName]; ok {
			defaultQueries = queries[policyName]["default"]
		} else {
			return nil, fmt.Errorf("no default queries for %s", policyName)
		}
		// omit default queries from queries
		delete(queries[policyName], "default")

		samples = append(samples, api.Sample{
			FileName: policyName,
			Content:  string(rego),
			DefaultInputs: api.SampleDefaultInputs{
				Default:         defaultInput,
				AdditionalProps: inputs[policyName],
			},
			DefaultData: api.SampleDefaultData{
				Default:         defaultData,
				AdditionalProps: data[policyName],
			},
			DefaultQueries: api.SampleDefaultQueries{
				Default:         defaultQueries,
				AdditionalProps: queries[policyName],
			},
		})
		// delete policyName from inputs, data, queries
		delete(inputs, policyName)
		delete(data, policyName)
		delete(queries, policyName)
	}

	return samples, nil
}

func readSample(name string, dir string) (string, error) {
	// read from samples/NAME
	var sample string
	samplePath := filepath.Join(dir, name)
	if _, err := os.Stat(samplePath); err != nil {
		return "", err
	}
	sampleBytes, err := os.ReadFile(samplePath)
	if err != nil {
		return "", err
	}
	sample = string(sampleBytes)
	return sample, nil
}
