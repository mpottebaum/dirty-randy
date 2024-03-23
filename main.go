package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/andybalholm/cascadia"
	"golang.org/x/net/html"
)

var MonthMap = map[string]string{
	"Jan": "01",
	"Feb": "02",
	"Mar": "03",
	"Apr": "04",
	"May": "05",
	"Jun": "06",
	"Jul": "07",
	"Aug": "08",
	"Sep": "09",
	"Oct": "10",
	"Nov": "11",
	"Dec": "12",
}

type LeagueEvent struct {
	Name, Location, Date, Time, TV string
}

func ParseInt(str string) (i int, err error) {
	parsedInt, err := strconv.ParseInt(str, 10, 64)
	i = int(parsedInt)
	return
}

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
	lightsOutSel, err := cascadia.Parse("td.winnerLightsOut__col")
	if err != nil {
		log.Fatal(err)
	}
	dateTimeSel, err := cascadia.Parse("span")
	if err != nil {
		log.Fatal(err)
	}
	tvSel, err := cascadia.Parse("td.tv__col")
	if err != nil {
		log.Fatal(err)
	}
	rowsResult := cascadia.QueryAll(doc, rowSel)
	events := make([]LeagueEvent, len(rowsResult))
	for i, row := range rowsResult {
		events[i] = LeagueEvent{}
		// name and location
		raceResult := cascadia.Query(row, raceSel)
		if raceResult != nil {
			raceNameResult := cascadia.Query(raceResult, raceNameSel)
			raceTrackResult := cascadia.Query(raceResult, raceTrackSel)
			//TODO: add nil result handling
			events[i].Name = raceNameResult.FirstChild.Data
			events[i].Location = raceTrackResult.FirstChild.Data
		}
		// date and time
		lightsOutResult := cascadia.Query(row, lightsOutSel)
		if lightsOutResult != nil {
			dateTimeResult := cascadia.Query(lightsOutResult, dateTimeSel)
			//if race in past, there's no span (just winner's name)
			if dateTimeResult != nil {
				//TODO: add nil result handling
				rawDateTimeStr := dateTimeResult.FirstChild.Data
				dateAndTime := strings.Split(rawDateTimeStr, " - ")
				monthAndDay := strings.Split(dateAndTime[0], " ")
				month := monthAndDay[0]
				formattedMonth := MonthMap[month]
				day := monthAndDay[1]
				formattedDay := day
				if len(day) == 1 {
					formattedDay = "0" + formattedDay
				}
				hourAndPeriod := strings.Split(dateAndTime[1], " ")
				hour := hourAndPeriod[0]
				period := hourAndPeriod[1]
				formattedPeriod := strings.ToUpper(period)
				yearInt := time.Now().Year()
				formattedDate := formattedMonth + "/" + formattedDay + "/" + strconv.Itoa(yearInt)
				formattedTime := hour + " " + formattedPeriod
				events[i].Date = formattedDate
				events[i].Time = formattedTime
			}
		}
		tvResult := cascadia.Query(row, tvSel)
		if tvResult != nil && tvResult.FirstChild != nil {
			events[i].TV = tvResult.FirstChild.Data
		}
	}
	fmt.Println("events: ", events)
	resp.Body.Close()
}
