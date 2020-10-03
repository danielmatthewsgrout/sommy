package main

//SelfOrgMap - a self organising map
type SelfOrgMap interface {
	GetBMU(vec []float64) (int, int)
}
