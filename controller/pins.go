package controller

// Map Paint Dispenser Controller names to the corresponding Pin Names used in Embd
const (
	PinDir = "PI_"
	PinStep = "PI_"
	PinEn1 = "PI_"
	PinEn2 = "PI_"
	PinEn3 = "PI_"
	PinEn4 = "PI_"
	PinEn5 = "PI_"
	PinLim1U = "P1_"
	PinLim1L = "P1_"
	PinLim2U = "P1_"
	PinLim2L = "P1_"
	PinLim3U = "P1_"
	PinLim3L = "P1_"
	PinLim4U = "P1_"
	PinLim4L = "P1_"
	PinLim5U = "P1_"
	PinLim5L = "P1_"
)

var PinEn = []string{PinEn1, PinEn2, PinEn3, PinEn4, PinEn5}
var PinLimU = []string{PinLim1U, PinLim2U, PinLim3U, PinLim4U, PinLim5U}
var PinLimL = []string{PinLim1L, PinLim2L, PinLim3L, PinLim4L, PinLim5L}
