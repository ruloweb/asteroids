package system

import (
	"github.com/ruloweb/asteroids/entity"
	"math"
	"os"
	"strconv"
)

type AI struct {
	width  uint
	height uint
	genome []uint8
}

func NewAIFromPath(path string) (*AI, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Width
	b := make([]byte, 1)
	_, err = f.Read(b)
	if err != nil {
		return nil, err
	}
	width, err := strconv.Atoi(string(b[0] + 48))
	if err != nil {
		return nil, err
	}

	// Height
	_, err = f.Read(b)
	if err != nil {
		return nil, err
	}
	height, err := strconv.Atoi(string(b[0] + 48))
	if err != nil {
		return nil, err
	}

	size := int(math.Pow(2, float64(width*(height+1))))
	genome := make([]byte, size)

	// Rest of the genome
	_, err = f.Read(genome)
	if err != nil {
		return nil, err
	}

	// Decompress
	genome2 := make([]uint8, size*4)
	for i, b := range genome {
		genome2[i*4] = decodeMove(b >> 6 & 0b00000011)
		genome2[i*4+1] = decodeMove(b >> 4 & 0b00000011)
		genome2[i*4+2] = decodeMove(b >> 2 & 0b00000011)
		genome2[i*4+3] = decodeMove(b & 0b00000011)
	}

	ai := &AI{width: uint(width), height: uint(height), genome: genome2}

	return ai, nil
}

func NewAIFromGenome(width uint, height uint, genome []uint8) *AI {
	return &AI{width: width, height: height, genome: genome}
}

func decodeMove(b byte) uint8 {
	switch b {
	case 0b00000001:
		return 1
	case 0b00000010:
		return 2
	default:
		return 0
	}
}

func (ai *AI) GetMove(g *entity.Grid) uint8 {
	pos := g.GetCellAsteroid()<<ai.width | uint64(g.GetCellShip())
	return ai.genome[pos]
}
