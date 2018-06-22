package controller

type GpioPin uint8

const(
	DIR_PIN GpioPin = 3
	STEP_PIN GpioPin = 18
	EN1_PIN GpioPin = 23
	EN2_PIN GpioPin = 2
	EN3_PIN GpioPin = 19
	EN4_PIN GpioPin = 26
	EN5_PIN GpioPin = 21
	LIM_1U_PIN GpioPin = 4
	LIM_1L_PIN GpioPin = 17
	LIM_2U_PIN GpioPin = 27
	LIM_2L_PIN GpioPin = 22
	LIM_3U_PIN GpioPin = 10
	LIM_3L_PIN GpioPin = 9
	LIM_4U_PIN GpioPin = 11
	LIM_4L_PIN GpioPin = 5
	LIM_5U_PIN GpioPin = 6
	LIM_5L_PIN GpioPin = 13
)
var pinInits = []map[GpioPin][]struct{uint32, uint32} {
	{"dir", gopigpio.OUTPUT, gopigpio.GPIO_LOW},
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
