package api

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"time"
)

type Application struct {
	Entries         map[string]*entry `json:"entries"`
	ActiveEntry     string            `json:"activeEntry"`
	ActiveStartTime time.Time         `json:"activeStartTime"`
}

func ReadFromFile() (Application, error) {
	app := Application{
		Entries: make(map[string]*entry),
	}

	file, err := os.Open("entryData.gob")
	if err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create("entryData.gob")
			if err != nil {
				return app, err
			}
		} else {
			return app, err
		}
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return app, err
	}

	if fileInfo.Size() == 0 {
		return app, nil
	}

	decoder := gob.NewDecoder(file)

	if err := decoder.Decode(&app); err != nil {
		return app, err
	}

	return app, nil
}

func SaveToFile(app Application) error {
	file, err := os.Create("entryData.gob")
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(app)
	if err != nil {
		return err
	}

	return nil
}

func (a *Application) HandleStart(name string, description string) error {

	if _, ok := a.Entries[name]; ok {
		return fmt.Errorf("error: name should be unique")
	}

	if a.ActiveEntry != "" {
		return fmt.Errorf("error: ther is an active entry")
	}

	// Should be unique

	startTime := time.Now()

	e := &entry{
		Name:         name,
		Descriptions: []string{description},
	}

	a.Entries[name] = e

	a.ActiveEntry = name
	a.ActiveStartTime = startTime

	return nil
}

func (a *Application) HandleGetActive() {
	e, ok := a.Entries[a.ActiveEntry]
	if !ok {
		fmt.Println("there is not an active entry")
	}

	fmt.Printf(
		"name: %v; description: %v; totalTime: %v\n",
		e.Name,
		e.formatDescriptions(),
		e.getElapsedTime(),
	)
}

func (a *Application) HandlePause() error {

	t := entryTime{
		StartTime: a.ActiveStartTime,
		EndTime:   time.Now(),
	}

	e := &entry{
		Name:  a.ActiveEntry,
		Times: []entryTime{t},
	}

	err := a.appendTime(e)
	if err != nil {
		return err
	}

	a.ActiveStartTime = time.Time{}

	return nil
}

func (a *Application) HandleResume() error {

	if a.ActiveEntry == "" {
		return fmt.Errorf("error: there is not active entry")
	}

	a.ActiveStartTime = time.Now()

	return nil
}

func (a *Application) HandleEnd() error {
	zeroTime := time.Time{}

	if a.ActiveEntry == "" {
		return fmt.Errorf("error: there is not active entry")
	}

	if a.ActiveStartTime == zeroTime {
		a.ActiveEntry = ""
		return nil
	}

	t := entryTime{
		StartTime: a.ActiveStartTime,
		EndTime:   time.Now(),
	}

	e := &entry{
		Name:  a.ActiveEntry,
		Times: []entryTime{t},
	}

	err := a.appendTime(e)
	if err != nil {
		return err
	}

	a.ActiveEntry = ""
	a.ActiveStartTime = time.Time{}

	return nil

}

func (a *Application) HandleAddDescription(desc string) error {
	if desc == "" {
		return fmt.Errorf("error: please provide a description")
	}

	e := &entry{
		Name:         a.ActiveEntry,
		Descriptions: []string{desc},
	}

	err := a.appendDescription(e)
	if err != nil {
		return err
	}

	return nil
}

func (a *Application) HandleSummary(startTimeStr string, endTimeStr string, year string) error {
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
			"name: %v; description: %v; totalTime: %v\n",
			e.Name,
			e.formatDescriptions(),
			e.getElapsedTime(),
		)
	}

	return nil
}

func (a *Application) HandleAdd(name string, desc string, startTimeStr string, endTimeStr string) error {
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

	startTime, err = time.ParseInLocation("2006-01-02 03:04 pm", startTimeStr, time.Local)
	if err != nil {
		return err
	}

	endTime, err = time.ParseInLocation("2006-01-02 03:04 pm", endTimeStr, time.Local)
	if err != nil {
		return err
	}

	t := entryTime{
		StartTime: startTime,
		EndTime:   endTime,
	}

	e := &entry{
		Name:         name,
		Descriptions: []string{desc},
		Times:        []entryTime{t},
	}

	a.saveEntry(e)

	return nil
}

func (a *Application) HandleUpdate(name string, desc string, startTimeStr string, endTimeStr string, timeIndex int) error {

	var startTime time.Time
	var endTime time.Time
	var err error

	if name == "" {
		return fmt.Errorf("error: name must be provided")
	}

	if startTimeStr != "" {
		startTime, err = time.ParseInLocation("2006-01-02 03:04 pm", startTimeStr, time.Local)
		if err != nil {
			return err
		}
	}

	if endTimeStr != "" {
		endTime, err = time.ParseInLocation("2006-01-02 03:04 pm", startTimeStr, time.Local)
		if err != nil {
			return err
		}
	}

	e := &entry{
		Name:         name,
		Descriptions: []string{desc},
		Times: []entryTime{
			{
				StartTime: startTime,
				EndTime:   endTime,
			},
		},
	}

	err = a.updateEntry(e, timeIndex)
	if err != nil {
		return err
	}

	return nil
}

func (a *Application) HandleDelete(name string) error {
	if name == "" {
		return fmt.Errorf("error: name must be provided")
	}

	a.deleteEntry(name)

	return nil
}
