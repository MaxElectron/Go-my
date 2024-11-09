package internal

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"text/tabwriter"
)

func (stats *stats) Print(format string) {
	switch format {
	case "tabular":
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
		_, err := fmt.Fprintln(w, "Name\tLines\tCommits\tFiles")
		if err != nil {
			log.Fatalf("tabular: %v", err)
		}

		for _, line := range stats.sortedData {
			_, err = fmt.Fprintln(w, line[0]+"\t"+line[1]+"\t"+line[2]+"\t"+line[3])
			if err != nil {
				log.Fatalf("tabular: %v", err)
			}
		}

		err = w.Flush()
		if err != nil {
			log.Fatalf("tabular: %v", err)
		}

	case "csv":
		header := []string{"Name", "Lines", "Commits", "Files"}
		w := csv.NewWriter(os.Stdout)
		var buff [][]string
		buff = append(buff, header)
		for _, line := range stats.sortedData {
			buff = append(buff, []string{line[0], line[1], line[2], line[3]})
		}
		err := w.WriteAll(buff)
		if err != nil {
			log.Fatalf("csv: %v", err)
		}

	case "json":
		var buff []map[string]interface{}
		for _, line := range stats.sortedData {
			lines, err := strconv.Atoi(line[1])
			if err != nil {
				log.Fatalf("json: could not convert num of lines: %v", err)
			}

			commits, err := strconv.Atoi(line[2])
			if err != nil {
				log.Fatalf("json: could not convert num of commits: %v", err)
			}

			files, err := strconv.Atoi(line[3])
			if err != nil {
				log.Fatalf("json: could not convert num of files: %v", err)
			}

			buff = append(buff, map[string]interface{}{
				"name":    line[0],
				"lines":   lines,
				"commits": commits,
				"files":   files,
			})
		}
		jsonData, err := json.Marshal(buff)
		if err != nil {
			log.Fatalf("json: could not marshal json: %v", err)
		}

		fmt.Println(string(jsonData))

	case "json-lines":
		for _, line := range stats.sortedData {
			lines, err := strconv.Atoi(line[1])
			if err != nil {
				log.Fatalf("json-lines: could not convert num of lines: %v", err)
			}

			commits, err := strconv.Atoi(line[2])
			if err != nil {
				log.Fatalf("json-lines: could not convert num of commits: %v", err)
			}

			files, err := strconv.Atoi(line[3])
			if err != nil {
				log.Fatalf("json-lines: could not convert num of files: %v", err)
			}

			jsonLine, err := json.Marshal(map[string]interface{}{
				"name":    line[0],
				"lines":   lines,
				"commits": commits,
				"files":   files,
			})
			if err != nil {
				log.Fatalf("json-lines: could not marshal json: %v", err)
			}

			fmt.Println(string(jsonLine))
		}

	default:
		log.Fatalf("Print: unsupported format %s", format)
	}
}
