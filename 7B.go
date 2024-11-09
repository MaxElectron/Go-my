//go:build !solution

package main

import (
	"fmt"
	"image"
	"image/color"
	res_png "image/png"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type TimeService struct {
	mux *http.ServeMux
}

type Response struct {
	URL string `json:"url"`
	Key string `json:"key"`
}

func matchTime(tstr string) bool {
	if iserr, _ := regexp.MatchString("[0-2][0-9]:[0-5][0-9]:[0-5][0-9]", tstr); !iserr {
		return false
	}
	hours, _ := strconv.Atoi(tstr[0:2])
	minutes, _ := strconv.Atoi(tstr[3:5])
	seconds, _ := strconv.Atoi(tstr[6:8])
	return hours < 24 && seconds < 60 && minutes < 60
}

func (sv *TimeService) DisplayTimeAsImage(wr http.ResponseWriter, rq *http.Request) {
	wr.Header().Add("Content-Type", "image/png")
	uro, _ := url.Parse(rq.URL.String())
	line := ""
	if len(uro.Query()["time"]) > 0 {
		line = uro.Query()["time"][0]
	}
	if line == "" {
		time_now := time.Now()
		line = fmt.Sprintf("%02d:%02d:%02d", time_now.Hour(), time_now.Minute(), time_now.Second())
	}
	if !matchTime(line) {
		wr.WriteHeader(http.StatusBadRequest)
		return
	}
	actual := 0
	if actual, _ = strconv.Atoi(uro.Query()["k"][0]); actual < 1 || actual > 30 {
		wr.WriteHeader(http.StatusBadRequest)
		return
	}
	x_coord, y_coord := 0, 0

	img := image.NewRGBA(image.Rect(0, 0, actual*56, actual*12))
	for i := 0; i < 8; i++ {
		// setting symbols
		allsym := make(map[byte]string)

		allsym[':'] = Colon
		allsym['0'] = Zero
		allsym['1'] = One
		allsym['2'] = Two
		allsym['3'] = Three
		allsym['4'] = Four
		allsym['5'] = Five
		allsym['6'] = Six
		allsym['7'] = Seven
		allsym['8'] = Eight
		allsym['9'] = Nine
		xx := x_coord
		yy := y_coord
		sym := allsym[line[i]] + "\n"
		new := ""
		line := ""
		for _, symb := range sym {
			if symb == '\n' {
				for i := 0; i < actual; i++ {
					new += line + "\n"
				}

				line = ""
			} else {
				line += strings.Repeat(string(symb), actual)
			}

		}
		dist := 0
		for _, cc := range new {
			if cc == '.' {
				img.Set(xx, yy, color.RGBA{255, 255, 255, 255})

				dist++

				xx++
			} else if cc == '1' {
				img.Set(xx, yy, Cyan)

				dist++

				xx++
			} else {
				xx -= dist

				dist = 0

				yy++
			}

		}

		if !(i == 5 || i == 2) {
			x_coord += 8 * actual
		} else {
			x_coord += 4 * actual
		}
	}

	output, _ := os.Create("/tmp/image.png")
	res_png.Encode(output, img)
	rd, _ := os.ReadFile("/tmp/image.png")
	wr.Write(rd)
}

func main() {
	service := TimeService{
		mux: http.NewServeMux(),
	}
	service.mux.HandleFunc("/", service.DisplayTimeAsImage)
	log.Fatal(http.ListenAndServe(":"+os.Args[2], service.mux))
}
