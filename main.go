package main

import (
	"github.com/alecthomas/kong"
)

var cli struct {
	Run   RunCmd   `cmd:"" help:"Run the game."`
	Train TrainCmd `cmd:"" help:"Train the genetic algorithm."`
}

func main() {
	ctx := kong.Parse(&cli)
	// Call the Run() method of the selected parsed command.
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
