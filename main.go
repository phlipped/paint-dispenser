package main

import (
	"fmt"

	"github.com/phlipped/paint-dispenser/controller"
)

func main() {
	fmt.Printf("Starting Paint Dispenser\n")
//  controller.PrintGPIOStatus()
//controller.ToggleAllPins()
controller.FlashLed()
}
