package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
)

type Athlete struct {
	Date         string `json:"date"`
	Sport        string `json:"sport"`
	GoldMedals   int    `json:"gold"`
	Country      string `json:"country"`
	Year         int    `json:"year"`
	Name         string `json:"athlete"`
	Age          int    `json:"age"`
	SilverMedals int    `json:"silver"`
	BronzeMedals int    `json:"bronze"`
	TotalMedals  int    `json:"total"`
}

type MedalCount struct {
	Silver int `json:"silver"`
	Bronze int `json:"bronze"`
	Gold   int `json:"gold"`
	Total  int `json:"total"`
}

type AthleteMedals struct {
	Country           string                `json:"country"`
	MedalCounts       MedalCount            `json:"medals"`
	Name              string                `json:"athlete"`
	MedalCountsByYear map[string]MedalCount `json:"medals_by_year"`
}

type CountryMedals struct {
	SilverMedals int    `json:"silver"`
	BronzeMedals int    `json:"bronze"`
	Country      string `json:"country"`
	GoldMedals   int    `json:"gold"`
	TotalMedals  int    `json:"total"`
}

type OlympicsHandler struct {
	mux *http.ServeMux

	nameToCountryMap      map[string]string
	nameToMedalCount      map[string]MedalCount
	nameToYearMedals      map[string]map[string]MedalCount
	nameToSportMap        map[string][]string
	nameToSportMedals     map[string]map[string]MedalCount
	nameToSportYearMedals map[string]map[string]map[string]MedalCount
	yearToMedalCount      map[string]map[string]MedalCount
	yearToCountries       map[string][]string
}

func main() {
	httpPort := flag.String("port", "", "HTTP port")
	dataFilePath := flag.String("data", "", "Data file path")
	flag.Parse()
	if err := Begin(*httpPort, *dataFilePath); err != nil {
		panic(err)
	}
}

func Begin(port, dataPath string) error {
	serverMux := http.NewServeMux()
	olympicsHandler := OlympicsHandler{
		mux:                   serverMux,
		nameToCountryMap:      make(map[string]string),
		nameToMedalCount:      make(map[string]MedalCount),
		nameToYearMedals:      make(map[string]map[string]MedalCount),
		nameToSportMap:        make(map[string][]string),
		nameToSportMedals:     make(map[string]map[string]MedalCount),
		nameToSportYearMedals: make(map[string]map[string]map[string]MedalCount),
		yearToMedalCount:      make(map[string]map[string]MedalCount),
		yearToCountries:       make(map[string][]string),
	}

	dataFile, err := os.Open(dataPath)
	if err != nil {
		return fmt.Errorf("can't open data file: %v", err)
	}
	defer dataFile.Close()

	dataContent, err := io.ReadAll(dataFile)
	if err != nil {
		return fmt.Errorf("can't read data file: %v", err)
	}

	var athletes []Athlete
	if err := json.Unmarshal(dataContent, &athletes); err != nil {
		return fmt.Errorf("can't unmarshal json: %v", err)
	}

	for _, athlete := range athletes {
		if _, exists := olympicsHandler.nameToCountryMap[athlete.Name]; !exists {
			olympicsHandler.nameToCountryMap[athlete.Name] = athlete.Country
		}

		currentMedalCount := olympicsHandler.nameToMedalCount[athlete.Name]
		currentMedalCount.Gold += athlete.GoldMedals
		currentMedalCount.Silver += athlete.SilverMedals
		currentMedalCount.Bronze += athlete.BronzeMedals
		currentMedalCount.Total += athlete.TotalMedals
		olympicsHandler.nameToMedalCount[athlete.Name] = currentMedalCount

		if _, exists := olympicsHandler.nameToYearMedals[athlete.Name]; !exists {
			olympicsHandler.nameToYearMedals[athlete.Name] = make(map[string]MedalCount)
		}
		yearMedalCount := olympicsHandler.nameToYearMedals[athlete.Name][strconv.Itoa(athlete.Year)]
		yearMedalCount.Gold += athlete.GoldMedals
		yearMedalCount.Silver += athlete.SilverMedals
		yearMedalCount.Bronze += athlete.BronzeMedals
		yearMedalCount.Total += athlete.TotalMedals
		olympicsHandler.nameToYearMedals[athlete.Name][strconv.Itoa(athlete.Year)] = yearMedalCount

		if _, exists := olympicsHandler.nameToSportMap[athlete.Sport]; !exists {
			olympicsHandler.nameToSportMap[athlete.Sport] = []string{}
		}
		olympicsHandler.nameToSportMap[athlete.Sport] = append(olympicsHandler.nameToSportMap[athlete.Sport], athlete.Name)

		if _, exists := olympicsHandler.nameToSportMedals[athlete.Name]; !exists {
			olympicsHandler.nameToSportMedals[athlete.Name] = make(map[string]MedalCount)
		}
		currentSportMedalCount := olympicsHandler.nameToSportMedals[athlete.Name][athlete.Sport]
		currentSportMedalCount.Gold += athlete.GoldMedals
		currentSportMedalCount.Silver += athlete.SilverMedals
		currentSportMedalCount.Bronze += athlete.BronzeMedals
		currentSportMedalCount.Total += athlete.TotalMedals
		olympicsHandler.nameToSportMedals[athlete.Name][athlete.Sport] = currentSportMedalCount

		if _, exists := olympicsHandler.nameToSportYearMedals[athlete.Name]; !exists {
			olympicsHandler.nameToSportYearMedals[athlete.Name] = make(map[string]map[string]MedalCount)
		}
		if _, exists := olympicsHandler.nameToSportYearMedals[athlete.Name][athlete.Sport]; !exists {
			olympicsHandler.nameToSportYearMedals[athlete.Name][athlete.Sport] = make(map[string]MedalCount)
		}
		sportYearMedalCount := olympicsHandler.nameToSportYearMedals[athlete.Name][athlete.Sport][strconv.Itoa(athlete.Year)]
		sportYearMedalCount.Gold += athlete.GoldMedals
		sportYearMedalCount.Silver += athlete.SilverMedals
		sportYearMedalCount.Bronze += athlete.BronzeMedals
		sportYearMedalCount.Total += athlete.TotalMedals
		olympicsHandler.nameToSportYearMedals[athlete.Name][athlete.Sport][strconv.Itoa(athlete.Year)] = sportYearMedalCount

		if _, exists := olympicsHandler.yearToMedalCount[strconv.Itoa(athlete.Year)]; !exists {
			olympicsHandler.yearToMedalCount[strconv.Itoa(athlete.Year)] = make(map[string]MedalCount)
		}
		countryMedalCount := olympicsHandler.yearToMedalCount[strconv.Itoa(athlete.Year)][athlete.Country]
		countryMedalCount.Gold += athlete.GoldMedals
		countryMedalCount.Silver += athlete.SilverMedals
		countryMedalCount.Bronze += athlete.BronzeMedals
		countryMedalCount.Total += athlete.TotalMedals
		olympicsHandler.yearToMedalCount[strconv.Itoa(athlete.Year)][athlete.Country] = countryMedalCount

		if _, exists := olympicsHandler.yearToCountries[strconv.Itoa(athlete.Year)]; !exists {
			olympicsHandler.yearToCountries[strconv.Itoa(athlete.Year)] = []string{}
		}
		olympicsHandler.yearToCountries[strconv.Itoa(athlete.Year)] = append(olympicsHandler.yearToCountries[strconv.Itoa(athlete.Year)], athlete.Country)
	}

	for sport := range olympicsHandler.nameToSportMap {
		uniqueAthleteNames := make(map[string]bool)
		filteredCount := 0

		for _, athleteName := range olympicsHandler.nameToSportMap[sport] {
			if _, exists := uniqueAthleteNames[athleteName]; !exists {
				uniqueAthleteNames[athleteName] = true
				olympicsHandler.nameToSportMap[sport][filteredCount] = athleteName
				filteredCount++
			}
		}

		olympicsHandler.nameToSportMap[sport] = olympicsHandler.nameToSportMap[sport][:filteredCount]
		sort.Slice(olympicsHandler.nameToSportMap[sport], func(i, j int) bool {
			athleteA := olympicsHandler.nameToSportMap[sport][i]
			athleteB := olympicsHandler.nameToSportMap[sport][j]
			if olympicsHandler.nameToSportMedals[athleteA][sport].Gold == olympicsHandler.nameToSportMedals[athleteB][sport].Gold {
				if olympicsHandler.nameToSportMedals[athleteA][sport].Silver == olympicsHandler.nameToSportMedals[athleteB][sport].Silver {
					if olympicsHandler.nameToSportMedals[athleteA][sport].Bronze == olympicsHandler.nameToSportMedals[athleteB][sport].Bronze {
						return athleteA < athleteB
					}
					return olympicsHandler.nameToSportMedals[athleteA][sport].Bronze > olympicsHandler.nameToSportMedals[athleteB][sport].Bronze
				}
				return olympicsHandler.nameToSportMedals[athleteA][sport].Silver > olympicsHandler.nameToSportMedals[athleteB][sport].Silver
			}
			return olympicsHandler.nameToSportMedals[athleteA][sport].Gold > olympicsHandler.nameToSportMedals[athleteB][sport].Gold
		})
	}

	for year := range olympicsHandler.yearToCountries {
		uniqueCountries := make(map[string]bool)
		filteredCount := 0
		for _, country := range olympicsHandler.yearToCountries[year] {
			if _, exists := uniqueCountries[country]; !exists {
				uniqueCountries[country] = true
				olympicsHandler.yearToCountries[year][filteredCount] = country
				filteredCount++
			}
		}

		olympicsHandler.yearToCountries[year] = olympicsHandler.yearToCountries[year][:filteredCount]
		sort.Slice(olympicsHandler.yearToCountries[year], func(i, j int) bool {
			countryA := olympicsHandler.yearToCountries[year][i]
			countryB := olympicsHandler.yearToCountries[year][j]
			if olympicsHandler.yearToMedalCount[year][countryA].Gold == olympicsHandler.yearToMedalCount[year][countryB].Gold {
				if olympicsHandler.yearToMedalCount[year][countryA].Silver == olympicsHandler.yearToMedalCount[year][countryB].Silver {
					if olympicsHandler.yearToMedalCount[year][countryA].Bronze == olympicsHandler.yearToMedalCount[year][countryB].Bronze {
						return countryA < countryB
					}
					return olympicsHandler.yearToMedalCount[year][countryA].Bronze > olympicsHandler.yearToMedalCount[year][countryB].Bronze
				}
				return olympicsHandler.yearToMedalCount[year][countryA].Silver > olympicsHandler.yearToMedalCount[year][countryB].Silver
			}
			return olympicsHandler.yearToMedalCount[year][countryA].Gold > olympicsHandler.yearToMedalCount[year][countryB].Gold
		})
	}

	serverMux.HandleFunc("/athlete-info", olympicsHandler.athleteInfoHandler)
	serverMux.HandleFunc("/top-athletes-in-sport", olympicsHandler.topAthletesInSportHandler)
	serverMux.HandleFunc("/top-countries-in-year", olympicsHandler.topCountriesInYearHandler)

	return http.ListenAndServe(":"+port, serverMux)
}

func (handler *OlympicsHandler) athleteInfoHandler(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Add("Content-Type", "application/json")

	parsedURL, err := url.Parse(req.URL.String())
	if err != nil {
		http.Error(writer, "URL parsing error", http.StatusBadRequest)
		return
	}

	names := parsedURL.Query()["name"]
	if len(names) != 1 {
		http.Error(writer, "Invalid request: exactly one name is required", http.StatusBadRequest)
		return
	}

	name := names[0]
	if _, exists := handler.nameToCountryMap[name]; !exists {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	athleteMedalInfo := &AthleteMedals{
		Name:              name,
		Country:           handler.nameToCountryMap[name],
		MedalCounts:       handler.nameToMedalCount[name],
		MedalCountsByYear: handler.nameToYearMedals[name],
	}

	response, err := json.Marshal(athleteMedalInfo)
	if err != nil {
		http.Error(writer, "Error marshaling response", http.StatusInternalServerError)
		return
	}

	writer.Write(response)
}

func (handler *OlympicsHandler) topAthletesInSportHandler(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Add("Content-Type", "application/json")

	parsedURL, err := url.Parse(req.URL.String())
	if err != nil {
		http.Error(writer, "URL parsing error", http.StatusBadRequest)
		return
	}

	sports := parsedURL.Query()["sport"]
	if len(sports) != 1 {
		http.Error(writer, "Invalid request: exactly one sport is required", http.StatusBadRequest)
		return
	}

	sport := sports[0]
	if _, exists := handler.nameToSportMap[sport]; !exists {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	actualLimit := 3
	limits := parsedURL.Query()["limit"]

	if len(limits) == 1 {
		actualLimit, err = strconv.Atoi(limits[0])
		if err != nil {
			http.Error(writer, "Invalid limit value", http.StatusBadRequest)
			return
		}
	}

	var athleteMedals []AthleteMedals
	for lim := 0; lim < actualLimit && lim < len(handler.nameToSportMap[sport]); lim++ {
		athleteName := handler.nameToSportMap[sport][lim]
		athleteMedals = append(athleteMedals, AthleteMedals{
			Name:              athleteName,
			Country:           handler.nameToCountryMap[athleteName],
			MedalCounts:       handler.nameToSportMedals[athleteName][sport],
			MedalCountsByYear: handler.nameToSportYearMedals[athleteName][sport],
		})
	}

	response, err := json.Marshal(athleteMedals)
	if err != nil {
		http.Error(writer, "Error marshaling response", http.StatusInternalServerError)
		return
	}

	writer.Write(response)
}

func (handler *OlympicsHandler) topCountriesInYearHandler(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Add("Content-Type", "application/json")

	parsedURL, err := url.Parse(req.URL.String())
	if err != nil {
		http.Error(writer, "URL parsing error", http.StatusBadRequest)
		return
	}

	years := parsedURL.Query()["year"]
	if len(years) != 1 {
		http.Error(writer, "Invalid request: exactly one year is required", http.StatusBadRequest)
		return
	}

	year := years[0]
	if _, exists := handler.yearToCountries[year]; !exists {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	currentLimit := 3
	limits := parsedURL.Query()["limit"]

	if len(limits) == 1 {
		currentLimit, err = strconv.Atoi(limits[0])
		if err != nil {
			http.Error(writer, "Invalid limit value", http.StatusBadRequest)
			return
		}
	}

	var countryMedals []CountryMedals
	for i := 0; i < currentLimit && i < len(handler.yearToCountries[year]); i++ {
		country := handler.yearToCountries[year][i]
		countryMedals = append(countryMedals, CountryMedals{
			Country:      country,
			GoldMedals:   handler.yearToMedalCount[year][country].Gold,
			SilverMedals: handler.yearToMedalCount[year][country].Silver,
			BronzeMedals: handler.yearToMedalCount[year][country].Bronze,
			TotalMedals:  handler.yearToMedalCount[year][country].Total,
		})
	}

	response, err := json.Marshal(countryMedals)
	if err != nil {
		http.Error(writer, "Error marshaling response", http.StatusInternalServerError)
		return
	}

	writer.Write(response)
}
