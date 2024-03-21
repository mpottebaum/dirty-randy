package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
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
	reggieTbody, _ := regexp.Compile("<tbody.*</tbody>")
	reggieRow, _ := regexp.Compile("<tr.*</tr>")
	//	reggieDee, _ := regexp.Compile("<td.*</td>")
	reggieRace, _ := regexp.Compile("<a.*/f1/race.*</a>")
	reggieTrack, _ := regexp.Compile("<div.*</div>")
	reggieDateTime, _ := regexp.Compile("<td.*winnerLightsOut__col.*<span></span>")
	if err != nil {
		log.Fatal(err)
	}
	tableBody := reggieTbody.FindString(string(content))
	scheduleRows := reggieRow.FindAllString(tableBody, -1)

	// look for class names
	for i := 0; i < len(scheduleRows); i++ {
		// store data based on class name match
		scheduleRow := scheduleRows[i]
		race := reggieRace.FindString(scheduleRow)
		track := reggieTrack.FindString(scheduleRow)
		dateTime := reggieDateTime.FindString(scheduleRow)
		fmt.Println("******************")
		fmt.Println("race", race)
		fmt.Println("track", track)
		fmt.Println("dateTime", dateTime)
		fmt.Println("******************")
	}
	resp.Body.Close()
}
