package pdga

import (
	"fmt"
	"math/rand"
	"sort"
)

// XXX reimplement as a container/heap, use goroutine per individual,
// let indivs pick mates, continuous evolution instead of generations,
// add support for path dependence and event stream

type Individual interface {
	GetDNA() []byte
	PutDNA([]byte)
	Breed(Individual) Individual
	Fitness(World) float64
	IsPerfect() bool
	String() string
}

type World interface {
	GetPopulation() []Individual
	String() string
}

func Evolve(w World, generations int, verbose bool) {
	pop := w.GetPopulation()
	popsize := len(pop)

	// run generations
	for gen := 0; gen < generations; gen++ {

		// sort
		sort.Slice(pop, func(i, j int) bool {
			ind := pop[i]
			indj := pop[j]
			return ind.Fitness(w) > indj.Fitness(w)
		})

		if verbose {
			fmt.Printf("%8d %s\n", gen, pop[0])
		}

		if pop[0].IsPerfect() {
			break
		}

		// breed
		halfsize := popsize / 2
		wc := WeightedChoice{world: w, cutoff: halfsize}
		wc.Init()
		for id := halfsize; id < popsize; id++ {
			a := wc.Choose()
			b := wc.Choose()
			c := a.Breed(b)
			// fmt.Println(c)
			pop[id] = c
		}
	}

	return
}

// WeightedChoice uses weighted random selection to return one of the supplied
// choices.  Weights of 0 are never selected.  All other weight values are
// relative.  E.g. if you have two choices both weighted 3, they will be
// returned equally often; and each will be returned 3 times as often as a
// choice weighted 1.
//
// Derived from:
// https://github.com/jmcvetta/randutil/blob/2bb1b664bcff/randutil.go#L121
// Based on this algorithm:
// http://eli.thegreenplace.net/2010/01/22/weighted-random-generation-in-python/
type WeightedChoice struct {
	world  World
	pop    []Individual
	cutoff int
	sum    float64
}

func (wc *WeightedChoice) Init() {
	wc.pop = wc.world.GetPopulation()
	for id := 0; id < wc.cutoff; id++ {
		ind := wc.pop[id]
		wc.sum += ind.Fitness(wc.world)
	}
}

func (wc *WeightedChoice) Choose() (ind Individual) {
	size := len(wc.pop)
	r := float64(rand.Intn(size))
	for id := 0; id < size; id++ {
		ind = wc.pop[id]
		r -= ind.Fitness(wc.world)
		if r < 0 {
			return
		}
	}
	return
}
