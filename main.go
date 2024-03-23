package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
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

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func NewSel(selector string) (cascadia.Sel, error) {
	return cascadia.Parse(selector)
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
	check(err)
	defer resp.Body.Close()
	content, err := io.ReadAll(resp.Body)
	fmt.Println("Resp status: ", resp.Status)
	// needed: date, time(localized?), event name/location
	// on something
	check(err)
	doc, _ := html.Parse(strings.NewReader(string(content)))
	rowSel, err := NewSel("tr.Table__TR")
	check(err)
	raceSel, err := NewSel("td.race__col")
	check(err)
	raceNameSel, err := NewSel("a")
	check(err)
	raceTrackSel, err := NewSel("div")
	check(err)
	lightsOutSel, err := NewSel("td.winnerLightsOut__col")
	check(err)
	dateTimeSel, err := NewSel("span")
	check(err)
	tvSel, err := NewSel("td.tv__col")
	check(err)
	rowsResult := cascadia.QueryAll(doc, rowSel)
	events := [][]string{}
ParseRows:
	for _, row := range rowsResult {
		newEvent := make([]string, 5)
		// name and location
		raceResult := cascadia.Query(row, raceSel)
		if raceResult != nil {
			raceNameResult := cascadia.Query(raceResult, raceNameSel)
			raceTrackResult := cascadia.Query(raceResult, raceTrackSel)
			//TODO: add nil result handling
			newEvent[0] = raceNameResult.FirstChild.Data
			newEvent[3] = raceTrackResult.FirstChild.Data
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
				newEvent[1] = formattedDate
				newEvent[2] = formattedTime
			}
		}
		tvResult := cascadia.Query(row, tvSel)
		if tvResult != nil && tvResult.FirstChild != nil {
			newEvent[4] = tvResult.FirstChild.Data
		}
		for _, column := range newEvent {
			// if column is empty
			if column == "" {
				// do not add to events
				continue ParseRows
			}
		}
		events = append(events, newEvent)
	}
	dir := "csv"
	//GIVE ME DEM WRITE PERMS YA DIG?
	os.Mkdir(dir, 0777)
	fileName := path.Join(dir, "/events.csv")
	file, err := os.Create(fileName)
	check(err)
	defer file.Close()
	headers := []string{
		"Subject",
		"Start date",
		"Start time",
		"Location",
		"Description",
	}
	records := append([][]string{headers}, events...)
	writer := csv.NewWriter(file)
	for _, record := range records {
		if err := writer.Write(record); err != nil {
			log.Fatalln("error writing record to csv: ", err)
		}
	}
	writer.Flush()
	check(writer.Error())
	fmt.Printf("file %s created\n", fileName)
	fmt.Println("import it into your google calendar")
	fmt.Println("https://go.dev/doc/tutorial/add-a-test")
}
