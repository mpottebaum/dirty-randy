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

var NewSel = cascadia.Parse
var Query = cascadia.Query
var QueryAll = cascadia.QueryAll

func UserInput() (league string) {
	args := os.Args
	if len(args) <= 1 {
		fmt.Println("Please enter a league name")
		return
	}
	//league = args[0]
	//no need to lie
	league = "f1"
	return
}

func GetScheduleHTML(league string) *html.Node {
	url := "https://www.espn.com/" + league + "/schedule"
	resp, err := http.Get(url)
	check(err)
	defer resp.Body.Close()
	fmt.Println("Resp status: ", resp.Status)
	content, err := io.ReadAll(resp.Body)
	check(err)
	doc, err := html.Parse(strings.NewReader(string(content)))
	check(err)
	return doc
}

func BuildSels() (rowSel, raceSel, raceNameSel, raceTrackSel, lightsOutSel, dateTimeSel, tvSel cascadia.Sel) {
	var err error
	rowSel, err = NewSel("tr.Table__TR")
	check(err)
	raceSel, err = NewSel("td.race__col")
	check(err)
	raceNameSel, err = NewSel("a")
	check(err)
	raceTrackSel, err = NewSel("div")
	check(err)
	lightsOutSel, err = NewSel("td.winnerLightsOut__col")
	check(err)
	dateTimeSel, err = NewSel("span")
	check(err)
	tvSel, err = NewSel("td.tv__col")
	check(err)
	return
}

func GetRows(doc *html.Node, sel cascadia.Sel) []*html.Node {
	return QueryAll(doc, sel)
}

func GetNameAndLocation(race *html.Node, nameSel cascadia.Sel, trackSel cascadia.Sel) (name string, location string) {
	// name and location
	raceNameResult := Query(race, nameSel)
	raceTrackResult := Query(race, trackSel)
	if raceNameResult.FirstChild != nil {
		name = raceNameResult.FirstChild.Data
	}
	if raceTrackResult.FirstChild != nil {
		location = raceTrackResult.FirstChild.Data
	}
	return
}

func GetDateAndTime(lightsOut *html.Node, dateTimeSel cascadia.Sel) (formattedDate string, formattedTime string) {
	dateTimeResult := Query(lightsOut, dateTimeSel)
	//if race in past, there's no span (just winner's name)
	if dateTimeResult != nil && dateTimeResult.FirstChild != nil {
		rawDateTimeStr := dateTimeResult.FirstChild.Data
		fullDateAndTime := strings.Split(rawDateTimeStr, " - ")
		fullDate := fullDateAndTime[0]
		monthAndDay := strings.Split(fullDate, " ")
		month := monthAndDay[0]
		formattedMonth := MonthMap[month]
		day := monthAndDay[1]
		formattedDay := day
		if len(day) == 1 {
			formattedDay = "0" + formattedDay
		}
		eventTime := fullDateAndTime[1]
		hourAndPeriod := strings.Split(eventTime, " ")
		hour := hourAndPeriod[0]
		period := hourAndPeriod[1]
		formattedPeriod := strings.ToUpper(period)
		year := strconv.Itoa(time.Now().Year())
		formattedDate = formattedMonth + "/" + formattedDay + "/" + year
		formattedTime = hour + " " + formattedPeriod
	}
	return
}

func main() {
	//STEP ONE: GET LEAGUE NAME
	//TODO: support other leagues
	yourLeague := UserInput()
	//STEP TWO: GET SCHEDULE HTML
	doc := GetScheduleHTML(yourLeague)
	//STEP THREE: PARSE HTML
	rowSel,
		raceSel,
		raceNameSel,
		raceTrackSel,
		lightsOutSel,
		dateTimeSel,
		tvSel := BuildSels()
	rowsResult := GetRows(doc, rowSel)
	events := [][]string{}
ParseRows:
	for _, row := range rowsResult {
		newEvent := make([]string, 5)
		// name and location
		raceResult := Query(row, raceSel)
		if raceResult != nil {
			name, location := GetNameAndLocation(raceResult, raceNameSel, raceTrackSel)
			newEvent[0] = name
			newEvent[3] = location
		}
		// date and time
		lightsOutResult := Query(row, lightsOutSel)
		if lightsOutResult != nil {
			date, eventTime := GetDateAndTime(lightsOutResult, dateTimeSel)
				newEvent[1] = date
				newEvent[2] = formattedTime
			}
		}
		tvResult := Query(row, tvSel)
		if tvResult != nil && tvResult.FirstChild != nil {
			newEvent[4] = tvResult.FirstChild.Data
		}
		for _, column := range newEvent {
			// if any column is empty
			if column == "" {
				// do not add to events
				continue ParseRows
			}
		}
		events = append(events, newEvent)
	}
	//STEP FOUR: WRITE TO CSV 
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
	fmt.Println("import it into your google calendar, son")
}
