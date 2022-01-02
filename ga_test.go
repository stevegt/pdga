package pdga

import (
	"fmt"
	"math/rand"
	"testing"
)

var ref = []byte("monkeyswithtypewriterscantypeallbooksintheuniverseeventuallygivenendlesstime")

type individual struct {
	dna     []byte
	fitness float64
}

func mkindividual() (c *individual) {
	c = &individual{}
	c = c.copy()
	return
}

func (a *individual) copy() (c *individual) {
	c = &individual{}
	c.dna = make([]byte, len(ref))
	c.fitness = -1
	return
}

func (a *individual) Breed(b Individual) (c Individual) {
	c = a.copy()
	adna := a.GetDNA()
	bdna := b.GetDNA()
	cdna := c.GetDNA()
	pos := rand.Intn(len(cdna))
	// crossover
	copy(cdna[:pos], adna[:pos])
	copy(cdna[pos:], bdna[pos:])
	// mutate at crossover point
	cdna[pos] = byte(rand.Intn(256))
	// fmt.Println(c)
	c.PutDNA(cdna)
	// fmt.Println(c)
	return
}

func (ind *individual) GetDNA() (dna []byte) {
	dna = ind.dna[:]
	return
}

func (ind *individual) PutDNA(dna []byte) {
	copy(ind.dna[:], dna)
	ind.fitness = -1
	return
}

func (ind *individual) Fitness(w World) (score float64) {
	if ind.fitness < 0 {
		score = 0
		for pos := 0; pos < len(ind.dna); pos++ {
			if ind.dna[pos] == ref[pos] {
				score++
			}
		}
		ind.fitness = score
	}
	return ind.fitness
}

func (ind *individual) IsPerfect() bool {
	return ind.fitness == float64(len(ref))
}

func (ind *individual) String() (out string) {
	var outbytes []byte
	for i := 0; i < len(ind.dna); i++ {
		g := ind.dna[i]
		if g < 33 || g > 126 {
			outbytes = append(outbytes, ' ')
		} else {
			outbytes = append(outbytes, g)
		}
	}
	out = fmt.Sprintf("%s %6.2f", string(outbytes), ind.fitness)
	return
}

type world struct {
	Population []Individual
	PopSize    int
	Proto      Individual
}

func (w *world) Init() {
	// initial population
	w.Population = make([]Individual, w.PopSize)
	for id := 0; id < w.PopSize; id++ {
		// new individual
		c := mkindividual()
		// randomize dna
		cdna := c.GetDNA()
		for pos := 0; pos < len(cdna); pos++ {
			cdna[pos] = byte(rand.Intn(256))
		}
		c.PutDNA(cdna)
		w.Population[id] = c
	}
	return
}

func (w world) GetPopulation() []Individual {
	return w.Population
}

func (w world) String() string {
	out := ""
	for _, x := range w.Population {
		out += fmt.Sprintf("%s\n", x)
	}
	return out
}

func TestSimple(t *testing.T) {
	// defer profile.Start().Stop()
	ind := individual{}
	w := world{PopSize: 1000, Proto: &ind}
	w.Init()
	Evolve(w, 10000, true)
	fmt.Println()
	// fmt.Println("        ", pop[len(pop)-1])
	fmt.Println("------------")
	// fmt.Println(w)
	return
}
