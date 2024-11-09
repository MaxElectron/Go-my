package internal

import (
	"log"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type stats struct {
	userToLines      map[string]int
	userToCommits    map[string]map[string]bool
	userToNumCommits map[string]int
	userToFiles      map[string]map[string]bool
	userToNumFiles   map[string]int
	combinedData     map[string][3]int
	sortedData       [][4]string
}

var totalStats stats

func CountStatistics(fp *FilesParams) stats {
	totalStats = stats{
		userToLines:      make(map[string]int),
		userToCommits:    make(map[string]map[string]bool),
		userToNumCommits: make(map[string]int),
		userToFiles:      make(map[string]map[string]bool),
		userToNumFiles:   make(map[string]int),
		combinedData:     make(map[string][3]int),
	}

	for _, path := range fp.FilesList {
		processFile(path, *fp.Cla.RepositoryPath, *fp.Cla.CommitPointer, *fp.Cla.UseCommitter)
	}

	totalStats.combineResults()
	return totalStats
}

func addLine(author string) {
	totalStats.userToLines[author]++
}

func addCommit(author, commit string) {
	if _, ok := totalStats.userToCommits[author]; !ok {
		totalStats.userToCommits[author] = make(map[string]bool)
	}
	if _, ok := totalStats.userToCommits[author][commit]; !ok {
		totalStats.userToCommits[author][commit] = true
		totalStats.userToNumCommits[author]++
	}
}

func addFile(author, path string) {
	if _, ok := totalStats.userToFiles[author]; !ok {
		totalStats.userToFiles[author] = make(map[string]bool)
	}
	if _, ok := totalStats.userToFiles[author][path]; !ok {
		totalStats.userToFiles[author][path] = true
		totalStats.userToNumFiles[author]++
	}
}

func processFile(path, gitDir, commitPointer string, useCommiter bool) {
	// Execute git blame command
	gitBlameCmd := exec.Command("git", "blame", "--line-porcelain", "-b", commitPointer, path)
	gitBlameCmd.Dir = gitDir

	var gitBlameCmdOutput strings.Builder
	gitBlameCmd.Stdout = &gitBlameCmdOutput

	err := gitBlameCmd.Run()
	if err != nil {
		log.Fatalf("processFile: %v", err)
	}

	commitersLog := gitBlameCmdOutput.String()
	statLines := strings.Split(commitersLog, "\n")
	author := ""
	commitHash := ""

	if len(statLines) == 1 && statLines[0] == "" {
		// Empty file, execute git log command
		gitLogCmd := exec.Command("git", "log", "-p", commitPointer, "--follow", "--", path)
		gitLogCmd.Dir = gitDir
		var gitLogCmdOutput strings.Builder
		gitLogCmd.Stdout = &gitLogCmdOutput
		err = gitLogCmd.Run()
		if err != nil {
			log.Fatalf("processFile: %v", err)
		}

		gitLog := gitLogCmdOutput.String()
		logLines := strings.Split(gitLog, "\n")
		commitHash = strings.Split(logLines[0], " ")[1]
		words := strings.Split(logLines[1], " ")
		author = strings.Join(words[1:len(words)-1], " ")

		addCommit(author, commitHash)
		addFile(author, path)
	}

	for i := 0; i < len(statLines); i++ {
		words := strings.Split(statLines[i], " ")

		if useCommiter {
			if words[0] == "committer" {
				commitHash = strings.Split(statLines[i-5], " ")[0]
				author = strings.Join(words[1:], " ")
			} else {
				continue
			}
		} else {
			if words[0] == "author" {
				commitHash = strings.Split(statLines[i-1], " ")[0]
				author = strings.Join(words[1:], " ")
			} else {
				continue
			}
		}

		addLine(author)
		addCommit(author, commitHash)
		addFile(author, path)
	}
}

func (stats *stats) combineResults() {
	for name, numCommits := range stats.userToNumCommits {
		numLines := 0

		if actualNumLines, ok := stats.userToLines[name]; ok {
			numLines = actualNumLines
		}

		stats.combinedData[name] = [3]int{
			numLines,
			numCommits,
			stats.userToNumFiles[name],
		}
	}
}

func (stats *stats) SortResults(sortKey string) {
	var users []string
	for user := range stats.userToNumCommits {
		if user != "Not Committed Yet" {
			users = append(users, user)
		}
	}

	if sortKey == "lines" {
		sort.SliceStable(users, func(i, j int) bool {
			if stats.combinedData[users[i]][0] == stats.combinedData[users[j]][0] {
				if stats.combinedData[users[i]][1] == stats.combinedData[users[j]][1] {
					if stats.combinedData[users[i]][2] == stats.combinedData[users[j]][2] {
						return users[i] < users[j]
					}
					return stats.combinedData[users[i]][2] > stats.combinedData[users[j]][2]
				}
				return stats.combinedData[users[i]][1] > stats.combinedData[users[j]][1]
			}
			return stats.combinedData[users[i]][0] > stats.combinedData[users[j]][0]
		})
	} else if sortKey == "commits" {
		sort.SliceStable(users, func(i, j int) bool {
			if stats.combinedData[users[i]][1] == stats.combinedData[users[j]][1] {
				if stats.combinedData[users[i]][0] == stats.combinedData[users[j]][0] {
					if stats.combinedData[users[i]][2] == stats.combinedData[users[j]][2] {
						return users[i] < users[j]
					}
					return stats.combinedData[users[i]][2] > stats.combinedData[users[j]][2]
				}
				return stats.combinedData[users[i]][0] > stats.combinedData[users[j]][0]
			}
			return stats.combinedData[users[i]][1] > stats.combinedData[users[j]][1]
		})
	} else if sortKey == "files" {
		sort.SliceStable(users, func(i, j int) bool {
			if stats.combinedData[users[i]][2] == stats.combinedData[users[j]][2] {
				if stats.combinedData[users[i]][0] == stats.combinedData[users[j]][0] {
					if stats.combinedData[users[i]][1] == stats.combinedData[users[j]][1] {
						return users[i] < users[j]
					}
					return stats.combinedData[users[i]][1] > stats.combinedData[users[j]][1]
				}
				return stats.combinedData[users[i]][0] > stats.combinedData[users[j]][0]
			}
			return stats.combinedData[users[i]][2] > stats.combinedData[users[j]][2]
		})
	}

	var sortedStats [][4]string

	for _, user := range users {
		sortedStats = append(sortedStats, [4]string{user, strconv.Itoa(stats.combinedData[user][0]), strconv.Itoa(stats.combinedData[user][1]), strconv.Itoa(stats.combinedData[user][2])})
	}
	stats.sortedData = sortedStats
}
