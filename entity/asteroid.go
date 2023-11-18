package entity

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"log"
)

var (
	ImgAsteroid *ebiten.Image
)

type Asteroid struct {
	img                  *ebiten.Image
	posX, posY, rotation float64
}

func init() {
	var err error
	img, err := assets.Open("assets/asteroid.png")
	ImgAsteroid, _, err = ebitenutil.NewImageFromReader(img)
	if err != nil {
		log.Fatal(err)
	}
}

func NewAsteroid(posX float64, posY float64, rotation float64) *Asteroid {
	return &Asteroid{
		img:      ImgAsteroid,
		posX:     posX,
		posY:     posY,
		rotation: rotation,
	}
}

func (a *Asteroid) Width() int {
	return a.img.Bounds().Dx()
}

func (a *Asteroid) Height() int {
	return a.img.Bounds().Dy()
}

func (a *Asteroid) CornerTopLeft() (float64, float64) {
	x := a.posX - float64(a.Width())/2
	y := a.posY - float64(a.Height())/2
	return x, y
}

func (a *Asteroid) CornerTopRight() (float64, float64) {
	x := a.posX - float64(a.Width())/2 + float64(a.Width())
	y := a.posY - float64(a.Height())/2
	return x, y
}

func (a *Asteroid) CornerBottomLeft() (float64, float64) {
	x := a.posX - float64(a.Width())/2
	y := a.posY - float64(a.Height())/2 + float64(a.Height())
	return x, y
}

func (a *Asteroid) CornerBottomRight() (float64, float64) {
	x := a.posX - float64(a.Width())/2 + float64(a.Width())
	y := a.posY - float64(a.Height())/2 + float64(a.Height())
	return x, y
}

func (a *Asteroid) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	s := a.img.Bounds().Size()
	op.GeoM.Reset()
	op.GeoM.Translate(-float64(s.X)/2, -float64(s.Y)/2)
	op.GeoM.Rotate(a.rotation)
	op.GeoM.Translate(a.posX, a.posY)
	screen.DrawImage(a.img, op)
}

func (a *Asteroid) Move(n float64) {
	a.posY += n
	a.rotation += 0.1
}

func (a *Asteroid) CheckCollision(ship *Ship) bool {
	shipX0, shipY0 := ship.CornerTopLeft()
	shipX1, _ := ship.CornerTopRight()

	x0, y0 := a.CornerBottomLeft()
	x1, _ := a.CornerBottomRight()

	return y0 >= shipY0 && ((x0 >= shipX0 && x0 <= shipX1) || (x1 >= shipX0 && x1 <= shipX1))
}
