//go:build !solution

package main

import (
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
)

type UrlShortener struct {
	connections     map[string]string
	allowed_servers map[string]string

	mux *http.ServeMux
}

func (sv *UrlShortener) MakeShort(wr http.ResponseWriter, req *http.Request) {
	wr.Header().Add("Content-Type", "application/json")
	text, _ := io.ReadAll(req.Body)
	rp := make(map[string]string)
	if json.Unmarshal(text, &rp) != nil {
		wr.WriteHeader(http.StatusBadRequest)
		return
	}
	random_number := strconv.Itoa(rand.Int())
	if _, alr_used := sv.connections[rp["url"]]; !alr_used {
		sv.connections[rp["url"]] = random_number
		sv.allowed_servers[random_number] = rp["url"]
	} else {
		random_number = sv.connections[rp["url"]]
	}
	rp["key"] = random_number
	serv_response, _ := json.Marshal(&rp)
	wr.Write(serv_response)
}

func (sv *UrlShortener) Go(wr http.ResponseWriter, rq *http.Request) {
	wr.Header().Add("Content-Type", "application/json")
	if url, alr_here := sv.allowed_servers[rq.URL.String()[4:]]; alr_here {
		http.Redirect(wr, rq, url, http.StatusFound)
		wr.Header().Add("Location", url)
	} else {
		wr.WriteHeader(http.StatusNotFound)
	}
}

func main() {
	shortener := UrlShortener{
		connections:     make(map[string]string),
		mux:             http.NewServeMux(),
		allowed_servers: make(map[string]string),
	}

	shortener.mux.HandleFunc("/shorten", shortener.MakeShort)
	shortener.mux.HandleFunc("/go/", shortener.Go)
	log.Fatal(http.ListenAndServe(":"+os.Args[2], shortener.mux))
}
