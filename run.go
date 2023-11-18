package main

import (
	"log"
	"os"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ruloweb/asteroids/entity"
	"github.com/ruloweb/asteroids/system"
)

type RunCmd struct {
	NoUI         bool   `help:"No draw UI."`
	Debug        bool   `default:"0" help:"Enable debug mode."`
	Ticks        int    `default:"60" help:"Ticks per second."`
	MaxAsteroids uint   `default:"50" help:"Max number of asteroids at the same time."`
	TicksLimit   uint64 `default:"65535" help:"Max number of ticks to play."`
	Genome       string `default:"" help:"Genome file path."`
	GridWidth    uint   `default:"4" help:"Number of columns in the grid."`
	GridHeight   uint   `default:"2" help:"Number of rows in the grid."`
	OutputFile   string `default:"" help:"Output file for the final ticks count"`
}

func (r *RunCmd) Run() error {
	// Initialize game
	ebiten.SetTPS(r.Ticks)
	ebiten.SetWindowSize(int(system.ScreenW), int(system.ScreenH))
	ebiten.SetWindowTitle("Asteroids")

	// Create ship
	ship := entity.NewShip(160, 440)

	// Genome
	var ai *system.AI
	if r.Genome != "" {
		var err error
		ai, err = system.NewAIFromPath(r.Genome)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Run game
	grid := entity.NewGrid(r.GridWidth, r.GridHeight, system.ScreenW, system.ScreenH)
	g := system.NewGame(!r.NoUI, r.Debug, r.TicksLimit, r.MaxAsteroids, ship, grid, ai)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}

	// Save output
	if r.OutputFile != "" {
		err := os.WriteFile(r.OutputFile, []byte(strconv.FormatUint(g.AsteroidsCrashed, 10)), 0644)
		if err != nil {
			log.Fatal(err)
		}
	}

	return nil
}
