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
	"io"
	"github.com/phlipped/gopigpio"
)

const (
	stepPulseWidthMicros = 1
	stepPulseMinPeriodMicros = 20
	stepPulseMaxPeriodMicros = 100
)

type StepCount uint32

// doSteps moves the steppers the number of steps specified in distances
// FIXME make it ease-in at the start
func doSteps(p io.ReadWriter, dists Distances, dir Direction) {
	// Set the direction
	var dirVal gopigpio.PinVal
	var limitPins []gopigpio.Pin
	if dir == Up {
		// FIXME check which way round this should be ...
		dirVal = gopigpio.GPIO_LOW
		limitPins = LimitUpperPins
	} else {
		dirVal = gopigpio.GPIO_HIGH
		limitPins = LimitLowerPins
	}
	if err := gopigpio.GpioWrite(p, DIR_PIN, dirVal); err != nil {
		panic(err) // FIXME do better
	}
	// Setting the pull direction should be done in some kind of init phase, IMHO, and shouldn't need to be redone each time we do steps
	for _, pin := range limitPins {
		if err := gopigpio.GpioSetPullUpDown(p, pin, gopigpio.GPIO_PULL_HIGH); err != nil { // FIXME check which way we want to pull them
			panic(err) // FIXME do better
		}
	}

	// Set step pin to definitely be turned off
	if err := gopigpio.GpioWrite(p, STEP_PIN, gopigpio.GPIO_LOW); err != nil {
		panic(err) // FIXME do better
	}

	// Init all Enable Pins to turn them on
	for _, pin := range EnablePins {
		if err := gopigpio.GpioWrite(p, pin, gopigpio.GPIO_HIGH); err != nil { // Low == Enabled
			panic(err) // FIXME do better
		}
	}

	// FIXME actually make the waveforms that we need to make for each pin

}
