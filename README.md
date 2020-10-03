# Sommy - Self Organising Map in Golang

## Features

* Light weight and very simple implementation
* Uses "divide and conquer" parallelism for BMU search and neighbour updates

## Limitations

* Only supports square grids - x==y
* Currently only supports Gausian Neighbour, and Exponential Radius and Learning Decay
* Best match is done using Squared Euclidean Distance

## How to Use

### Generating a Map from Vectors

```go
TrainBasicSOM(dimensions, xySize, maxSteps int, learningDecayRate float64, vectors [][]float64) SelfOrgMap
```

* __dimensions__ - the dimensionality of the vectors

* __xySize__ - the dimensions of the grid

* __maxSteps__ - the number of iterations to perform

* __learningDecayRate__ - the variable fed into the exponential decay - controls the convergence onto the fit

* __vectors__ - two dimensional array of the vectors to map

### Using the Generated Map

It returns an interace to SelfOrgMap which provides the following function:

```go
GetBMU(vec []float64) (int, int)
```

Pass a vector to this and it will return the x,y coordinates of the best matching unit on the grid.

---

## Details

__Author__ @danielmatthewsgrout <dan@dmg.dev>  
__Licence__ MIT  
__Current Version__ 1.0.0  
