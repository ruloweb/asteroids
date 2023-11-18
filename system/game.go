package system

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/ruloweb/asteroids/entity"
)

const (
	ScreenW float64 = 320
	ScreenH float64 = 480
)

var (
	debugMessage string
)

type Game struct {
	drawUI           bool
	debug            bool
	ticksLimit       uint64
	maxAsteroids     uint
	ship             *entity.Ship
	asteroids        []*entity.Asteroid
	ticks            uint64
	AsteroidsCrashed uint64
	grid             *entity.Grid
	ai               *AI
}

func NewGame(drawUI bool, debug bool, ticksLimit uint64, maxAsteroids uint, ship *entity.Ship, grid *entity.Grid, ai *AI) *Game {
	return &Game{drawUI: drawUI, debug: debug, ticksLimit: ticksLimit, maxAsteroids: maxAsteroids, ship: ship, asteroids: []*entity.Asteroid{},
		ticks: 0, AsteroidsCrashed: 0, grid: grid, ai: ai}
}

func (g *Game) Update() error {
	g.ticks++
	if g.ticks >= g.ticksLimit {
		return ebiten.Termination
	}

	// Reset the grid
	g.grid.Reset()

	// Create new asteroid
	if g.ticks%(30+uint64(rand.Intn(90))) == 0 {
		g.createAsteroid()
	}

	for _, a := range g.asteroids {
		// Update asteroids position
		a.Move(2)
		// Update grid asteroids
		g.grid.UpdateGridAsteroid(a)
	}

	// Update grid ship
	g.grid.UpdateGridShip(g.ship)

	// Check out of screen
	for i, a := range g.asteroids {
		_, y := a.CornerBottomLeft()
		if y > ScreenH {
			g.deleteAsteroid(i)
		}
	}

	// Collision detection
	for i, a := range g.asteroids {
		if a.CheckCollision(g.ship) {
			g.deleteAsteroid(i)
			g.AsteroidsCrashed++
		}
	}

	// Game over
	/*if g.AsteroidsCrashed >= 100 {
		return ebiten.Termination
	}*/

	// Move ship
	x, _ := g.ship.CornerTopLeft()
	if g.ai != nil {
		m := g.ai.GetMove(g.grid)
		if m == 1 && x > float64(g.ship.Width()) {
			g.ship.Move(-2)
		}
		if m == 2 && x < ScreenW-float64(g.ship.Width()) {
			g.ship.Move(2)
		}
	}
	if inpututil.KeyPressDuration(ebiten.KeyArrowRight) > 0 && x < ScreenW-float64(g.ship.Width()) {
		g.ship.Move(2)
	}
	if inpututil.KeyPressDuration(ebiten.KeyArrowLeft) > 0 && x > 0 {
		g.ship.Move(-2)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if !g.drawUI {
		return
	}

	// Draw background
	screen.Fill(color.RGBA{A: 0xff})
	ebitenutil.DebugPrint(screen,
		fmt.Sprintf(
			"Ticks: %d, # asteroids: %d, Chashes: %d\n%s",
			g.ticks, len(g.asteroids), g.AsteroidsCrashed, debugMessage))

	// Draw ship
	g.ship.Draw(screen)

	// Draw asteroids
	for _, a := range g.asteroids {
		a.Draw(screen)
	}

	// Debug
	if g.debug {
		drawBoxes(screen, g)
		drawGrid(screen, g)
	}
}

func drawBoxes(screen *ebiten.Image, g *Game) {
	for _, a := range g.asteroids {
		x, y := a.CornerTopLeft()
		vector.DrawFilledRect(
			screen,
			float32(x), float32(y),
			float32(a.Width()), float32(a.Height()),
			color.RGBA{R: 100, G: 100, A: 1}, false)
	}

	x, y := g.ship.CornerTopLeft()
	vector.DrawFilledRect(screen, float32(x), float32(y), float32(g.ship.Width()), float32(g.ship.Height()), color.RGBA{R: 100, G: 100, A: 1}, false)
}

func drawGrid(screen *ebiten.Image, g *Game) {
	g.grid.DrawGridAsteroid(screen)
	g.grid.DrawGridShip(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return int(ScreenW), int(ScreenH)
}

func (g *Game) createAsteroid() {
	if len(g.asteroids) >= int(g.maxAsteroids) {
		return
	}
	b := entity.ImgAsteroid.Bounds()
	newAsteroid := entity.NewAsteroid(rand.Float64()*ScreenW, -float64(b.Dy()), math.Pi*rand.Float64())
	g.asteroids = append(g.asteroids, newAsteroid)
}

func (g *Game) deleteAsteroid(i int) {
	if i < len(g.asteroids) {
		g.asteroids[i], g.asteroids = g.asteroids[len(g.asteroids)-1], g.asteroids[:len(g.asteroids)-1]
	}
}
