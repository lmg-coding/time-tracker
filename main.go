package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"
)

// time-tracker start -n Task#1 -d <description>
// time-tracker end -n Task#1
// time-tracker summary   (default this week)
// time-tracker summary -st 01/01/2024 -et 01/07/2024
// time-tracker update Task#1 -d description -st startTime -eT endTime
// time-tracker add Task#1 -d description -st startTime -eT endTime

// Better to use something like Cobra, but we are going to build just for the sake of learning
func main() {

	// Get command
	if len(os.Args) < 2 {
		fmt.Println("error: no command provided")
		os.Exit(1)
	}

	command := os.Args[1]
	os.Args = os.Args[1:]

	name := flag.String("n", "", "Name of the entry")
	description := flag.String("d", "", "Description of the entry")
	startTimeStr := flag.String("st", "", "Start time of the entry")
	endTimeStr := flag.String("et", "", "End time of the entry")
	yearStr := flag.String("y", "", "Summary year filter")

	flag.Parse()

	app := &application{}

	entriesSlice := []*entry{}

	bytes, err := os.ReadFile("entryData.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = json.Unmarshal(bytes, &entriesSlice)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	entries := sliceToMap(entriesSlice)

	app.entries = entries

	switch command {
	case "start":
		app.handleStart(*name, *description)
	case "end":
		err := app.handleEnd(*name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case "summary":
		err := app.handleSummary(*startTimeStr, *endTimeStr, *yearStr)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case "update":
		err := app.handleUpdate(*name, *description, *startTimeStr, *endTimeStr)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case "add":
		err := app.handleAdd(*name, *description, *startTimeStr, *endTimeStr)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	default:

	}

	err = saveToFile(app.entries)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Validate command
	// Handle Command
}

func saveToFile(entries map[string]*entry) error {
	entriesSlice := mapToSlice(entries)

	entriesJson, err := json.Marshal(entriesSlice)
	if err != nil {
		return err
	}

	err = os.WriteFile("entryData.json", entriesJson, 0644)
	if err != nil {
		return err
	}

	return nil
}

func mapToSlice(entriesMap map[string]*entry) []*entry {
	entriesSlice := []*entry{}

	for _, e := range entriesMap {
		entriesSlice = append(entriesSlice, e)
	}

	return entriesSlice
}

func sliceToMap(entriesSlice []*entry) map[string]*entry {
	entriesMap := make(map[string]*entry)

	for _, e := range entriesSlice {
		entriesMap[e.Name] = e
	}

	return entriesMap
}

func (a *application) handleStart(name string, description string) {

	// Should be unique

	startTime := time.Now()

	e := &entry{
		Name:        name,
		Description: description,
		StartTime:   startTime,
	}

	a.entries[name] = e
}

func (a *application) handleEnd(name string) error {

	e := &entry{
		Name:    name,
		EndTime: time.Now(),
	}

	err := a.updateEntry(e)
	if err != nil {
		return err
	}

	return nil

}

func (a *application) handleSummary(startTimeStr string, endTimeStr string, year string) error {
	// I need a per day description of the task I've done

	entries := []*entry{}

	if startTimeStr != "" || endTimeStr != "" {
		startDate, endDate, err := getDates(startTimeStr, endTimeStr)
		if err != nil {
			return err
		}

		entries = a.getByDateRange(startDate, endDate)
	} else {
		entries = a.getActualWeek()
	}

	for _, e := range entries {
		fmt.Printf(
			"name: %v; description: %v; startTime: %v; endTime: %v; totalTime: %v\n",
			e.Name,
			e.Description,
			e.StartTime.Format("2006-01-02"),
			e.EndTime.Format("2006-01-02"),
			e.EndTime.Sub(e.StartTime),
		)
	}

	return nil
}

// time-tracker add Task#1 -d description -st startTime -eT endTime
func (a *application) handleAdd(name string, desc string, startTimeStr string, endTimeStr string) error {
	var startTime time.Time
	var endTime time.Time
	var err error

	if name == "" {
		return fmt.Errorf("error: name must be provided")
	}

	if desc == "" {
		return fmt.Errorf("error: description must be provided")
	}

	if startTimeStr == "" {
		return fmt.Errorf("error: start time must be provided")
	}

	if endTimeStr == "" {
		return fmt.Errorf("error: end time must be provided")
	}

	startTime, err = time.ParseInLocation("2006-01-02", startTimeStr, time.Local)
	if err != nil {
		return err
	}

	endTime, err = time.ParseInLocation("2006-01-02", endTimeStr, time.Local)
	if err != nil {
		return err
	}

	e := &entry{
		Name:        name,
		Description: desc,
		StartTime:   startTime,
		EndTime:     endTime,
	}

	a.saveEntry(e)

	return nil
}

func (a *application) handleUpdate(name string, desc string, startTimeStr string, endTimeStr string) error {

	var startTime time.Time
	var endTime time.Time
	var err error

	if name == "" {
		return fmt.Errorf("error: name must be provided")
	}

	if startTimeStr != "" {
		startTime, err = time.ParseInLocation("2006-01-02", startTimeStr, time.Local)
		if err != nil {
			return err
		}
	}

	if endTimeStr != "" {
		endTime, err = time.ParseInLocation("2006-01-02", startTimeStr, time.Local)
		if err != nil {
			return err
		}
	}

	e := &entry{
		Name:        name,
		Description: desc,
		StartTime:   startTime,
		EndTime:     endTime,
	}

	err = a.updateEntry(e)
	if err != nil {
		return err
	}

	return nil
}

func getDates(startTimeStr string, endTimeStr string) (time.Time, time.Time, error) {

	startDate, err := time.ParseInLocation("2006-01-02", startTimeStr, time.Local)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	endDate, err := time.ParseInLocation("2006-01-02", endTimeStr, time.Local)
	if err != nil {
		return time.Time{}, time.Time{}, nil
	}

	endDate = endDate.Add(24*time.Hour - time.Nanosecond)

	return startDate, endDate, err
}

type entry struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
}

type application struct {
	entries map[string]*entry
}

func (a *application) saveEntry(e *entry) {
	a.entries[e.Name] = e
}

func (a *application) deleteEntry(e *entry) {
	delete(a.entries, e.Name)
}

func (a *application) updateEntry(e *entry) error {
	en, ok := a.entries[e.Name]
	if !ok {
		return fmt.Errorf("error: entry not found")
	}

	if e.Description != "" {
		en.Description = e.Description
	}

	if !e.StartTime.IsZero() {
		en.StartTime = e.StartTime
	}

	if !e.EndTime.IsZero() {
		en.EndTime = e.EndTime
	}

	return nil
}

func (a *application) getByDateRange(paramStartDate time.Time, paramEndDate time.Time) []*entry {

	entries := []*entry{}

	for _, e := range a.entries {
		if (paramStartDate.Before(e.StartTime) || paramStartDate.Equal(e.StartTime)) &&
			(paramEndDate.After(e.EndTime) || paramEndDate.Equal(e.EndTime)) {

			entries = append(entries, e)

		}
	}

	return entries
}

func (a *application) getActualWeek() []*entry {
	offset := -int(time.Now().Weekday())

	weekStart := time.Now().AddDate(0, 0, offset).Truncate(24 * time.Hour)
	weekEnd := weekStart.AddDate(0, 0, 7).Add(24*time.Hour - time.Nanosecond)

	return a.getByDateRange(weekStart, weekEnd)
}

// Save to File
// Get last saved
// name, startTime, endTime, description

// Connect with API?
// Get Bitbucket commit data to automatically fill tasks?
