package main

import (
	"cmp"
	"fmt"
	"log"
	"math"
	"os"
	"pgregory.net/rand"
	"runtime"
	"time"

	"github.com/ruloweb/asteroids/entity"
	"github.com/ruloweb/asteroids/system"
	"golang.org/x/exp/slices"

	"github.com/tomcraven/goga"
)

// TODO: this file is a mess, refactor it

type TrainCmd struct {
	NumThreads     int    `default:"8" help:"Number of threads."`
	NumIter        int    `default:"100" help:"Number of iterations to run."`
	PopulationSize int    `default:"500" help:"Number of executions per iteration."`
	MaxAsteroids   uint   `default:"50" help:"Max number of asteroids at the same time."`
	TicksLimit     uint64 `default:"65535" help:"Max number of ticks to play."`
	GridWidth      uint   `default:"4" help:"Number of columns in the grid."`
	GridHeight     uint   `default:"2" help:"Number of rows in the grid."`
	OutputGen      string `default:"output.gen" help:"Output genome file."`
	OutputResults  string `default:"output.txt" help:"Output training results file."`
}

func runGame(genome []uint8, gridWidth uint, gridHeight uint, ticksLimit uint64, maxAsteroids uint) uint64 {
	// Create ship
	ship := entity.NewShip(160, 440)

	// Genome
	var ai *system.AI
	var err error
	ai = system.NewAIFromGenome(gridWidth, gridHeight, genome)
	if err != nil {
		log.Fatal(err)
	}

	grid := entity.NewGrid(gridWidth, gridHeight, system.ScreenW, system.ScreenH)
	g := system.NewGame(true, false, ticksLimit, maxAsteroids, ship, grid, ai)

	for i := 0; i < int(ticksLimit); i++ {
		g.Update()
	}

	return g.AsteroidsCrashed
}

type FitnessChan struct {
	iter    int
	fitness uint64
}

type stringMaterSimulator struct {
	c           *TrainCmd
	fitness     []FitnessChan
	fitnessChan *chan FitnessChan
	size        int
}

func (sms *stringMaterSimulator) OnBeginSimulation() {
}
func (sms *stringMaterSimulator) OnEndSimulation() {
}
func (sms *stringMaterSimulator) Simulate(g goga.Genome) {
	buffer := make([]uint8, sms.size)
	bits := g.GetBits()
	for i := 0; i < bits.GetSize(); i++ {
		buffer[i] = uint8(bits.Get(i))
	}

	fitness := runGame(buffer, sms.c.GridWidth, sms.c.GridHeight, sms.c.TicksLimit, sms.c.MaxAsteroids)
	*sms.fitnessChan <- FitnessChan{iter: currentIter, fitness: fitness}

	g.SetFitness(int(fitness))
}
func (sms *stringMaterSimulator) ExitFunc(g goga.Genome) bool {
	if currentIter >= sms.c.NumIter {
		genome = g
		return true
	}
	return false
}

type myBitsetCreate struct {
	size int
}

func (bc *myBitsetCreate) Go() goga.Bitset {
	b := goga.Bitset{}
	b.Create(bc.size)
	for i := 0; i < bc.size; i++ {
		b.Set(i, rand.Intn(3))
	}
	return b
}

type myEliteConsumer struct {
	sms *stringMaterSimulator
}

func (ec *myEliteConsumer) OnElite(g goga.Genome) {
	var acc uint64
	var count float64

	// Close fitness channel
	close(*ec.sms.fitnessChan)

	for f := range *ec.sms.fitnessChan {
		acc += f.fitness
		count++
		ec.sms.fitness = append(ec.sms.fitness, f)
	}

	// Open the fitness channel again
	*ec.sms.fitnessChan = make(chan FitnessChan, ec.sms.c.PopulationSize)

	currentIter++

	log.Printf("Iter %d, avg fitness: %f\n", currentIter, float64(acc)/count)
}

func genomeCmp(a, b goga.Genome) int {
	return cmp.Compare(a.GetFitness(), b.GetFitness())
}

func Roulette(genomeArray []goga.Genome, totalFitness int) goga.Genome {
	if len(genomeArray) == 0 {
		panic("genome array contains no elements")
	}

	slices.SortFunc(genomeArray, genomeCmp)
	elite := int(math.Ceil(float64(len(genomeArray)) * 0.2))

	r := genomeArray[rand.Intn(elite)]
	return r
}

func Mutate(g1, g2 goga.Genome) (goga.Genome, goga.Genome) {
	g1Bits := g1.GetBits().CreateCopy()
	for i := 0; i < 1; i++ {
		randomBit := rand.Intn(g1Bits.GetSize())
		g1Bits.Set(randomBit, rand.Intn(3))
	}

	g2Bits := g2.GetBits().CreateCopy()
	for i := 0; i < 1; i++ {
		randomBit := rand.Intn(g2Bits.GetSize())
		g2Bits.Set(randomBit, rand.Intn(3))
	}
	return goga.NewGenome(g1Bits), goga.NewGenome(g2Bits)
}

var (
	genome      goga.Genome
	currentIter int
	fitnessChan chan FitnessChan
)

func (c *TrainCmd) Run() error {
	runtime.GOMAXPROCS(c.NumThreads)

	size := int(math.Pow(2, float64(c.GridWidth*(c.GridHeight+1))))
	fitnessChan = make(chan FitnessChan, c.PopulationSize)

	genAlgo := goga.NewGeneticAlgorithm()

	simulator := stringMaterSimulator{c: c, fitnessChan: &fitnessChan, size: size}
	genAlgo.Simulator = &simulator
	genAlgo.BitsetCreate = &myBitsetCreate{size: size}
	genAlgo.EliteConsumer = &myEliteConsumer{sms: &simulator}
	genAlgo.Mater = goga.NewMater(
		[]goga.MaterFunctionProbability{
			//{P: 1.0, F: goga.TwoPointCrossover},
			{P: 1.0, F: goga.OnePointCrossover},
			{P: 1.0, F: Mutate},
			//{P: 1.0, F: goga.UniformCrossover, UseElite: true},
		},
	)
	genAlgo.Selector = goga.NewSelector(
		[]goga.SelectorFunctionProbability{
			{P: 1.0, F: Roulette},
		},
	)

	genAlgo.Init(c.PopulationSize, c.NumThreads)

	startTime := time.Now()
	genAlgo.Simulate()
	fmt.Println(time.Since(startTime))

	// Save fitness
	file, err := os.Create(c.OutputResults)
	defer file.Close()
	_, err = file.WriteString("iter,fitness\n")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range simulator.fitness {
		line := fmt.Sprintf("%d,%d\n", f.iter, f.fitness)
		_, err := file.WriteString(line)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Save final genome
	buffer := make([]byte, size/4+2)
	buffer[0] = byte(c.GridWidth)
	buffer[1] = byte(c.GridHeight)
	bits := genome.GetBits()
	for i := 0; i < bits.GetSize(); i += 4 {
		for j := 0; j < 4; j++ {
			b := bits.Get(i + j)
			buffer[2+i/4] |= (byte(b) & 0xFF) << (6 - j*2)
		}
	}
	err = os.WriteFile(c.OutputGen, buffer, 0644)
	if err != nil {
		log.Println(err.Error())
	}

	return nil
}
