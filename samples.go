package main

import (
	"os"
	"path/filepath"
)

func listSamples() ([]string, error) {
	// list files in samples/ recursively
	var samples []string
	err := filepath.Walk("samples", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		samples = append(samples, path)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return samples, nil
}

func readSample(name string) (string, error) {
	// read from samples/NAME
	var sample string
	samplePath := filepath.Join("samples", name)
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
