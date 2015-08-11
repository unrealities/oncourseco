package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
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

type Stats struct {
	Name        string             `json:"name"`
	Departments map[string]float64 `json:"departments"`
	Wasted      Wasted             `json:"wasted"`
}

type Wasted struct {
	TwentyEightDays float64 `json:"twentyeight"`
	SevenDays       float64 `json:"seven"`
}

func dumpStats(events *calendar.Events) (stats Stats, err error) {
	err = nil
	stats = Stats{}
	stats.Departments = make(map[string]float64)
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

	// parse JSON
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
		for _, team := range category.Teams {
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

	// Analyze Events
	var colors = make(map[string]int)
	var departmentCount = make(map[string]int)
	var times = make(map[int]int)
	workDaySeconds := 8 * 3600
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
			timestamp := int(day.Unix())
			if end.IsZero() {
				times[timestamp] = workDaySeconds
			} else {
				duration := int(end.Sub(start).Seconds())
				times[timestamp] += duration
				if times[timestamp] > workDaySeconds {
					times[timestamp] = workDaySeconds
				}
			}
			colors[i.ColorId] += 1

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
		return
	}

	// Generate Stats
	var total int
	// count colors
	fmt.Printf("COLORS\n")
	for _, count := range colors {
		total += count
	}
	for color, count := range colors {
		fmt.Printf("'%s' %f%%\n", color, (float32(count) / float32(total)))
	}

	// count departments
	total = 0
	for _, count := range departmentCount {
		total += count
	}
	for dpt, count := range departmentCount {
		if dpt == "" {
			dpt = "Unknown"
		}
		stats.Departments[dpt] = (float64(count) / float64(total))
	}

	// Time in meetings
	timestamps := make([]int, 0, len(times))
	for k := range times {
		timestamps = append(timestamps, k)
	}
	sort.Ints(timestamps)
	now := time.Now()
	nowDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0,
		0, now.Location())
	nowTimestamp := int(nowDay.Unix())
	oneWeekSeconds := 86400 * 7
	oneWeekWorkSeconds := workDaySeconds * 5
	fourWeekWorkSeconds := 4 * oneWeekWorkSeconds
	oneWeekCutoff := nowTimestamp - oneWeekSeconds
	fourWeekCutoff := nowTimestamp - (4 * oneWeekSeconds)
	var oneWeekWasted int
	var fourWeeksWasted int
	for i := len(timestamps) - 1; i >= 0; i-- {
		if timestamps[i] > nowTimestamp {
			continue
		}
		if timestamps[i] > oneWeekCutoff {
			oneWeekWasted += times[timestamps[i]]
		}
		if timestamps[i] > fourWeekCutoff {
			fourWeeksWasted += times[timestamps[i]]
		}
	}

	stats.Wasted.SevenDays = (float64(oneWeekWasted) / float64(oneWeekWorkSeconds))
	stats.Wasted.TwentyEightDays = (float64(fourWeeksWasted) / float64(fourWeekWorkSeconds))

	return
}
