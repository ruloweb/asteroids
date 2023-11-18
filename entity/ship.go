package entity

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"log"
)

var (
	imgShip *ebiten.Image
)

type Ship struct {
	img        *ebiten.Image
	posX, posY float64
}

func init() {
	var err error
	img, err := assets.Open("assets/ship.png")
	imgShip, _, err = ebitenutil.NewImageFromReader(img)
	if err != nil {
		log.Fatal(err)
	}
}

func NewShip(posX float64, posY float64) *Ship {
	return &Ship{
		img:  imgShip,
		posX: posX,
		posY: posY,
	}
}

func (s *Ship) Width() int {
	return s.img.Bounds().Dx()
}

func (s *Ship) Height() int {
	return s.img.Bounds().Dy()
}

func (s *Ship) CornerTopLeft() (float64, float64) {
	x := s.posX
	y := s.posY
	return x, y
}

func (s *Ship) CornerTopRight() (float64, float64) {
	x := s.posX + float64(s.Width())
	y := s.posY
	return x, y
}

func (s *Ship) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(s.posX, s.posY)
	screen.DrawImage(s.img, op)
}

func (s *Ship) Move(n float64) {
	s.posX += n
}
