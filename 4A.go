package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"gitlab.com/slon/shad-go/gitfame/internal"
)

func main() {
	// Initialize command line arguments
	args := internal.NewCommandLineArgs()
	err1 := args.GetCommandLineArgs()
	if err1 != nil {
		log.Fatalf("Failed to get command line arguments:\n%v", err1)
	}

	// Load language mapping
	currentDir, err2 := os.Getwd()
	if err2 != nil {
		log.Fatalf("Failed to get current working directory:\n%v", err2)
	}

	// Find the root directory containing the go.mod file
	rootMarker := "go.mod"
	rootPath := currentDir
	for {
		if _, err3 := os.Stat(filepath.Join(rootPath, rootMarker)); err3 == nil {
			break
		}
		rootPath = filepath.Dir(rootPath)
	}

	mappingData, err4 := os.Open(filepath.Join(rootPath, "gitfame/configs/language_extensions.json"))
	if err4 != nil {
		log.Fatalf("Failed to open language extensions JSON file:\n%v", err4)
	}
	defer func(mappingData *os.File) {
		err5 := mappingData.Close()
		if err5 != nil {
			log.Fatalf("Failed to close language extensions JSON file:\n%v", err5)
		}
	}(mappingData)

	var languageMapping []internal.MappingEntity
	err6 := json.NewDecoder(mappingData).Decode(&languageMapping)
	if err6 != nil {
		log.Fatalf("Failed to decode language mapping JSON:\n%v", err6)
	}

	// Initialize file parameters
	filesParams := internal.NewFilesParams(languageMapping, args)
	filesParams.GetAllFiles(*filesParams.Cla.CommitPointer, *filesParams.Cla.RepositoryPath)

	// Count statistics
	stats := internal.CountStatistics(filesParams)

	// Sort and print results based on command line arguments
	stats.SortResults(args.SortOrderKey)
	stats.Print(args.OutputFormat)
}
