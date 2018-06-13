package controller

import (
	"github.com/stianeikeland/go-rpio"
)

type pinInit struct {
	name string
	direction rpio.Mode
	value rpio.State
}

// Map Paint Dispenser Controller names to the corresponding Pin Names used in Embd
var pinInits = []pinInit {
	{"dir", rpio.Output, rpio.Low},
	{"step", rpio.Output, rpio.Low},
	{"en1", rpio.Output, rpio.High},
	{"en2", rpio.Output, rpio.High},
	{"en3", rpio.Output, rpio.High},
	{"en4", rpio.Output, rpio.High},
	{"en5", rpio.Output, rpio.High},
	{"lim1U", rpio.Input, rpio.High},
	{"lim1L", rpio.Input, rpio.High},
	{"lim2U", rpio.Input, rpio.High},
	{"lim2L", rpio.Input, rpio.High},
	{"lim3U", rpio.Input, rpio.High},
	{"lim3L", rpio.Input, rpio.High},
	{"lim4U", rpio.Input, rpio.High},
	{"lim4L", rpio.Input, rpio.High},
	{"lim5U", rpio.Input, rpio.High},
	{"lim5L", rpio.Input, rpio.High},
}

var EnablePinNames = []string{"en1", "en2", "en3", "en4", "en5"}
var LimitUpperPinNames = []string{"lim1U", "lim2U", "lim3U", "lim4U", "lim5U"}
var LimitLowerPinNames = []string{"lim1L", "lim2L", "lim3L", "lim4L", "lim5L"}

var Pins = map[string]rpio.Pin{
	"dir": 3,
	"step": 18,
	"en1": 23,
	"en2": 2,
	"en3": 19,
	"en4": 26,
	"en5": 21,
	"lim1U": 4,
	"lim1L": 17,
	"lim2U": 27,
	"lim2L": 22,
	"lim3U": 10,
	"lim3L": 9,
	"lim4U": 11,
	"lim4L": 5,
	"lim5U": 6,
	"lim5L": 13,
}

func initPins() {
	for _, pinInit := range pinInits {
		pin := Pins[pinInit.name]
		pin.Mode(pinInit.direction)
		if pinInit.direction == rpio.Output {
			pin.Write(pinInit.value)
		} else { // else assume input
			if pinInit.value == rpio.High {
				pin.PullUp()
			} else { // else assume pull down
				pin.PullDown()
			}
		}
	}
}
