package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"google.golang.org/api/calendar/v3"
)

// gather stats from the list of events
type Employee struct {
	Email      string `json:"email"`
	Department string `json:"department"`
}

type Team struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type TeamCategory struct {
	Id    int     `json:"team_category_id"`
	Name  string  `json:"team_category_name"`
	Teams []*Team `json:"teams"`
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
	teamsFile, err := os.Open("teams.json")
	if err != nil {
		return
	}
	jsonParser := json.NewDecoder(teamsFile)
	teams := []TeamCategory{}
	err = jsonParser.Decode(&teams)
	if err != nil {
		return
	}
	var departmentNames = make(map[string]string)
	for _, category := range teams {
		fmt.Printf("category %s (%d)\n", category.Name, category.Id)
		for _, team := range category.Teams {
			fmt.Printf("team %s (%d)\n", team.Name, team.Id)
			_, ok := departmentNames[team.Id]
			if ok {
				fmt.Printf("DUPLICATE TEAM IDS")
				return
			}
			departmentNames[team.Id] = team.Name
		}
	}

	employeesFile, err := os.Open("employees.json")
	if err != nil {
		return
	}
	employees := []Employee{}
	jsonParser = json.NewDecoder(employeesFile)
	err = jsonParser.Decode(&employees)
	if err != nil {
		return
	}

	var departments = make(map[string]string)
	for _, employee := range employees {
		department := departmentNames[employee.Department]
		departments[employee.Email] = department
	}

	var colors = make(map[string]int)
	var departmentCount = make(map[string]int)
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
					departmentCount[department] += 1
				} else {
					if strings.Contains(strings.ToLower(attendee.Email), "sendgrid") {
						departmentCount["Sendgrid Other"] += 1
					} else {
						departmentCount["External"] += 1
					}

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
	for _, count := range departmentCount {
		total += count
	}
	for dpt, count := range departmentCount {
		if dpt == "" {
			dpt = "Unknown"
		}
		fmt.Printf("'%s' %f%%\n", dpt, 100.0*(float32(count)/float32(total)))
	}

	// show business
	fmt.Printf("WASTED\n")
	for date, seconds := range times {
		fmt.Printf("on date %d, %d seconds in meetings\n", date, seconds)
	}

	return err
}
