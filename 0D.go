//go:build !solution

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func urlFetch(link string) string {
	response, err := http.Get(link)

	if err != nil {
		panic(err)
	}

	body, err := io.ReadAll(response.Body)
	response.Body.Close()

	if err != nil {
		panic(err)
	}

	return string(body)
}

func main() {
	links := os.Args[1:]

	for _, link := range links {
		fmt.Printf("%s", urlFetch(link))
	}
}
