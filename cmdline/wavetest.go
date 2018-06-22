package main

import (
	"context"
	"flag"
	"github.com/google/subcommands"
	"github.com/phlipped/paint-dispenser/controller"

)


type wavetestCmd struct {}

func (w *wavetestCmd) Name() string { return "wavetest" }
func (w *wavetestCmd) Synopsis() string { return "Test using the DMA waveform stuff." }
func (w *wavetestCmd) Usage() string {
	return `wavetest:
	Test using pigpio waveform stuff.
`
}

func (w *wavetestCmd) SetFlags(_ *flag.FlagSet) {}

func (w *wavetestCmd) Execute(_ context.Context, _ *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	controller.WaveTest()
	return subcommands.ExitSuccess
}
