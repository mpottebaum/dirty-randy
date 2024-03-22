package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/andybalholm/cascadia"
	"golang.org/x/net/html"
)

func main() {
	args := os.Args
	if len(args) <= 1 {
		fmt.Println("Please enter a league name")
		return
	}
	yourLeague := args[1]
	url := "https://www.espn.com/" + yourLeague + "/schedule"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	content, err := io.ReadAll(resp.Body)
	fmt.Println("Resp status: ", resp.Status)
	// needed: date, time(localized?), event name/location
	// on something
	if err != nil {
		log.Fatal(err)
	}
	doc, _ := html.Parse(strings.NewReader(string(content)))
	rowSel, err := cascadia.Parse("tr.Table__TR")
	if err != nil {
		log.Fatal(err)
	}
	raceSel, err := cascadia.Parse("td.race__col")
	if err != nil {
		log.Fatal(err)
	}
	raceNameSel, err := cascadia.Parse("a")
	if err != nil {
		log.Fatal(err)
	}
	raceTrackSel, err := cascadia.Parse("div")
	if err != nil {
		log.Fatal(err)
	}
	rowsResult := cascadia.QueryAll(doc, rowSel)
	for i, row := range rowsResult {
		raceResult := cascadia.Query(row, raceSel)
		if raceResult != nil {
			raceNameResult := cascadia.Query(raceResult, raceNameSel)
			raceTrackResult := cascadia.Query(raceResult, raceTrackSel)
			//TODO: add nil result handling
			fmt.Println("race name: ", raceNameResult.FirstChild.Data)
			fmt.Println("race track: ", raceTrackResult.FirstChild.Data)
			fmt.Println("num: ", i)
		}
	}
	resp.Body.Close()
}
