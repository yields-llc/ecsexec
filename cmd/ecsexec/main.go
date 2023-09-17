package main

import (
	"context"
	"flag"
	"os"

	"github.com/google/subcommands"
	"github.com/yields-llc/ecsexec/internal/commands"
)

func main() {
	subcommands.Register(subcommands.HelpCommand(), "help")
	subcommands.Register(commands.NewStartCommand(), "ecsexec")
	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
