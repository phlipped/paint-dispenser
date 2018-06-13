package main

import (
//	"fmt"
	"context"
	"flag"
	"sync"
//	"time"

	"github.com/google/subcommands"
	"github.com/phlipped/paint-dispenser/controller"
	"github.com/stianeikeland/go-rpio"
)

type pinTestCmd struct {}

func (s *pinTestCmd) Name() string { return "pintest" }
func (s *pinTestCmd) Synopsis() string { return "Test the GPIO pins" }
func (s *pinTestCmd) Usage() string {
	return `pintest:
	Runs a test pattern on output GPIO pins, and prints the status of the input pins.
`
}

func (s *pinTestCmd) SetFlags(_ *flag.FlagSet) {}

func (s *pinTestCmd) Execute(_ context.Context, _ *flag.FlagSet, _ ...interface{}) (status subcommands.ExitStatus) {
	//timeout := time.NewTicker(time.Second * 10)
	cleanup, err := controller.Init()
	if err != nil {
		return subcommands.ExitFailure
	}
	defer func() {
		if err := cleanup(); err != nil {
			status = subcommands.ExitFailure
		}
	}()


	// Set up a 1 us square wave on the direction, step and enable pins.
	// Run them as go routines and coordinate them using a channel that sends
	// pulses to them
//	triggerChans := []chan struct{}{}
//	outputPins := []string{"step", "dir"}
//	outputPins = append(outputPins, controller.EnablePins...)
//	wg := sync.WaitGroup{} // Used to indicate receipt of trigger
//	wg.Add(len(outputPins))
//	for _, pin := range outputPins {
//		triggerChan := make(chan struct{})
//		triggerChans = append(triggerChans, triggerChan)
//		go runSquareWave(controller.Pins[pin], triggerChan, wg)
//	}

	// Run a ticker and trigger each trigger channel at each tick
//	ticker := time.NewTicker(time.Microsecond)
//	fmt.Printf("hi mum")
//	for {
//		select {
//		case t := <-timeout:
//			break
//		case tick := <-ticker:
//			for _, trigger := range triggerChans {
//				trigger <- struct{}{}
//			}
//			// Wait for everyone to have received the trigger
//			wg.Wait()
//			// And then reset the wait group
//			wg.Add(len(outputPins))
//		}
//	}
//	for _, tc := range triggerChans {
//		close(tc)
//	}


	// Close the trigger channels, wait for all the goroutines to finish

	//outputPins := []string{"step", "dir"}
	//outputPins = append(outputPins, controller.EnablePins...)

	//microTicker := time.NewTicker(time.Nanosecond * 1)
	//endTime := time.Now().Add(time.Second * 30)
	//for ; time.Now().Before(endTime); {
	dir := controller.Pins["dir"]
	//step := controller.Pins["step"]

	for {
		//_ = <-microTicker.C
		rpio.WritePin(dir, rpio.High)
		rpio.WritePin(dir, rpio.Low)
		//for _, p := range outputPins {
		//	controller.Pins[p].Toggle()
		//}
	}

	return subcommands.ExitSuccess
}

func runSquareWave(pin rpio.Pin, trigger <-chan struct{}, wg sync.WaitGroup) {
	for _ = range trigger {
		wg.Done()
		pin.Toggle()
	}
}

