package api

import (
	"fmt"
	"time"
)

type entry struct {
	Name         string      `json:"name"`
	Descriptions []string    `json:"descriptions"`
	Status       status      `json:"status"`
	Times        []entryTime `json:"times"`
}

type status string

const (
	Active   status = "Active"
	Finished status = "Finished"
)

type entryTime struct {
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}

func (e *entry) getElapsedTime() time.Duration {
	var elapsedTime time.Duration

	for _, et := range e.Times {
		elapsedTime += et.EndTime.Sub(et.StartTime)
	}

	return elapsedTime
}

func (e *entry) formatDescriptions() string {
	d := ""
	for _, desc := range e.Descriptions {
		d = fmt.Sprintf("%v\n%v", d, desc)
	}

	return d
}
