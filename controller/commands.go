package controller

// controller provides functions and data structures for interacting with the
// device
// It is the primary interface used by Application code to interact with the Dispenser.

// The following operations are supported
// - Dispense paint from one or more syringes
// - Move one or more syringes to the top or bottom End positions
//    - (ie Home In or Home Out)
// - Get the current location of any Syringe
//    - A Location is one of Top, Middle or Bottom
//    - It is not possible to find out exactly where a syringe is when it is
//      is in the middle.
// - Get the status/state of the dispenser - ie what action is it performing,
//   what's the state of that action
import (
	"time"
)


// A Controller is used to interact with the dispenser.
// (FIXME Make it so that only one controller can be open at the same time)
// FIXME add actual contents to the Controller struct
type Controller struct {
}

// Open will open control 
func Open() (Controller, error) {
	// Try to open the gpio exclusively
	// Return error on fail
	c := Controller{}
	return c, nil
}

// Location is in Enum that indicates one of three locations that a syringe might exist in
type Location uint64
const (
	Top Location = iota
	Middle
	Bottom
)

// Colour refers to a Syringe by the colour it contains
type Colour uint64
const (
	C Colour = iota
	M
	Y
	K
	W
)

// Distance represents an amount of dispensing to be done
type Distance uint64

// Action represents an Action that the device is performing
type Action int
const (
	Dispense Action = iota
	HomeOut
	HomeIn
)

// Locations and Distances types hold the corresponding values
// for the full set of syringes available.
type Locations [5]Location
type Distances [5]Distance

// Status holds a time-stamped status of various aspects of the device
type Status struct {
	Time      time.Time
	Locations Locations
	Action    Action
	Distances Distances // Only populated if Action is "Dispense"
}

// Status returns a snapshot of the status of the device.
func (c *Controller) Status() Status {
	return Status{}
}

// Dispense will cause paint to be dispensed from syringes.
// The amount of paint dispensed is determined by the values in P
func (c *Controller) Dispense(d Distances) error {
	return nil
}

// HomeOut will dispense all paint from selected syringes
// Syringes are selected by having a non-zero value in Distances
func (c *Controller) HomeOut(d Distances) error {

	return nil
}

// HomeIn will draw up all syringes.
func (c *Controller) HomeIn(d Distances) error {

	return nil
}
