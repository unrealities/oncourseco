package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

// gather stats from the list of events
type Employee struct {
	Email      string `json:"email"`
	Department string `json:"department"`
}

func dumpStats(events *calendar.Events) (err error) {
	err = nil
	//parkletKey := os.Getenv("PARKLET_KEY")
	//req, err := http.NewRequest("GET", "https://app.parklet.co/api/v1/employees?page=1", nil)
	//req.SetBasicAuth("nat.thompson@sendgrid.com", "HH9cZpv1Ehx3WbuHFSAB")
	//req.Header.Add("Content-Type", "application/json")
	//client := &http.Client{}
	//resp, err := client.Do(req)
	//if err != nil {
	//	return err
	//}
	//fmt.Printf("%s\n", resp)
	employeesFile, err := os.Open("employees.json")
	if err != nil {
		return
	}
	employees := []Employee{}
	jsonParser := json.NewDecoder(employeesFile)
	err = jsonParser.Decode(&employees)
	if err != nil {
		return
	}
	var departments = make(map[string]string)
	for _, employee := range employees {
		departments[employee.Email] = employee.Department
	}

	var colors = make(map[string]int)
	var department_count = make(map[string]int)
	var times = make(map[int64]int)
	fmt.Println("Upcoming events:")
	if len(events.Items) > 0 {
		for _, i := range events.Items {
			// If the DateTime is an empty string the Event is an all-day Event.
			// So only Date is available.
			var start time.Time
			var end time.Time
			// pull out start and end times
			if i.Start.DateTime != "" {
				start, _ = time.Parse(time.RFC3339, i.Start.DateTime)
				end, _ = time.Parse(time.RFC3339, i.End.DateTime)
			} else {
				start, _ = time.Parse("2000-01-22", i.Start.Date)
			}
			// for each calendar stay calculate the time spent in meetings
			day := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0,
				0, start.Location())
			timestamp := day.Unix()
			maxTime := 8 * 3600
			if end.IsZero() {
				times[timestamp] = maxTime
			} else {
				duration := int(end.Sub(start).Seconds())
				times[timestamp] += duration
				if times[timestamp] > maxTime {
					times[timestamp] = maxTime
				}
			}
			colors[i.ColorId] += 1
			fmt.Printf("%s (%s - %s)\n", i.Summary, start, end)

			for _, attendee := range i.Attendees {
				department, ok := departments[attendee.Email]
				if ok {
					department_count[department] += 1
				}
			}
		}
	} else {
		fmt.Printf("No upcoming events found.\n")
	}

	var total int
	// count colors
	fmt.Printf("COLORS\n")
	for _, count := range colors {
		total += count
	}
	for color, count := range colors {
		fmt.Printf("'%s' %f%%\n", color, 100.0*(float32(count)/float32(total)))
	}

	// count departments
	total = 0
	fmt.Printf("DEPARTMENTS\n")
	for _, count := range department_count {
		total += count
	}
	for dpt, count := range department_count {
		fmt.Printf("'%s' %f%%\n", dpt, 100.0*(float32(count)/float32(total)))
	}

	// show business
	fmt.Printf("WASTED\n")
	for date, seconds := range times {
		fmt.Printf("on date %d, %d seconds in meetings\n", date, seconds)
	}

	return err
}
