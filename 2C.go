//go:build !solution

package externalsort

import (
	"bytes"
	"container/heap"
	"errors"
	"io"
	"os"
	"sort"
	"strings"
)

type MyReader struct {
	r io.Reader
}

type MyWriter struct {
	w io.Writer
}

func NewReader(r io.Reader) LineReader {
	return &MyReader{r}
}

func NewWriter(w io.Writer) LineWriter {
	return &MyWriter{w}
}

// structs for mergesort
type Line struct {
	value  string
	reader LineReader
}

type LineHeap []*Line

// interface methods required by container/heap
func (h LineHeap) Len() int           { return len(h) }
func (h LineHeap) Less(i, j int) bool { return h[i].value < h[j].value }
func (h LineHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *LineHeap) Push(x interface{}) {
	*h = append(*h, x.(*Line))
}
func (h *LineHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func Merge(w LineWriter, readers ...LineReader) error {
	// read lines
	var lines LineHeap

	for _, r := range readers {
		line, err := r.ReadLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				if line != "" {
					lines = append(lines, &Line{value: line, reader: r})
				}
				continue
			}
			return err
		}
		lines = append(lines, &Line{value: line, reader: r})
	}

	// heapify
	heap.Init(&lines)

	// write lines
	for len(lines) > 0 {
		line := heap.Pop(&lines).(*Line)
		if err := w.Write(line.value); err != nil {
			return err
		}
		nextLine, err := line.reader.ReadLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				if nextLine != "" {
					heap.Push(&lines, &Line{value: nextLine, reader: line.reader})
				}
				continue
			}
			return err
		}
		heap.Push(&lines, &Line{value: nextLine, reader: line.reader})
	}

	return nil
}

func Sort(w io.Writer, in ...string) error {
	// create readers for each file
	readers := make([]LineReader, len(in))

	// assign readers to files
	for i, filename := range in {
		// sort file
		if err := sortFile(filename); err != nil {
			return err
		}

		// open file
		f, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer f.Close()

		// assign reader
		readers[i] = NewReader(f)
	}

	// offload to Merge
	return Merge(NewWriter(w), readers...)
}

func sortFile(filename string) error {
	// open file for reading and writing
	f, err := os.OpenFile(filename, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// read the file
	lines, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	// sort the lines
	lineSlice := strings.Split(string(lines), "\n")
	sort.Strings(lineSlice)

	// truncate the file
	err = f.Truncate(0)
	if err != nil {
		return err
	}

	// rewind the file
	_, err = f.Seek(0, 0)
	if err != nil {
		return err
	}

	// some newline count magic here, because the newlines are weird
	if len(string(lines)) > 0 && string(lines)[len(string(lines))-1] == '\n' && lineSlice[len(lineSlice)-1] != "" {
		_, err = f.Write([]byte(strings.Join(lineSlice, "\n"))[1:])
	} else {
		_, err = f.Write([]byte(strings.Join(lineSlice, "\n")))
	}
	return err
}

func (r *MyReader) ReadLine() (string, error) {
	// make buffer to store the line in
	var buf bytes.Buffer

	// read until the eof
	for {
		b := make([]byte, 1)
		n, err := r.r.Read(b)

		// check for errors and EOF
		if err != nil {
			if errors.Is(err, io.EOF) {
				if n != 0 {
					buf.WriteByte(b[0])
				}
				return buf.String(), io.EOF
			}
			return "", err
		}

		// check for new line
		if b[0] == '\n' {
			return buf.String(), nil
		}

		// add to buffer
		buf.WriteByte(b[0])
	}
}

func (w *MyWriter) Write(l string) error {
	// attempt writing
	_, err := w.w.Write([]byte(l))
	if err != nil {
		return err
	}

	// write new line
	_, err = w.w.Write([]byte{'\n'})
	return err
}
