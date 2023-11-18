package entity

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Grid struct {
	width, height uint
	cells         uint64
	cellsShip     uint8
	screenW       float64
	screenH       float64
}

func NewGrid(width uint, height uint, screenW float64, screenH float64) *Grid {
	return &Grid{width: width, height: height, screenW: screenW, screenH: screenH}
}

func (g *Grid) Reset() {
	g.cells = 0
	g.cellsShip = 0
}

func (g *Grid) UpdateGridAsteroid(a *Asteroid) {
	g.updateGridAsteroid(a.CornerBottomLeft())
	g.updateGridAsteroid(a.CornerBottomRight())
	g.updateGridAsteroid(a.CornerTopLeft())
	g.updateGridAsteroid(a.CornerTopRight())
}

func (g *Grid) updateGridAsteroid(x, y float64) {
	ratioX := g.screenW / float64(g.width)
	ratioY := g.screenH / float64(g.height)
	pos := uint(y/ratioY)*g.width + uint(x/ratioX)
	if x >= 0 && x < g.screenW && y >= 0 && y < g.screenH {
		g.cells |= 1 << pos
	}
}

func (g *Grid) UpdateGridShip(s *Ship) {
	ratioX := g.screenW / float64(g.width)
	x, _ := s.CornerTopLeft()
	pos := uint(x / ratioX)
	if x >= 0 && x < g.screenW {
		g.cellsShip |= 1 << pos
	}

	x, _ = s.CornerTopRight()
	pos = uint(x / ratioX)
	if x >= 0 && x < g.screenW {
		g.cellsShip |= 1 << pos
	}
}

func (g *Grid) DrawGridAsteroid(screen *ebiten.Image) {
	ratioX := g.screenW / float64(g.width)
	ratioY := g.screenH / float64(g.height)
	for i := 0; i < 64; i++ {
		if g.cells&(1<<uint(i)) != 0 {
			x := float64(i%int(g.width)) * ratioX
			y := float64(i/int(g.width)) * ratioY
			vector.DrawFilledRect(screen, float32(x), float32(y), float32(ratioX), float32(ratioY), color.RGBA{R: 70, A: 1}, false)
		}
	}
}

func (g *Grid) DrawGridShip(screen *ebiten.Image) {
	ratioX := g.screenW / float64(g.width)
	ratioY := g.screenH / float64(g.height)
	for i := 0; i < 8; i++ {
		if g.cellsShip&(1<<uint(i)) != 0 {
			x := float64(i%int(g.width)) * ratioX
			y := g.screenH - ratioY
			vector.DrawFilledRect(screen, float32(x), float32(y), float32(ratioX), float32(ratioY), color.RGBA{G: 70, A: 1}, false)
		}
	}
}

func (g *Grid) GetCellAsteroid() uint64 {
	return g.cells
}

func (g *Grid) GetCellShip() uint8 {
	return g.cellsShip
}
