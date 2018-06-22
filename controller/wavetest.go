package controller

import (
	"fmt"
	"net"
	"time"

	"github.com/phlipped/gopigpio"
)

const (
	pin = 18
)

func WaveTest() {
	p, err := net.Dial("tcp", "localhost:8888")
	if err != nil {
		panic(err)
	}

	gopigpio.GpioSetMode(p, pin, gopigpio.OUTPUT)

	gopigpio.WaveClear(p)
	gopigpio.WaveAddNew(p)

	if err := gopigpio.GpioSetMode(p, pin, gopigpio.OUTPUT); err != nil {
		panic(err)
	}

	pulses := []gopigpio.Pulse{
		{[]uint{18}, []uint{}, time.Microsecond * 10},
		{[]uint{}, []uint{18}, time.Microsecond * 5},
	}
	result, err := gopigpio.WaveAddGeneric(p, pulses)
	if err != nil {
		panic(err)
	}
	fmt.Printf("pulses is %d\n", result)

	waveID, err := gopigpio.WaveCreate(p)
	if err != nil {
		panic(err)
	}
	fmt.Printf("WaveID is %d\n", waveID)

	waveChain := gopigpio.ChainWaveID(waveID)
	withLoop := gopigpio.ChainLoopN{waveChain, 3}
	withDelay := gopigpio.Chainers([]gopigpio.Chainer{withLoop, gopigpio.ChainDelay{14 * time.Microsecond}})
	chainLoop := gopigpio.ChainLoopForever{withDelay}

	res, err := gopigpio.WaveChain(p, chainLoop)
	fmt.Printf("Result of wave chain command is %d\n", res)
}

