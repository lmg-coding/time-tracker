package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/lmg-coding/time-tracker/api"
)

// time-tracker start -n Task#1 -d <description>
// time-tracker end -n Task#1
// time-tracker summary   (default this week)
// time-tracker summary -st 01/01/2024 -et 01/07/2024
// time-tracker update Task#1 -d description -st startTime -eT endTime
// time-tracker add Task#1 -d description -st startTime -eT endTime

func main() {

	command := os.Args[1]
	os.Args = os.Args[1:]

	name := flag.String("n", "", "Name of the entry")
	description := flag.String("d", "", "Description of the entry")
	startTimeStr := flag.String("st", "", "Start time of the entry")
	endTimeStr := flag.String("et", "", "End time of the entry")
	yearStr := flag.String("y", "", "Summary year filter")
	index := flag.Int("i", 0, "Time index for update")

	flag.Parse()

	app, err := api.ReadFromFile()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	switch command {
	case "start":
		err = app.HandleStart(*name, *description)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = api.SaveToFile(app)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case "end":
		err := app.HandleEnd()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = api.SaveToFile(app)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case "summary":
		err := app.HandleSummary(*startTimeStr, *endTimeStr, *yearStr)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	case "update":
		err := app.HandleUpdate(*name, *description, *startTimeStr, *endTimeStr, *index)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = api.SaveToFile(app)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case "pause":
		err := app.HandlePause()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = api.SaveToFile(app)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case "resume":
		err = app.HandleResume()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = api.SaveToFile(app)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case "delete":
		err = app.HandleDelete(*name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = api.SaveToFile(app)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case "add":
		err := app.HandleAdd(*name, *description, *startTimeStr, *endTimeStr)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = api.SaveToFile(app)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case "add-description":
		err := app.HandleAddDescription(*description)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = api.SaveToFile(app)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case "active":
		app.HandleGetActive()
	default:

	}
}
