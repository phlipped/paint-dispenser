package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/subcommands"
	"github.com/phlipped/paint-dispenser/controller"
)

type statusCmd struct {}

func (s *statusCmd) Name() string { return "status" }
func (s *statusCmd) Synopsis() string { return "Print status of dispenser." }
func (s *statusCmd) Usage() string {
	return `status:
	Print status of dispenser.
`
}

func (s *statusCmd) SetFlags(_ *flag.FlagSet) {}

func (s *statusCmd) Execute(_ context.Context, _ *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	status, err := controller.GetStatus()
	if err != nil {
		return subcommands.ExitFailure
	}
	fmt.Printf("%v\n", status)
	return subcommands.ExitSuccess
}
