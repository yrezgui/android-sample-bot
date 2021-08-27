package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
)

type Sample struct {
	Id   string `json:"id"`
	Path string `json:"path"`
}

type SampleCollection struct {
	Items []Sample `json:"items"`
}

//go:embed "sample.json"
var rawSampleList []byte

//go:embed "excluded_folders"
var rawExcludedFolders string

func acrossAllSamples() {
	excluded_items := strings.Split(rawExcludedFolders, "\n")
	validSampleID := regexp.MustCompile(`^android\/`)
	unregisteredCount := 0

	// Get working directory path
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize sample collection and decode raw JSON file into it
	sampleCollection := SampleCollection{}
	if err := json.Unmarshal(rawSampleList, &sampleCollection); err != nil {
		panic(err)
	}

	samplesDB := make(map[string]Sample)
	for _, sample := range sampleCollection.Items {
		if !validSampleID.MatchString(sample.Id) {
			panic("Not supposed to happen: " + sample.Id)
		}

		simplifiedId := strings.Split(sample.Id, "android/")[1]
		samplesDB[simplifiedId] = sample
	}

	// fmt.Println(len(samplesDB))

	// List all repositories of working directory
	repositories, err := ioutil.ReadDir(workingDir)
	if err != nil {
		log.Fatal(err)
	}

	// Iterate through each file
	for _, repository := range repositories {
		if !repository.IsDir() || sliceSearch(repository.Name(), excluded_items) {
			continue
		}

		// List all files of repository directory
		files, err := ioutil.ReadDir(path.Join(workingDir, repository.Name()))
		if err != nil {
			log.Fatal(err)
		}

		for _, file := range files {
			// We ignore the item if it's not a directory or part of the excluded folders
			if !file.IsDir() || sliceSearch(file.Name(), excluded_items) {
				continue
			}

			sampleId := repository.Name() + "/" + file.Name()

			if _, ok := samplesDB[sampleId]; !ok {
				// fmt.Println("This sample (" + sampleId + ") isn't registered")
				fmt.Println(sampleId)
				unregisteredCount++
			}
		}
	}

	fmt.Println("There are " + strconv.Itoa(unregisteredCount) + " unregistered samples")
}

func sliceSearch(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
