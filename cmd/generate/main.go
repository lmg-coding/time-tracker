package main

import (
	"encoding/gob"
	"fmt"
	"os"
	"time"
)

// Task struct to hold task details
type Task struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
}

func main() {
	const numOfTasks = 10_000_000
	tasks := make([]Task, numOfTasks) // Create a slice to hold 10,000 tasks

	// Define a start time for the first task
	startTime := time.Date(2024, 2, 5, 0, 0, 0, 0, time.UTC)

	for i := 0; i < numOfTasks; i++ {
		// Generate task details
		task := Task{
			Name:        fmt.Sprintf("Task %d", i+1),
			Description: fmt.Sprintf("Description for Task %d", i+1),
			StartTime:   startTime.Add(time.Hour * time.Duration(i)),                // Each task starts an hour after the previous one
			EndTime:     startTime.Add(time.Hour*time.Duration(i) + time.Minute*30), // Each task ends 30 minutes after it starts
		}

		tasks[i] = task
	}

	// Marshal the tasks slice into JSON
	f, err := os.Create("entryData.gob")
	if err != nil {
		fmt.Println(err)
		return
	}
	encoder := gob.NewEncoder(f)
	encoder.Encode(tasks)
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return
	}

	fmt.Println("Tasks saved to entryData.gob")
}
