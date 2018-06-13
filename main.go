package main

import (
	//"context"
	"fmt"
	"time"

	//"github.com/google/subcommands"
	//"github.com/phlipped/paint-dispenser/controller"
)

const (
	j int = 100
	k int = 100000
)
func main() {
	fmt.Printf("Starting Paint Dispenser\n")
	//  controller.PrintGPIOStatus()
	//controller.ToggleAllPins()
        // controller.FlashLed()
	times := [j]time.Time{}
	for i := 0; i < j; i++ {
		times[i] = time.Now()
	}
	for i := 0; i < j-1; i++ {
		fmt.Printf("%d\n", times[i].Nanosecond() - times[i+1].Nanosecond())
	}

	start := time.Now()
	vals := [k]int{}
	for i := 0; i < k; i++ {
		vals[i] = i
	}
	end := time.Now()
	fmt.Printf("%v\n%v\n", start, end)
}
