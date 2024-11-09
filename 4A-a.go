package internal

import (
	"fmt"
	"os"
	"os/exec"

	flag "github.com/spf13/pflag"
)

type CommandLineArgs struct {
	RepositoryPath       *string
	CommitPointer        *string
	SortOrderKey         string
	UseCommitter         *bool
	OutputFormat         string
	FileExtensions       *[]string
	ProgrammingLanguages *[]string
	ExcludePatterns      *[]string
	IncludePatterns      *[]string
}

func NewCommandLineArgs() *CommandLineArgs {
	return &CommandLineArgs{}
}

func (cla *CommandLineArgs) GetCommandLineArgs() error {
	cla.RepositoryPath = flag.String("repository", "./", "Path to the repository.")
	cla.CommitPointer = flag.String("revision", "HEAD", "Reference to a specific commit.")
	sortOrder := flag.String("order-by", "lines", "Key for sorting.")
	cla.UseCommitter = flag.Bool("use-committer", false, "Enable use of committer information.")
	outputFormat := flag.String("format", "tabular", "Format for output.")
	cla.FileExtensions = flag.StringSlice("extensions", []string{}, "File extensions to search for.")
	cla.ProgrammingLanguages = flag.StringSlice("languages", []string{}, "Programming languages to search for.")
	cla.ExcludePatterns = flag.StringSlice("exclude", []string{}, "Glob patterns to exclude.")
	cla.IncludePatterns = flag.StringSlice("restrict-to", []string{}, "Patterns to include in the search.")

	flag.Parse()

	// Validate repository path
	if _, err := os.Stat(*cla.RepositoryPath); os.IsNotExist(err) {
		return fmt.Errorf("repository path missing: %s", *cla.RepositoryPath)
	}

	// Validate commit pointer by executing git show command
	gitShowCmd := exec.Command("git", "show", *cla.CommitPointer)
	gitShowCmd.Dir = *cla.RepositoryPath

	err := gitShowCmd.Run()
	if err != nil {
		return fmt.Errorf("commit missing: %s", *cla.CommitPointer)
	}

	// Validate sort order key
	switch *sortOrder {
	case "lines", "commits", "files":
		cla.SortOrderKey = *sortOrder
	default:
		return fmt.Errorf("key error: %s. Permitted keys: 'lines', 'commits', 'files'", *sortOrder)
	}

	// Validate output format
	switch *outputFormat {
	case "tabular", "csv", "json", "json-lines":
		cla.OutputFormat = *outputFormat
	default:
		return fmt.Errorf("file format error: %s. Permitted formats: 'tabular', 'csv', 'json', 'json-lines'", *outputFormat)
	}

	return nil
}
