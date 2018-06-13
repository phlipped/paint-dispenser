/*
Thoughts ...

200 steps per revolution, with x16 microstepping, means 3200 steps per revolution.
Thread pitch is 1.25mm
Max speed needed is ... say ... 10mm per second ...

10mm/s =
(10mm/1.25mm)revs/s =
8revs/s =
(8 * 3200)steps/s =
25600 steps/s = 
40 micros per step = 

So ... 40 micros to do everything we need to do between steps.

HARDWARE vs SOFTWARE PWM

Hardware Pros ...
- No CPU cost
- Very precise timings

Hardware Cons
- Needs to be actively turned off - risk damage if software goes AWOL
- Perhaps complicated to change the frequency of pulses to allow slow ramp-up ramp down.


Software tests ....
I ran a tight loop in go and got ~1 microsecond granularity on times

Filling an array of 10000 ints seemed to take anywhere from 400 to 3400 micros ..
so that's .... 
0.03 micros to 0.24 micros

100000 took 11000 micros, pretty consistently ...


so each loop iteration took ... 0.1 micros

So it's probably a bit lumpy, but it might not be a big deal, and might be only in the order of 1-2 micros?
Which is fine for my application, because I've got 40 micros to play with ...

So ... software PWM it is ...
*/
package controller

import (
	"github.com/stianeikeland/go-rpio"
)

const (
	stepPulseWidthMicros = 1
	stepPulseMinPeriodMicros = 20
	stepPulseMaxPeriodMicros = 100
)

type StepCount uint32

// doSteps moves the steppers the number of steps specified in distances
// FIXME make it ease-in at the start
func doSteps(dists Distances, dir Direction) {
	// Set the direction
	dirPin := Pins["dir"]
	if dir == Up {
		// FIXME check which way round this should be ...
		dirPin.Write(rpio.High)
	} else {
		dirPin.Write(rpio.Low)
	}

	// Init the step pin
	stepPin := Pins["step"]
	stepPin.Write(rpio.Low)

	// Init array of Enable Pins
	enablePins := [5]rpio.Pin{}
	for i, pinName := range EnablePinNames {
		enablePins[i] = Pins[pinName]
		enablePins[i].Write(rpio.Low) // Low == Enabled
	}

	// Init Limit Pin Inputs
	limitPinBank := []string{}
	if dir == Up {
		limitPinBank = LimitUpperPinNames
	} else {
		limitPinBank = LimitLowerPinNames
	}
	limitPins := [5]rpio.Pin{}
	for i, pinName := range limitPinBank {
		limitPins[i] = Pins[pinName]
	}

	// This loop has to be tight, yo
	// while any of the pins are still enabled ...
	enableMask := 0x1f // 0b11111 - 1 bit for each of the steppers
	for enableMask != 0 {
		// Check distances - if zero, disable corresponding enable pin, and update the enable Mask
		for i := uint8(0); i < 5; i++ {
			// If the distance 
			if dists[i] == 0 || limitPins[i].Read() == rpio.High {
				enableMask ^= 1 << i
				enablePins[i].Write(rpio.High) // High == Disabled
			}
		}
		// Check if any of the limit switches are triggered - if so, disable the corresponding pin.
		// Push the step pin up ... wait a microsecond and then pull it down again.
		// Wait out the timeout until the next step is due ...
	}

}
