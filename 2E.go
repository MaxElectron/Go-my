//go:build !solution

package ciletters

import (
	"bytes"
	"strings"
	"text/template"
)

// its 9 am after a sleepless night, I won't cleanup that code
func truncate(s string, length int) string {
	if length == 0 {
		return s
	}
	if len(s) < length {
		return s
	}
	return s[:length]
}

func present(s string, indent int) string {
	var result string
	result += "\n"
	stringsSlice := strings.Split(s, "\n")
	if len(stringsSlice) >= 10 {
		stringsSlice = stringsSlice[len(stringsSlice)-10:]
	}
	for _, line := range stringsSlice {
		result += strings.Repeat(" ", indent) + line + "\n"
	}
	return result
}

func MakeLetter(n *Notification) (string, error) {
	letterTemplate := `Your pipeline #{{ .Pipeline.ID }} {{ if isFailed .Pipeline.Status }}has failed{{ else }}passed{{ end }}!
    Project:      {{ .Project.GroupID }}/{{ .Project.ID }}
    Branch:       🌿 {{ .Branch }}
    Commit:       {{ truncate .Commit.Hash 8 }} {{ truncate .Commit.Message 0}}
    CommitAuthor: {{ .Commit.Author}}{{ range .Pipeline.FailedJobs }}
        Stage: {{ .Stage }}, Job {{ .Name }}{{ if .RunnerLog }}{{ present .RunnerLog 12 }}{{ end }}{{ end }}`

	funcMap := template.FuncMap{
		"truncate": truncate,
		"present":  present,
		"isFailed": func(status PipelineStatus) bool { return status == PipelineStatusFailed },
	}
	tmpl := template.Must(template.New("letter").Funcs(funcMap).Parse(letterTemplate))
	buf := bytes.NewBufferString("")
	err := tmpl.Execute(buf, n)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
