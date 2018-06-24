package controller

import (
	"io"
	"github.com/phlipped/gopigpio"
)

const(
	DIR_PIN gopigpio.Pin = 3
	STEP_PIN gopigpio.Pin = 18
	EN1_PIN gopigpio.Pin = 23
	EN2_PIN gopigpio.Pin = 2
	EN3_PIN gopigpio.Pin = 19
	EN4_PIN gopigpio.Pin = 26
	EN5_PIN gopigpio.Pin = 21
	LIM_1U_PIN gopigpio.Pin = 4
	LIM_1L_PIN gopigpio.Pin = 17
	LIM_2U_PIN gopigpio.Pin = 27
	LIM_2L_PIN gopigpio.Pin = 22
	LIM_3U_PIN gopigpio.Pin = 10
	LIM_3L_PIN gopigpio.Pin = 9
	LIM_4U_PIN gopigpio.Pin = 11
	LIM_4L_PIN gopigpio.Pin = 5
	LIM_5U_PIN gopigpio.Pin = 6
	LIM_5L_PIN gopigpio.Pin = 13
)

type pinInit struct {
	mode gopigpio.PinMode
	val gopigpio.PinVal
}

var pinInits = map[gopigpio.Pin]pinInit {
	DIR_PIN: pinInit{gopigpio.OUTPUT, gopigpio.GPIO_LOW},
	STEP_PIN: pinInit{gopigpio.OUTPUT, gopigpio.GPIO_LOW},
	EN1_PIN: pinInit{gopigpio.OUTPUT, gopigpio.GPIO_HIGH},
	EN2_PIN: pinInit{gopigpio.OUTPUT, gopigpio.GPIO_HIGH},
	EN3_PIN: pinInit{gopigpio.OUTPUT, gopigpio.GPIO_HIGH},
	EN4_PIN: pinInit{gopigpio.OUTPUT, gopigpio.GPIO_HIGH},
	EN5_PIN: pinInit{gopigpio.OUTPUT, gopigpio.GPIO_HIGH},
	LIM_1U_PIN: pinInit{mode: gopigpio.INPUT},
	LIM_1L_PIN: pinInit{mode: gopigpio.INPUT},
	LIM_2U_PIN: pinInit{mode: gopigpio.INPUT},
	LIM_2L_PIN: pinInit{mode: gopigpio.INPUT},
	LIM_3U_PIN: pinInit{mode: gopigpio.INPUT},
	LIM_3L_PIN: pinInit{mode: gopigpio.INPUT},
	LIM_4U_PIN: pinInit{mode: gopigpio.INPUT},
	LIM_4L_PIN: pinInit{mode: gopigpio.INPUT},
	LIM_5U_PIN: pinInit{mode: gopigpio.INPUT},
	LIM_5L_PIN: pinInit{mode: gopigpio.INPUT},
}

var EnablePins = []gopigpio.Pin{EN1_PIN, EN2_PIN, EN3_PIN, EN4_PIN, EN5_PIN}
var LimitUpperPins = []gopigpio.Pin{LIM_1U_PIN, LIM_2U_PIN, LIM_3U_PIN, LIM_4U_PIN, LIM_5U_PIN}
var LimitLowerPins = []gopigpio.Pin{LIM_1L_PIN, LIM_2L_PIN, LIM_3L_PIN, LIM_4L_PIN, LIM_5L_PIN}


func initPins(p io.ReadWriter) {
	for pin, init := range pinInits {
		gopigpio.GpioSetMode(p, pin, init.mode)
		switch init.mode {
		case gopigpio.OUTPUT:
			gopigpio.GpioWrite(p, pin, init.val)
		case gopigpio.INPUT:
			//FIXME set pull direction base on init.val
		default:
			panic("unsupported pin mode")
		}
	}
}
