package main

import (
	"time"
)

// DaysLeft returns the number of days between today and project’s target date, both ends included
func (p *Project) DaysLeft() int {
	return int(p.TargetDate.Add(24*time.Hour).Truncate(24*time.Hour).Sub(time.Now().Truncate(24*time.Hour)) / (24 * time.Hour))
}

// MeasuresTotals returns a map of the sum of values for each measures of a project
func (p *Project) MeasuresTotals() map[int]int {
	ret := make(map[int]int)
	for _, note := range p.Notes {
		for mID, mValue := range note.MeasuresValues {
			if _, ok := ret[mID]; !ok {
				ret[mID] = mValue
			} else {
				ret[mID] += mValue
			}
		}
	}
	for mID := range p.Measures {
		if _, ok := ret[mID]; !ok {
			ret[mID] = 0
		}
	}
	return ret
}
