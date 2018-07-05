package main

import (
	"fmt"
	"net"
	"time"

	"github.com/phlipped/paint-dispenser/controller"


)

func yEqualsX(x float64) float64 {
	return x
}

func main() {
	wrappedFunc := controller.WrapFunction(yEqualsX, 0.0, 1.0, 0.0, 1.0, 0.0, 50000.0)
	pdis := controller.CalcPulseDelayIntervals(wrappedFunc, time.Duration(1000000000), 10)
	for _, pdi := range pdis {
		fmt.Printf("%v\n", pdi)
	}

	s, err := net.Dial("tcp", "localhost:8888")
	if err != nil {
		panic(err)
	}

	dists := controller.Distances{0, 10, 15, 200, 10000}
	err = controller.DoSteps(s, dists, controller.Down)
	if err != nil {
		panic(err)
	}
}
