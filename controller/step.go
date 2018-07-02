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
	"time"
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
func doSteps(p io.ReadWriter, dists Distances, dir Direction) error {
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
		return err // FIXME do better - custom error and logging
	}
	// Setting the pull direction should be done in some kind of init phase, IMHO, and shouldn't need to be redone each time we do steps
	for _, pin := range limitPins {
		if err := gopigpio.GpioSetPullUpDown(p, pin, gopigpio.GPIO_PULL_HIGH); err != nil { // FIXME confirm which way we want to pull them
			return err // FIXME do better - custom error and logging
		}
	}

	// Set step pin to definitely be turned off
	if err := gopigpio.GpioWrite(p, STEP_PIN, gopigpio.GPIO_LOW); err != nil {
		return err // FIXME do better - custom error and logging
	}

	// Init all Enable Pins to turn them on
	for _, pin := range EnablePins {
		if err := gopigpio.GpioWrite(p, pin, gopigpio.GPIO_HIGH); err != nil { // Low == Enabled
			return err // FIXME do better - custom error and logging
		}
	}

	// FIXME actually make the waveforms that we need to make for each pin
	// do we loop over <dists>, and create pigpio.WaveLoop and gopigpio.WavePulses, and change the delay value in each pulse over the course of the loop
	// Possibly need to make use of waveform loops structures to make the pulse structures smaller

	// Ok, so here's the plan ... we create Pulse objects ...
	// Some pulses are used to generate a STEP.
	// Other pulses are used to TURN OFF the enable pins of a particular stepper motor.
	// Hmmm, we can pre-make the pulses that turn off stepper motors, and just interject them into the waveform/loop structure at suitable times.
	// So basically, we need to define each of the STEP pulses we are using for the ramp-up, and also the count for each of them.
	// Then we need to work out where to inject the pulses that disable each stepper at the right time.
	// And then we need to build small loops of step structures with the right number of repititions, and then a tiny litte pulse that disables a particular stepper
	// And then have another small step-loop section with the right delay in it, until either the delay changes, or we need to disable another stepper.

	// What we want to do is be able to move steppers according to a rate function
	// The rate function defines rate vs time
	// But we can't "smoothly" change the rate with each step, so we are going to approximate the rate function with a series of straight lines
	// Now ... we could just draw some straight lines to connect some points on the rate function, but we risk accumulating errors and ending up with an overall
	// rate that doesn't match the expected rate.
	// So instead ... we do this ...
	// Then, we evaluate the definite integral of the rate function at various intervals.
	// The definite integral is, by definition, the distance. Therefore we can calculate an equivalent rate for that interval - distance/time
	// Now we just have a series of flat rates that we should follow one after the other.
	// Because rates will be implemented as series of discrete pulses, the full time interval may not be filled up completely. Therefore the next time interval
	// will need to by determined based on the end of the previous time interval dynamically.


	// While there are still some steppers to move ...
	for ! dists.AllZero() {
		// 
	}

	return nil
}

// PulseDelayInterval describes an interval, and a PulseDelay to use in that interval
// A series of PulseDelayIntervals can be used to describe a changing rate of pulses
type PulseDelayInterval struct {
	start time.Duration // The start time of this Interval, assuming the first interval in a sequence of
			    //   intervals starts at 0. Each subsequent interval should start exactly at
			    //   <prevInterval>.start + <prevInterval>.width
	width time.Duration // The duration of the interval
	delay time.Duration // How long the delay between pulse starts should be. This does NOT account for
			    //   the width of the pulse itself
}


// calcRampRates calculates a series of constant rates that approximate a continuously variable rate function
// In addition to the rate function, startX and endX values must be provided that
//   describe which part of the rate function to use.
// To calulate a suitable rate, the definite integral of the rate function is calculated over each interval,
//   which is used to calculate
// an equivalent constant rate over that interval.
// The actual end time of each interval may be adjusted so that there is an integer number of pulses in that interval.
// The adjustment of the end time of an interval will affect the start time of
//   the next interval. The normal end time of the next interval is still used
//   as the target end time for that interval, but it may also be subject to
//   adjustment when the rate for that interval is calculated.
func calcPulseDelayIntervals(
	rateFunc func(float64) float64, // The rate function
	startX float64, // The input value to use for the start of the rate function
	endX float64, // The input value to use for the end of the rate function
	startRate float64, // The starting rate (typically, this will be zero)
	endRate float64, // The rate to reach by the end of the ramp up time
	rampTime time.Duration, // Total duration of ramp up
	intervals int, // Number of intervals to generate
) []PulseDelayInterval {
	var pdis []PulseDelayInterval


	return pdis
}
