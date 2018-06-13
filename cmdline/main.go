package main

import (
	"context"
	"flag"
	"github.com/google/subcommands"
	"os"
)

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(&statusCmd{}, "")
	subcommands.Register(&pinTestCmd{}, "")

	flag.Parse()
	ctx := context.Background()
        os.Exit(int(subcommands.Execute(ctx)))

}
