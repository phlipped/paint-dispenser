package controller

import (
	"fmt"
	"io"
	"math"
	"time"

	"github.com/phlipped/gopigpio"
	"github.com/TheDemx27/calculus"
)

const (
	stepPulseWidth = time.Duration(2) * time.Microsecond
	stepPulseWait = time.Duration(2) * time.Microsecond
	stepPulseMinPeriod = time.Duration(20) * time.Microsecond
	stepPulseMaxPeriod = time.Duration(100) * time.Microsecond
	delayUntilDisable = time.Duration(2) * time.Microsecond
)

type StepCount uint32

///////////////////////////////
// Ramp Function Stuff
func RampShapeFunc(x float64) float64 {
	return x
}

var (
	rampFuncStart = 0.0
	rampFuncEnd = 1.0

	rampDuration = time.Duration(1) * time.Second
	maxRate = 50000.0
	rampIntervals = 10
	steadyPDI = PulsePeriodInterval{
		period: time.Duration(int(1000000000.0 / maxRate + 0.5)), // +0.5 to make it round, rather than truncate
		count: math.MaxUint64,
	}
)

// Now create the actual rampFunction based on those other parameters
var rampFunc = WrapFunction(RampShapeFunc, rampFuncStart, rampFuncEnd, 0.0, float64(rampDuration.Seconds()), 0.0, maxRate)

//
///////////////////////////////




// DoSteps moves the steppers the number of steps specified in <dists>
func DoSteps(pigpio io.ReadWriter, dists Distances, dir Direction) error {

	// FIXME Clear out all the current waves

	// Set Direction and LimitPin variables based on <dir>
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
	if err := gopigpio.GpioWrite(pigpio, DIR_PIN, dirVal); err != nil {
		return err // FIXME do better - custom error and logging
	}

	// Setting the GPIO pull on the limit pin inputs should be done in some kind of init phase, IMHO,
	// and shouldn't need to be redone each time we do steps
	// FIXME make this a function, if nothing else
	for _, pin := range limitPins {
		if err := gopigpio.GpioSetPullUpDown(pigpio, pin, gopigpio.GPIO_PULL_HIGH); err != nil { // FIXME confirm which way we want to pull them
			return err // FIXME do better - custom error and logging
		}
	}

	// Set step pin to definitely be turned off
	if err := gopigpio.GpioWrite(pigpio, STEP_PIN, gopigpio.GPIO_LOW); err != nil {
		return err // FIXME do better - custom error and logging
	}

	// Init all Enable Pins so the motors are then enabled
	for _, pin := range EnablePins {
		if err := gopigpio.GpioWrite(pigpio, pin, gopigpio.GPIO_HIGH); err != nil { // Low == Enabled
			return err // FIXME do better - custom error and logging
		}
	}

	// Generate the series of waveforms we want to send
	waveChain, err := makeWaveChain(pigpio, dists)
	if err != nil {
		return err // FIXME Log, do better
	}

	// Set up a callback for any of the limit switches being hit

	// Transmit the waveChain
	_ = waveChain // FIXME implement

	// Wait for the waveform to be finished, or for a limit switch to be hit

	return nil
}

func makeWaveChain(pigpio io.ReadWriter, dists Distances) (gopigpio.Chainers, error) {
	// Get the pdis for the ramp up phase
	pdis := CalcPulseDelayIntervals(rampFunc, rampDuration, rampIntervals)
	pdis = append(pdis, steadyPDI) // steadyPDI.Count should be maxUint64 FIXME do something better to indicate infinity

	// Set the "current pulse width" to the first pdi in <pdis>
	currPDIIndex := 0
	// Set a counter for the remaining number of pulses in the current "pdi"
	currPDIPulseCount := pdis[currPDIIndex].count
	currPulseWidth := pdis[currPDIIndex].period

	// FIXME Implement Start the waveChain with a Wave that turns off any of the motors that have dists of 0

	// FIXME This whole "Chainer" interface feels clunky - particularly how a Slice of Chainers
	// itself is also a Chainer. I think it needs to be reworked or re-thought or something.
	chain := []gopigpio.Chainer{} // this is the final result we return - a slice of Chainers

	// While there's still some steps in <dists>:
	for ! dists.AllZero() {
		debug.Printf("\n----------------------------\n")
		debug.Printf("dists: %v\n", dists)
		debug.Printf("currPDIPulseCount: %d", currPDIPulseCount)
		debug.Printf("currPulseWidth: %d", currPulseWidth)

		// Update the currPDI if the pulse count has hit zero
		for currPDIPulseCount == 0 {
			debug.Printf("Recalulating currPDIPulseWidth")
			// advance the currPDIIndex, panic if overflow
			if currPDIIndex += 1; currPDIIndex == len(pdis) {
				panic("Overflow of PDI structure. Means we used up all the steps in the last PDI in the PDI set. This should never happen. Means at least one distance was bigger than uint64")
			}
			// reload currPulseWidth
			currPulseWidth = pdis[currPDIIndex].period
			// reload currPDIPulseCount
			currPDIPulseCount = pdis[currPDIIndex].count
			debug.Printf("currPDIPulseCount: %d", currPDIPulseCount)
			debug.Printf("currPulseWidth: %d", currPulseWidth)
		}

		// If any of the <dists> are about to Expire (ie, have a value of 1):
		// 	- Do a single step pulse, which incorporates a disable pulse for any of those pins
		//	- Then start a new iteration of the overall loop
		if uint64(dists.MinNotZero()) == 1 {
			debug.Printf("At least one dist is about to expire - making a special pulse")
			// FIXME IMPLEMENT
			disablePins := []gopigpio.Pin{}
			for i, d := range dists {
				if d == 1 {
					disablePins = append(disablePins, EnablePins[i])
				}

			}

			subChain, err := makeStepWithPinDisableChain(pigpio, currPulseWidth, disablePins)
			if err != nil {
				return gopigpio.Chainers(chain), err // FIXME log, do better etc
			}
			debug.Printf("SingleStepDisablechain: %v\n", subChain)

			// Add the new chain to the overall chain
			chain = append(chain, subChain)

			// Subtract 1 from various step counters
			currPDIPulseCount--
			for i := range dists {
				if dists[i] > 0 {
					dists[i]--
				}
			}

			// Start a new loop iteration
			continue
		}

		// Set "loopCount" to the minimum of currPDIPulseCount and dists.MintNotZero()-1
		// set NextLoopCount to currPDIPulseCount
		loopCount := uint64(dists.MinNotZero() - 1)
		debug.Printf("MinNotZero - 1: %d\n", loopCount)

		if currPDIPulseCount < loopCount {
			loopCount = currPDIPulseCount
		}
		debug.Printf("loopCount: %d\n", loopCount)

		// Create a loop based on the current PDI width and loop count
		subChain, err := makeStepLoopChain(pigpio, currPulseWidth, loopCount)
		if err != nil {
			return gopigpio.Chainers(chain), err // FIXME log, do better, etc	
		}
		debug.Printf("LoopNChain: %v\n", subChain)

		// Add the new chain to the overall chain
		chain = append(chain, subChain)

		// Subtract <loopcount> from various counters
		currPDIPulseCount -= loopCount
		for i := range dists {
			if dists[i] > 0 {
				dists[i] -= Distance(loopCount)
			}
		}
	}

	debug.Printf("chain result is: %v\n", chain)
	return gopigpio.Chainers(chain), nil
}

func makeStepWithPinDisableChain(pigpio io.ReadWriter, period time.Duration, disablePins []gopigpio.Pin) (gopigpio.ChainWaveID, error) {

	if err := addStepPulseToWave(pigpio, period); err != nil {
		return gopigpio.ChainWaveID(0), err
	}

	// Add the disable pulse to the current wave being built in pigpio
	pulseDelay := gopigpio.Pulse {
		OnPins: nil,
		OffPins: nil,
		Delay: stepPulseWidth + stepPulseWait + delayUntilDisable,
	}
	pulseDisable := gopigpio.Pulse {
		OnPins: disablePins, // Enable pin is active low, so turning a pin on DISABLES the stepper
		OffPins: nil,
		Delay: 0,
	}

	pulses := []gopigpio.Pulse{pulseDelay, pulseDisable}

	pulseCount, err := gopigpio.WaveAddGeneric(pigpio, pulses)
	if err != nil { // FIXME validate the pulse count we receive back
		return gopigpio.ChainWaveID(0), err // FIXME log, do better, whatever
	}
	debug.Printf("Number of Pulses in current wave is %d\n", pulseCount)

	// Create a wave from the pulses added so far
	waveID, err := gopigpio.WaveCreate(pigpio)
	if err != nil {
		return gopigpio.ChainWaveID(waveID), err // FIXME lo, do better, whatever
	}
	debug.Printf("WaveID is %d\n", waveID)

	return gopigpio.ChainWaveID(waveID), nil

}

func makeStepLoopChain(pigpio io.ReadWriter, period time.Duration, count uint64) (chain gopigpio.ChainLoopN, err error){

	// Clean out any previous waves
	/* COMMENTED OUT, BECAUSE THERE SHOULD NEVER BE ANY PREVIOUS WAVES, AND IF THERE ARE
	 IT'S A BuG SO FIX IT

	if err = gopigpio.WaveAddNew(pigpio); err != nil {
		return chain, err // FIXME log error
	}
	*/

	if err = addStepPulseToWave(pigpio, period); err != nil {
		return chain, err
	}

	// Create a wave from the pulses added so far
	waveID, err := gopigpio.WaveCreate(pigpio)
	if err != nil {
		return chain, err // FIXME lo, do better, whatever
	}
	debug.Printf("WaveID is %d\n", waveID)

	// FIXME this is crap, do better
	// Make sure we haven't overflowed the capability of our loop structure
	// To fix this, we need to embed multiple loops. But that's a hassle right now
	if count > math.MaxUint16 {
		panic("count is too big - max size is 2^16-1")
	}

	chain = gopigpio.ChainLoopN {
		Chain: gopigpio.ChainWaveID(waveID),
		Count: uint16(count),
	}

	return chain, nil
}

func addStepPulseToWave(pigpio io.ReadWriter, period time.Duration) error {
	pulseUp := gopigpio.Pulse {
		OnPins: []gopigpio.Pin{STEP_PIN},
		OffPins: nil,
		Delay: stepPulseWidth,
	}
	pulseDown := gopigpio.Pulse {
		OnPins: nil,
		OffPins: []gopigpio.Pin{STEP_PIN},
		Delay: period - stepPulseWidth,
	}

	pulses := []gopigpio.Pulse{pulseUp, pulseDown}

	pulseCount, err := gopigpio.WaveAddGeneric(pigpio, pulses)
	if err != nil { // FIXME validate the pulse count we receive back
		return err // FIXME log, do better, whatever
	}
	debug.Printf("Number of Pulses in current wave is %d\n", pulseCount)

	return nil
}

// PulseDelayInterval describes an interval, and a PulseDelay to use in that interval
// A series of PulseDelayIntervals can be used to describe a changing rate of pulses
type PulsePeriodInterval struct {
	period time.Duration // Time between pulses within this interval. Note this is from pulse-start to pulse-start
	count uint64	     // The number of pulses in this interval. Multiply this by <period> to get total duration
			     // 	of this PDI
			     // FIXME should count be uint? but if uint, how to represent infinity. With a zero?
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
func CalcPulseDelayIntervals(
	rateFunc func(float64) float64, // The rate function - returns steps/second vs seconds
	duration time.Duration, // Total duration of intervals
	intervals int, // Number of intervals to generate
) []PulsePeriodInterval {
	var pdis []PulsePeriodInterval

	intervalWidthNanos := float64(duration.Nanoseconds()) / float64(intervals)
	fmt.Printf("IntervalWidthNanos is %f\n", intervalWidthNanos)
	pulsesSoFar := 0

	intervalStartNanos := 0.0 // At each loop iteration, this will get set to the the end of the previous interval
	for i := 0; i < intervals; i++ {
		// Calculate the exact end of the interval
		intervalEndNanos := intervalWidthNanos * float64(i + 1)

		intervalTargetWidthNanos := intervalEndNanos - intervalStartNanos

		// Calculate the total number of steps we want to have achieved by the end of this interval
		pulsesToReachFloat := calculus.AntiDiff(rateFunc, 0, intervalEndNanos / 1000000000)

		// Calculate how many pulses to send in this interval
		pulsesInIntervalFloat := pulsesToReachFloat - float64(pulsesSoFar)

		// Calculate ideal period for pulses
		pulseIdealPeriod := intervalTargetWidthNanos / pulsesInIntervalFloat

		// Round the ideal period to nearest microsecond
		pulsePeriodMicros := int(pulseIdealPeriod / 1000 + 0.5) // +0.5 causes it to round rather than truncate

		// Calculate exact number of pulses needed based on the rounded PulsePeriodMicros
		pulsesInIntervalExact := intervalTargetWidthNanos / float64(pulsePeriodMicros * 1000)

		// Round the number of pulses
		pulsesInIntervalRounded := int(pulsesInIntervalExact + 0.5)

		// Calculate the 'adjusted' interval width, which is the product of the integer number of pulses
		// and the integer period
		adjustedIntervalWidthNanos := pulsesInIntervalRounded * pulsePeriodMicros * 1000

		pdis = append(pdis, PulsePeriodInterval{
			period: time.Duration(pulsePeriodMicros) * time.Microsecond,
			count: uint64(pulsesInIntervalRounded),
		})

		// Update values for next iteration
		intervalStartNanos = intervalStartNanos + float64(adjustedIntervalWidthNanos)
		pulsesSoFar = pulsesSoFar + pulsesInIntervalRounded
	}

	return pdis
}

func TranslateRange(val, inMin, inMax, outMin, outMax float64) float64 {
	inRange := inMax - inMin

	// Scale val to a value between 0 and 1
	val = (val - inMin) / inRange

	outRange := outMax - outMin
	// Now scale val to a value in the output range
	val = val * outRange + outMin

	return val
}

// wrapFunction wraps a function such that the inputs and outputs have different ranges
func WrapFunction(f func(float64) float64, oldMin, oldMax, newMin, newMax, outMin, outMax float64) (func(float64) float64) {
	oldOutMin := f(oldMin)
	oldOutMax := f(oldMax)

	wrapped := func(x float64) float64 {
		newX := TranslateRange(x, newMin, newMax, oldMin, oldMax)
		val := f(newX)
		return TranslateRange(val, oldOutMin, oldOutMax, outMin, outMax)
	}

	return wrapped
}
