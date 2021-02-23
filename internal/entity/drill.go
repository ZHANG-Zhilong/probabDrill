package entity

import (
	"fmt"
	"log"
)

const Ground int64 = 0

type Layer_t int64

//Drill data
type Drill struct {
	Name              string
	X, Y, Z           float64
	length            float64
	Layers            []int64   //layers' seq id.
	LayerFloorHeights []float64 //layer's bottom height.
	weight            float64
}

func (drill Drill) Print() {
	log.SetFlags(log.Lshortfile)
	log.Printf("name: %s\nPosition:%.2f, %.2f, %.2f\nLength:%.2f\n",
		drill.Name, drill.X, drill.Y, drill.Z, drill.GetLength())
	fmt.Print("Layers: ")
	printInt64s(drill.Layers)

	fmt.Print("Heights:")
	printFloat64s(drill.LayerFloorHeights)
	fmt.Printf("Weights:%.4f\n\n", drill.weight)
}

func (drill Drill) MakeDrill(name string, x, y, z float64) Drill {
	return Drill{
		Name:              name,
		X:                 x,
		Y:                 y,
		Z:                 z,
		Layers:            []int64{Ground},
		LayerFloorHeights: []float64{z},
	}
}

// getter and setter
func (drill *Drill) SetZ(z float64) {
	drill.Z = z
}
func (drill *Drill) SetWeight(weight float64) {
	drill.weight = weight
}
func (drill Drill) GetWeight() float64 {
	return drill.weight
}

func (drill *Drill) GetLength() (length float64) {
	if drill.length == 0.0 {
		drill.length = drill.Z - drill.LayerFloorHeights[len(drill.LayerFloorHeights)-1]
	}
	return drill.length
}
func (drill *Drill) SetLength(length float64) {
	drill.length = length
}
func (drill Drill) GetBottomHeight() (bottom float64) {
	bottom = drill.Z - drill.GetLength()
	return
}
func (drill Drill) IsValid() (valid bool) {
	if drill.length != drill.Z-drill.GetBottomHeight() {
		return false
	}
	if len(drill.LayerFloorHeights) != len(drill.Layers) {
		return false
	}
	if len(drill.LayerFloorHeights) > 1 {
		for idx := 1; idx < len(drill.LayerFloorHeights); idx++ {
			if drill.LayerFloorHeights[idx]-drill.LayerFloorHeights[idx-1] > 0 {
				return false
			}
		}
	}
	return true
}
func printFloat64s(s []float64) () {
	fmt.Print("[")
	for _, v := range s {
		fmt.Printf("%+2.2f\t", v)
	}
	fmt.Print("]\n")
}
func printInt64s(s []int64) () {
	fmt.Print("[")
	for _, v := range s {
		fmt.Printf("%4d\t", v)
	}
	fmt.Print("]\n")
}
