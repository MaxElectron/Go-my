package internal

import (
	"log"
	"os/exec"
	"path/filepath"
	"strings"
)

type MappingEntity struct {
	Name       string
	Type       string
	Extensions []string
}

type FilesParams struct {
	FilesList []string
	Cla       *CommandLineArgs
	Mapping   []MappingEntity
}

func NewFilesParams(mapping []MappingEntity, cla *CommandLineArgs) *FilesParams {
	return &FilesParams{Cla: cla, Mapping: mapping}
}

func (fp *FilesParams) GetAllFiles(commitPointer, gitDir string) {
	// Execute git ls-tree command
	gitLsTreeCmd := exec.Command("git", "ls-tree", "-r", "--name-only", commitPointer)
	gitLsTreeCmd.Dir = gitDir

	var gitLsTreeCmdOutput strings.Builder
	gitLsTreeCmd.Stdout = &gitLsTreeCmdOutput

	err := gitLsTreeCmd.Run()
	if err != nil {
		log.Fatalf("Get all files %v", err)
	}

	gitTree := gitLsTreeCmdOutput.String()
	filesInfo := strings.Split(gitTree, "\n")
	for _, file := range filesInfo {
		if file == "" {
			continue
		}

		// Check file extension
		if len(*fp.Cla.FileExtensions) > 0 {
			ext := filepath.Ext(file)
			extensionMatch := false
			for _, e := range *fp.Cla.FileExtensions {
				if strings.EqualFold(ext, e) {
					extensionMatch = true
					break
				}
			}
			if !extensionMatch {
				continue
			}
		}

		// Determine file language
		fileLanguage := ""
		fileExtension := filepath.Ext(file)
		for _, mappingEntity := range fp.Mapping {
			for _, extension := range mappingEntity.Extensions {
				if strings.EqualFold(fileExtension, extension) {
					fileLanguage = mappingEntity.Name
					break
				}
			}
			if fileLanguage != "" {
				break
			}
		}

		// Check if language is acceptable
		if len(*fp.Cla.ProgrammingLanguages) > 0 {
			languageMatch := false
			for _, language := range *fp.Cla.ProgrammingLanguages {
				if strings.EqualFold(language, fileLanguage) {
					languageMatch = true
					break
				}
			}
			if !languageMatch {
				continue
			}
		}

		// Check exclude patterns
		if len(*fp.Cla.ExcludePatterns) > 0 {
			excludeMatch := false
			for _, pattern := range *fp.Cla.ExcludePatterns {
				match, _ := filepath.Match(pattern, file)
				if match {
					excludeMatch = true
					break
				}
			}
			if excludeMatch {
				continue
			}
		}

		// Check include patterns
		if len(*fp.Cla.IncludePatterns) > 0 {
			includeMatch := false
			for _, pattern := range *fp.Cla.IncludePatterns {
				match, _ := filepath.Match(pattern, file)
				if match {
					includeMatch = true
					break
				}
			}
			if !includeMatch {
				continue
			}
		}

		fp.FilesList = append(fp.FilesList, file)
	}
}
