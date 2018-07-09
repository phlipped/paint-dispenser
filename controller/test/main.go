package main

import (
	"net"

	"github.com/phlipped/paint-dispenser/controller"


)

func yEqualsX(x float64) float64 {
	return x
}

func main() {
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
