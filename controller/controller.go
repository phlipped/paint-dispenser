package controller

// controller provides functions and data structures for interacting with the
// device
// It is the primary interface used by Application code to interact with the Dispenser.

// The following operations are supported
// - Dispense paint from one or more syringes
// - Move one or more syringes to the top or bottom End positions
//    - (ie Home In or Home Out)
// - Get the status/state of the dispenser
//    - ie what action is it performing,
//    - what's the state of that action
//    - where is each syringe right now.

// Location is an Enum that indicates one of three locations that a syringe might exist in
// Note that it is not possible to know exactly where a syringe is when it is in the Middle


type Location uint
const (
	Top Location = iota
	Middle
	Bottom
)

type Direction uint
const (
	Up Direction = 1 // The plunger is drawn up, ie takes in paint
	Down Direction = 0 // The plunger is pushed down, ie dispenses paint
)

// Colour refers to a Syringe by the colour it contains
type Colour uint
const (
	C Colour = iota
	M
	Y
	K
	W
)

// Action represents an Action that the device is performing
type Action uint
const (
	Idle Action = iota
	Dispensing
	HomingOut
	HomingIn
)

// Distance represents an amount of dispensing to be done
// The units are not yet defined, but most likely it will
// represent the number of steps that the stepper motor needs to make
type Distance uint64

// Locations and Distances types hold the corresponding values
// for all Syringes.
type Locations [5]Location
type Distances [5]Distance


// Dispense will cause paint to be dispensed from syringes.
// The amount of paint dispensed is determined by the values in d
func Dispense(dists Distances) error {
	return nil
}

// HomeOut will dispense all paint from selected syringes
// Syringes are selected by having a non-zero value in d
func HomeOut(dists Distances) error {

	return nil
}

// HomeIn will draw up the selected syringes.
// Syringes are selected by having a non-zero value in d
func HomeIn(dists Distances) error {

	return nil
}
