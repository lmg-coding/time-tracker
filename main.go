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

	switch command {
	case "start":
		app.handleStart(*name, *description)
	case "end":
		app.handleEnd(*name)
	case "summary":
		app.handleSummary(*startTimeStr, *endTimeStr, *weekStr, *yearStr)
	default:

	}

	fmt.Println(app.entries)
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

func (a *application) handleSummary(startTimeStr string, endTimeStr string, weekStr string, year string) {
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

// Save to File
// Get last saved
// name, startTime, endTime, description

// Connect with API?
// Get Bitbucket commit data to automatically fill tasks?
