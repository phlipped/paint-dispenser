package controller

import (
	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/rpi"
	"time"
	"fmt"
)

func PrintGPIOStatus() error {
	fmt.Printf("Printing GPIO PinMap:\n")
	var host *embd.Descriptor
	var err error

	if host, err = embd.DescribeHost(); err != nil {
		return err
	}
	fmt.Printf("(got host descriptor)\n")

	// Ok now we've got the host descriptor - what do?
	// we can get a GPIODriver out of it, right?
	gpio_driver := host.GPIODriver()

	// Print out the GPIO pin map, perhaps ...
	for i, PinDesc := range gpio_driver.PinMap() {
			fmt.Printf("%d: %v\n", i, *PinDesc)
	}

	return nil
}

func ToggleAllPins() {
	if err := embd.InitGPIO(); err != nil {
		panic(err)
	}

	var host *embd.Descriptor
	var err error

	if host, err = embd.DescribeHost(); err != nil {
		fmt.Printf("Failed to get host description: %v\n", err)
		return
	}

	gpio_driver := host.GPIODriver()

	pins := []embd.DigitalPin{}
	for _, pin_desc := range gpio_driver.PinMap() {
		pin, err := embd.NewDigitalPin(pin_desc.DigitalLogical)
		if err != nil {
			fmt.Printf("failed to open pin %d: %v\n", pin_desc.DigitalLogical, err)
			return
		}
		defer pin.Close()

		if err = pin.SetDirection(embd.Out); err != nil {
			fmt.Printf("failed to set direction on pin %d: %v\n", pin_desc.DigitalLogical, err)
			return
		}

		pins = append(pins, pin)
	}

	for i := 0; i < 10000; i ++ {
		fmt.Printf("%d\n", i)

		for _, pin := range(pins) {
			if err := pin.Write(embd.High); err != nil {
				fmt.Printf("failed to Write HIGH to pin %d: %v\n", pin.N(), err)
			}

		}

		time.Sleep(time.Millisecond * 250)

		for _, pin := range(pins) {
			if err := pin.Write(embd.Low); err != nil {
				fmt.Printf("failed to Write LOW to pin %d: %v\n", pin.N(), err)
			}

		}

		time.Sleep(time.Millisecond * 250)
	}

}

func FlashLed() {
	ledPin, err := embd.NewDigitalPin("P1_11")
	defer ledPin.Close()
	fmt.Printf("pin N is %d\n", ledPin.N())
	if err != nil {
		fmt.Printf("failed to open pin\n")
		return
	}

	if err = ledPin.SetDirection(embd.Out); err != nil {
		fmt.Printf("failed to Set direction on pin\n")
		return
	}
	for i := 0; i < 50; i++ {
		fmt.Printf("%d\n", i)
		if err := ledPin.Write(1); err != nil {
			fmt.Printf("failed to pullup\n")
		}
		time.Sleep(time.Second * 1)
		if err := ledPin.Write(0); err != nil {
			fmt.Printf("failed to pulldown\n")
		}
		time.Sleep(time.Second * 1)
	}
}
