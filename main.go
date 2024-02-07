package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

// time-tracker start -n Task#1 -d <description>
// time-tracker end -n Task#1
// time-tracker summary   (default this week)
// time-tracker summary -st 01/01/2024 -et 01/07/2024
// time-tracker summary -w 1 -y 2024
// time-tracker update Task#1 -d description -st startTime -eT endTime -n name
// time-tracker add Task#1 -d description -st startTime -eT endTime -n name

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
	weekStr := flag.String("w", "", "Summary week filter")
	yearStr := flag.String("y", "", "Summary year filter")

	flag.Parse()

	app := &application{
		entries: make(map[string]*entry),
	}

	e1St, _ := time.ParseInLocation("2006-01-02", "2024-02-05", time.Local)
	e1Et, _ := time.ParseInLocation("2006-01-02", "2024-02-05", time.Local)

	e1 := &entry{
		name:        "Task 1",
		description: "Description 1",
		startTime:   e1St,
		endTime:     e1Et,
	}

	app.entries[e1.name] = e1

	e2 := &entry{
		name:        "Task 2",
		description: "Description 2",
		startTime:   time.Now().UTC(),
		endTime:     time.Now().UTC(),
	}

	app.entries[e2.name] = e2

	e3 := &entry{
		name:        "Task 3",
		description: "Description 3",
		startTime:   time.Now().Add(24 * time.Hour),
		endTime:     time.Now().Add(25 * time.Hour),
	}

	app.entries[e3.name] = e3

	e4 := &entry{
		name:        "Task 4",
		description: "Description 3",
		startTime:   time.Now().Add(-7 * 24 * time.Hour),
		endTime:     time.Now().Add(-7 * 24 * time.Hour),
	}

	app.entries[e4.name] = e4

	switch command {
	case "start":
		app.handleStart(*name, *description)
	case "end":
		app.handleEnd(*name)
	case "summary":
		app.handleSummary(*startTimeStr, *endTimeStr, *weekStr, *yearStr)
	default:

	}

	// Validate command
	// Handle Command
}

func (a *application) handleStart(name string, description string) {

	// Should be unique

	startTime := time.Now()

	e := &entry{
		name:        name,
		description: description,
		startTime:   startTime,
	}

	a.entries[name] = e
}

func (a *application) handleEnd(name string) error {

	e := &entry{
		name:    name,
		endTime: time.Now(),
	}

	err := a.updateEntry(e)
	if err != nil {
		return err
	}

	return nil

}

func (a *application) handleSummary(startTimeStr string, endTimeStr string, weekStr string, year string) error {
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
			"name: %v; startTime: %v; endTime: %v\n",
			e.name,
			e.startTime.Format("2006-01-02"),
			e.endTime.Format("2006-01-02"),
		)
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
	name        string
	description string
	startTime   time.Time
	endTime     time.Time
}

type application struct {
	entries map[string]*entry
}

func (a *application) saveEntry(e *entry) {
	a.entries[e.name] = e
}

func (a *application) deleteEntry(e *entry) {
	delete(a.entries, e.name)
}

func (a *application) updateEntry(e *entry) error {
	en, ok := a.entries[e.name]
	if !ok {
		return fmt.Errorf("error: entry not found")
	}

	if e.description != "" {
		en.description = e.description
	}

	if !e.startTime.IsZero() {
		en.startTime = e.startTime
	}

	if !e.endTime.IsZero() {
		en.endTime = e.endTime
	}

	return nil
}

func (a *application) getByDateRange(paramStartDate time.Time, paramEndDate time.Time) []*entry {

	entries := []*entry{}

	for _, e := range a.entries {
		if (paramStartDate.Before(e.startTime) || paramStartDate.Equal(e.startTime)) &&
			(paramEndDate.After(e.endTime) || paramEndDate.Equal(e.endTime)) {

			entries = append(entries, e)

		}
	}

	return entries
}

func (a *application) getActualWeek() []*entry {
	offset := -int(time.Now().Weekday())

	weekStart := time.Now().AddDate(0, 0, offset).Truncate(24 * time.Hour)
	weekEnd := weekStart.AddDate(0, 0, 4).Add(24*time.Hour - time.Nanosecond)

	return a.getByDateRange(weekStart, weekEnd)
}

// Save to File
// Get last saved
// name, startTime, endTime, description

// Connect with API?
// Get Bitbucket commit data to automatically fill tasks?
