package entity

import (
	"fmt"
	"log"
	"math"
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
func (drill Drill) hasLayer(layer int64) (num int) {
	for _, l := range drill.Layers {
		if l == layer {
			num++
		}
	}
	return num
}
func (drill Drill) hasBlock(ceil, floor float64) (has bool) {
	if ceil <= floor {
		return false
	}
	drillCeil, drillFloor := drill.Z, drill.GetBottomHeight()
	//已经规定block范围小于最小层厚
	if ceil <= drillCeil && floor >= drillFloor ||
		ceil > drillCeil && floor < drillCeil ||
		ceil > drillFloor && floor < drillFloor {
		has = true
		return
	}
	return false
}
func (drill *Drill) Merge() {
	var (
		layers  []int64
		heights []float64
	)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if drill.LayerFloorHeights[len(drill.LayerFloorHeights)-1] != drill.GetBottomHeight() {
		log.Fatal("error: drill.LayerFloorHeights[len(drill.LayerFloorHeights)-1] != drill.GetBottomHeight()")
	}
	if len(drill.LayerFloorHeights) != len(drill.Layers) {
		drill.Print()
		log.Printf("%d, %d\n", len(drill.LayerFloorHeights), len(drill.Layers))
		log.Fatal("error: len(drill.LayerFloorHeights) != len(drill.Layers)")
	}

	layers = append(layers, drill.Layers[0])
	heights = append(heights, drill.LayerFloorHeights[0])

	//37 84 149
	for idx := 1; idx < len(drill.LayerFloorHeights); idx++ {
		if layers[len(layers)-1] == drill.Layers[idx] {
			heights[len(heights)-1] = drill.LayerFloorHeights[idx]
		} else {
			layers = append(layers, drill.Layers[idx])
			heights = append(heights, drill.LayerFloorHeights[idx])
		}
	}
	drill.Layers = layers
	drill.LayerFloorHeights = heights
}
func (drill Drill) DistanceBetween(drill2 Drill) (dist float64) {
	x1, y1, x2, y2 := drill.X, drill.Y, drill2.X, drill2.Y
	dist = math.Sqrt((x2-x1)*(x2-x1) + (y2-y1)*(y2-y1))
	if math.IsNaN(dist) || math.IsInf(dist, 0) || dist < 0 {
		return -1
	}
	return dist
}
func (drill Drill) Explode(blocks []float64) (scattered Drill) {
	if blocks == nil {
		return
	}
	if blocks[0] < drill.Z {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Println("abnormal Drills")
		drill.Print()
	}
	scattered = scattered.MakeDrill(drill.Name, drill.X, drill.Y, drill.Z)

	var idxa int
	for idx, h := range blocks {
		if h < drill.Z {
			idxa = idx
			break
		}
	}
	var drillBlocks []float64 = []float64{drill.Z}
	var drillLayers []int64 = []int64{Ground}

	for idx := idxa; idx < len(blocks); idx++ {
		if blocks[idx] <= drill.Z && blocks[idx] >= drill.Z-drill.GetLength() {
			drillBlocks = append(drillBlocks, blocks[idx])
			if seq, ok := drill.GetLayerSeq(
				drillBlocks[len(drillBlocks)-2], drillBlocks[len(drillBlocks)-1]); ok {
				drillLayers = append(drillLayers, seq)
			}
		}
	}
	if drillBlocks[len(drillBlocks)-1] != drill.Z-drill.GetLength() {
		drillBlocks = append(drillBlocks, drill.Z-drill.GetLength())
		if seq, ok := drill.GetLayerSeq(
			drillBlocks[len(drillBlocks)-2], drillBlocks[len(drillBlocks)-1]); ok {
			drillLayers = append(drillLayers, seq)
		}
	}
	scattered.LayerFloorHeights = drillBlocks
	scattered.Layers = drillLayers
	scattered.SetWeight(drill.GetWeight())
	return
}
func (drill Drill) GetLayerSeq(ceil, floor float64) (seq int64, ok bool) {
	// drill top >=ceil >= floor >= drill bottom
	if floor > drill.Z || ceil < drill.LayerFloorHeights[len(drill.LayerFloorHeights)-1] {
		return
	}
	//case1: 1 or less layer in block
	for idx := 1; idx < len(drill.LayerFloorHeights); idx++ {
		if drill.LayerFloorHeights[idx] <= floor &&
			drill.LayerFloorHeights[idx-1] >= ceil && idx < len(drill.Layers) {
			return drill.Layers[idx], true
		}
	}

	//case2: 2 layers in block
	if ceil <= drill.Z && floor >= drill.LayerFloorHeights[len(drill.LayerFloorHeights)-1] {
		//here suppose that resolution z < min layer thick,
		//so there are 2 layers in the block at most.
		var bidx []int
		var thick []float64

		//layer surface in block.
		for idx, h := range drill.LayerFloorHeights {
			if h < ceil && h > floor {
				bidx = append(bidx, idx)
			}
		}

		if len(bidx) < 1 {
			return -1, false
		}

		l := len(bidx)
		thick = append(thick, ceil-drill.LayerFloorHeights[bidx[0]])
		for idx := 1; idx < l; idx++ {
			thick = append(thick,
				drill.LayerFloorHeights[bidx[idx]]-drill.LayerFloorHeights[bidx[idx-1]])
		}

		//!!
		bidx = append(bidx, bidx[l-1]+1)
		thick = append(thick, drill.LayerFloorHeights[bidx[l-1]]-floor)
		if len(bidx) > 2 {
			log.SetFlags(log.Lshortfile | log.LstdFlags)
			log.Println("Warning, the resolution z is too large!")
			log.Printf("param: ceil %.2f, floor %.2f, block %.2f", ceil, floor, ceil-floor)
			log.Println(drill)
		}

		var maxThick float64 = -math.MaxFloat64
		var maxIndex int = 0
		for idx, thick := range thick {
			if math.Abs(thick) > maxThick {
				maxThick = math.Abs(thick)
				maxIndex = bidx[idx]
			}
		}
		if maxIndex < len(drill.Layers) {
			return drill.Layers[maxIndex], true
		}
	}

	//case3.1: boundary
	if ceil > drill.Z && floor < drill.Z {
		return drill.GetLayerSeq(drill.Z, floor)
	}

	//case 3.2
	if ceil > drill.GetBottomHeight() && floor < drill.GetWeight() {
		return drill.GetLayerSeq(ceil, drill.GetBottomHeight())
	}
	return -1, false
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
