package main

//@danielmatthewsgrout

//A very basic Self Organising Map implementation with limited parallelism

import (
	"math"
	"math/rand"
	"sync"
)

//DistanceFunction a function to measure distance between 2 vectors of same len
type DistanceFunction func(a1, a2 []float64) float64

//InitFunction - the function used to initialise the SOM
type InitFunction func(dimensions, xySize int, vectors [][]float64) [][]float64

//tweak this if needed - I found this was the best value on a 4c/8t Intel
const concurrent = 32

//BasicSOM a basicSom
type BasicSOM struct {
	nodes  [][]float64
	xySize int
	df     DistanceFunction
}

//TrainBasicSOM ronseal
func TrainBasicSOM(dimensions, xySize, maxSteps int, learningDecayRate float64, vectors [][]float64) SelfOrgMap {

	initf := RandomInit

	df := EuclideanDistanceSquared

	nodes := initf(dimensions, xySize, vectors)

	println("SOM processing...")

	fXYSize := float64(xySize)
	fMaxSteps := float64(maxSteps)
	v := 0
	for i := 0; i < maxSteps; i++ {
		fI := float64(i) + 1

		selection := vectors[v] //select vector for this run

		v++
		if v == len(vectors) {
			v = 0
		}

		winner := getNearest(selection, nodes, df)                            //find the nearest node to this selection
		radius := fXYSize * math.Exp(-(fI / (fMaxSteps / math.Log(fXYSize)))) //calculate the radius for neighbour sarch
		lRate := learningDecayRate * math.Exp(-fI/fMaxSteps)                  //calculate the learning rate for this run

		sq := (radius * radius)
		x2 := winner % xySize
		y2 := winner / xySize

		i := 0
		wg := sync.WaitGroup{}

		k := len(nodes) / concurrent
		//find the smallest
		for i < len(nodes) {

			if i+k >= len(nodes) {
				k = len(nodes) - i
			}
			wg.Add(1)
			go func(i, k int) { //dividde and conquer the area that needs updating

				for z := i; z < i+k; z++ {
					n := nodes[z]
					x1 := z % xySize
					y1 := z / xySize

					if dist := float64((x1-x2)*(x1-x2) + (y1-y2)*(y1-y2)); dist < sq { //caclulate the distance from our node to this node
						neighbour := math.Exp(-dist / (2 * sq)) //neighbour modification value - decays as it gets further away

						for v, s := range selection {
							n[v] += lRate * neighbour * (s - n[v])
						}

					}
				}
				wg.Done()
			}(i, k)
			i += k
		}
		wg.Wait()
	}

	println("SOM done.")

	return &BasicSOM{
		nodes:  nodes,
		xySize: xySize,
		df:     df,
	}
}

//GetBMU get best matching unit x and y coords
func (s *BasicSOM) GetBMU(vec []float64) (int, int) {
	winner := getNearest(vec, s.nodes, s.df)
	return (winner % s.xySize) + 1, (winner / s.xySize) + 1
}

func getNearest(vec []float64, nodes [][]float64, df DistanceFunction) int {

	var minDistance float64 = math.MaxFloat64
	var winner int
	i := 0
	wg := sync.WaitGroup{}
	var lck sync.Mutex

	k := len(nodes) / concurrent
	//find the smallest
	for i < len(nodes) {

		if i+k >= len(nodes) {
			k = len(nodes) - i
		}

		wg.Add(1)
		go func(i, k int) { //divide and conquer the seach
			var md float64 = math.MaxFloat64
			var w int

			for x := i; x < i+k; x++ {
				if d := df(vec, nodes[x]); d < md {
					md = d
					w = x
				}
			}
			lck.Lock()
			if md < minDistance {
				minDistance = md
				winner = w
			}
			lck.Unlock()
			wg.Done()
		}(i, k)
		i += k
	}

	wg.Wait()

	return winner
}

//EuclideanDistanceSquared squared Euclidean distance function
func EuclideanDistanceSquared(a1, a2 []float64) float64 {
	var d float64

	for i := range a1 {
		t := a1[i] - a2[i]
		d += t * t
	}

	return d
}

//RandomInit - init with 0 to 1 random values
func RandomInit(dimensions, xySize int, vectors [][]float64) [][]float64 {
	//init with random values between 0,0 and 1.0
	nodes := make([][]float64, xySize*xySize)

	println("initialising map")
	i := 0
	for x := 0; x < xySize; x++ {
		for y := 0; y < xySize; y++ {
			weight := make([]float64, dimensions)
			for w := range weight {
				weight[w] = rand.Float64()
			}
			nodes[i] = weight
			i++
		}
	}
	return nodes
}
