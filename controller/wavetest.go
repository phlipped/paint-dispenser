package controller

import (
	"fmt"
	"net"
	"time"

	"github.com/phlipped/gopigpio"
)

const (
	TEST_PIN gopigpio.Pin = STEP_PIN
)

func WaveTest() {
	p, err := net.Dial("tcp", "localhost:8888")
	if err != nil {
		panic(err)
	}

	if err := gopigpio.GpioSetMode(p, TEST_PIN, gopigpio.OUTPUT); err != nil {
		panic(err)
	}

	gopigpio.WaveClear(p)
	gopigpio.WaveAddNew(p)


	pulses := []gopigpio.Pulse{
		{[]gopigpio.Pin{TEST_PIN}, []gopigpio.Pin{}, time.Microsecond * 10},
		{[]gopigpio.Pin{}, []gopigpio.Pin{TEST_PIN}, time.Microsecond * 5},
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

