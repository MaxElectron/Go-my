//go:build !solution

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func urlFetch(link string, channel chan<- string) {
	startTime := time.Now()

	response, err := http.Get(link)

	if err != nil {
		channel <- fmt.Sprint(err)
		return
	}

	byteCount, err := io.Copy(io.Discard, response.Body)
	response.Body.Close()

	if err != nil {
		channel <- fmt.Sprint(err)
		return
	}

	channel <- fmt.Sprintf("%.2fs%9d  %s", time.Since(startTime).Seconds(), byteCount, link)
}

func main() {
	links := os.Args[1:]

	startTime := time.Now()
	channel := make(chan string)

	for _, link := range links {
		go urlFetch(link, channel)
	}

	for range links {
		fmt.Println(<-channel)
	}

	fmt.Printf("%.2fs elapsed\n", time.Since(startTime).Seconds())
}
