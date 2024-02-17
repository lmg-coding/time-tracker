package api

import (
	"fmt"
	"time"
)

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

func (a *Application) saveEntry(e *entry) {
	a.Entries[e.Name] = e
}

func (a *Application) deleteEntry(name string) {
	delete(a.Entries, name)
}

func (a *Application) appendTime(e *entry) error {
	en, ok := a.Entries[e.Name]
	if !ok {
		return fmt.Errorf("error: entry not found")
	}

	for _, t := range e.Times {
		en.Times = append(en.Times, t)
	}

	return nil
}

func (a *Application) appendDescription(e *entry) error {
	en, ok := a.Entries[e.Name]
	if !ok {
		return fmt.Errorf("error: entry not found")
	}

	for _, desc := range e.Descriptions {
		en.Descriptions = append(en.Descriptions, desc)
	}

	return nil
}

func (a *Application) updateEntry(e *entry, index int) error {

	en, ok := a.Entries[e.Name]
	if !ok {
		return fmt.Errorf("error: entry not found")
	}

	if index >= len(en.Times) {
		return fmt.Errorf("error: timeIndex not found")
	}

	if len(e.Descriptions) != 0 {
		en.Descriptions[index] = e.Descriptions[0]
	}

	if len(e.Times) > 0 && !e.Times[0].StartTime.IsZero() {
		en.Times[index].StartTime = e.Times[0].StartTime
	}

	if len(e.Times) > 0 && !e.Times[0].EndTime.IsZero() {
		en.Times[index].EndTime = e.Times[0].EndTime
	}

	return nil
}

func (a *Application) getByDateRange(paramStartDate time.Time, paramEndDate time.Time) []*entry {

	entries := []*entry{}

	for _, e := range a.Entries {
		if (paramStartDate.Before(e.Times[0].StartTime) || paramStartDate.Equal(e.Times[0].StartTime)) &&
			(paramEndDate.After(e.Times[len(e.Times)-1].EndTime) || paramEndDate.Equal(e.Times[len(e.Times)-1].EndTime)) {

			entries = append(entries, e)

		}
	}

	return entries
}

func (a *Application) getActualWeek() []*entry {
	offset := -int(time.Now().Weekday())

	weekStart := time.Now().AddDate(0, 0, offset).Truncate(24 * time.Hour)
	weekEnd := weekStart.AddDate(0, 0, 7).Add(24*time.Hour - time.Nanosecond)

	return a.getByDateRange(weekStart, weekEnd)
}
