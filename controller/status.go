package controller

import (
	"time"
)

// Status holds a time-stamped status of various aspects of the device
type Status struct {
	Time      time.Time
	Locations Locations
	Action    Action
	Distances Distances // Only populated if Action is "Dispense"
}

// Status returns a snapshot of the status of the device.
func GetStatus() (Status, error) {
	return Status{}, nil
}

