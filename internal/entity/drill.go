package entity

import (
	"fmt"
	"log"
	"math"
	"runtime/debug"
)

const Ground int = 0

type Layer_t int

//Drill data
type Drill struct {
	Name         string
	X, Y, Z      float64
	length       float64
	Layers       []int     //layers' seq id.
	LayerHeights []float64 //layer's bottom height.
	weight       float64
}

func (drill Drill) Print() {
	log.SetFlags(log.Lshortfile)
	//debug.PrintStack()
	log.Printf("name: %s\nPosition:%.2f, %.2f, %.2f\nLength:%.2f\n",
		drill.Name, drill.X, drill.Y, drill.Z, drill.GetLength())
	fmt.Print("Layers: ")
	printInts(drill.Layers)

	fmt.Print("Heights:")
	printFloat64s(drill.LayerHeights)
	fmt.Printf("Weights:%.4f\n\n", drill.weight)
}

func (drill Drill) MakeDrill(name string, x, y, z float64) Drill {
	return Drill{
		Name:         name,
		X:            x,
		Y:            y,
		Z:            z,
		Layers:       []int{Ground},
		LayerHeights: []float64{z},
	}
}
func (drill *Drill) AddLayer(layer int, layerDepthHeight float64) {
	log.SetFlags(log.Lshortfile)
	drill.Layers = append(drill.Layers, layer)
	if layerDepthHeight > drill.LayerHeights[len(drill.LayerHeights)-1] {
		log.Fatal("error")
	}
	drill.LayerHeights = append(drill.LayerHeights, layerDepthHeight)
}

// getter and setter
func (drill *Drill) SetZ(z float64) {
	drill.Z = z
}
func (drill *Drill) SetWeight(weight float64) {
	drill.weight = weight
}
func (drill Drill) GetWeight() float64 {
	if drill.weight <= 0.001 || math.IsInf(drill.weight, 10) || math.IsNaN(drill.weight) {
		log.SetFlags(log.Lshortfile)
		drill.Print()
		debug.PrintStack()
		log.Fatal("error")
	}
	return drill.weight
}
func (drill Drill) GetBottomHeightByLayer(layer int) (height []float64) {
	for idx := 0; idx < len(drill.Layers); idx++ {
		if layer == drill.Layers[idx] {
			height = append(height, drill.LayerHeights[idx])
		}
	}
	return
}
func (drill *Drill) GetLength() (length float64) {
	if drill.length-0 < 10e-5 && drill.Z > drill.LayerHeights[len(drill.LayerHeights)-1] {
		drill.length = drill.Z - drill.LayerHeights[len(drill.LayerHeights)-1]
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
	if drill.GetLength() != drill.Z-drill.GetBottomHeight() {
		return false
	}
	if len(drill.LayerHeights) != len(drill.Layers) {
		return false
	}
	if len(drill.LayerHeights) > 1 {
		for idx := 1; idx < len(drill.LayerHeights); idx++ {
			if drill.LayerHeights[idx]-drill.LayerHeights[idx-1] > 0 {
				return false
			}
		}
	}
	return true
}
func (drill Drill) HasLayer(layer int) (num int) {
	for _, l := range drill.Layers {
		if l == layer {
			num++
		}
	}
	return num
}
func (drill Drill) HasBlock(ceil, floor float64) (has bool) {
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
		layers  []int
		heights []float64
	)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if drill.LayerHeights[len(drill.LayerHeights)-1] != drill.GetBottomHeight() {
		log.Fatal("error: drill.LayerHeights[len(drill.LayerHeights)-1] != drill.GetBottomHeight()")
	}
	if len(drill.LayerHeights) != len(drill.Layers) {
		drill.Print()
		log.Printf("%d, %d\n", len(drill.LayerHeights), len(drill.Layers))
		log.Fatal("error: len(drill.LayerHeights) != len(drill.Layers)")
	}

	layers = append(layers, drill.Layers[0])
	heights = append(heights, drill.LayerHeights[0])

	//37 84 149
	for idx := 1; idx < len(drill.LayerHeights); idx++ {
		if layers[len(layers)-1] == drill.Layers[idx] {
			heights[len(heights)-1] = drill.LayerHeights[idx]
		} else {
			layers = append(layers, drill.Layers[idx])
			heights = append(heights, drill.LayerHeights[idx])
		}
	}
	drill.Layers = layers
	drill.LayerHeights = heights
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
	var drillLayers []int = []int{Ground}

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
	scattered.LayerHeights = drillBlocks
	scattered.Layers = drillLayers
	scattered.SetWeight(drill.GetWeight())
	return
}
func (drill Drill) GetLayerSeq(ceil, floor float64) (seq int, ok bool) {
	// drill top >=ceil >= floor >= drill bottom
	if floor > drill.Z || ceil < drill.LayerHeights[len(drill.LayerHeights)-1] {
		return
	}
	//case1: 1 or less layer in block
	for idx := 1; idx < len(drill.LayerHeights); idx++ {
		if drill.LayerHeights[idx] <= floor &&
			drill.LayerHeights[idx-1] >= ceil && idx < len(drill.Layers) {
			return drill.Layers[idx], true
		}
	}

	//case2: 2 layers in block
	if ceil <= drill.Z && floor >= drill.LayerHeights[len(drill.LayerHeights)-1] {
		//here suppose that resolution z < min layer thick,
		//so there are 2 layers in the block at most.
		var bidx []int
		var thick []float64

		//layer surface in block.
		for idx, h := range drill.LayerHeights {
			if h < ceil && h > floor {
				bidx = append(bidx, idx)
			}
		}

		if len(bidx) < 1 {
			return -1, false
		}

		l := len(bidx)
		thick = append(thick, ceil-drill.LayerHeights[bidx[0]])
		for idx := 1; idx < l; idx++ {
			thick = append(thick,
				drill.LayerHeights[bidx[idx]]-drill.LayerHeights[bidx[idx-1]])
		}

		//!!
		bidx = append(bidx, bidx[l-1]+1)
		thick = append(thick, drill.LayerHeights[bidx[l-1]]-floor)
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
	if ceil > drill.GetBottomHeight() && floor < drill.GetBottomHeight() {
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
func printInts(s []int) () {
	fmt.Print("[")
	for _, v := range s {
		fmt.Printf("%4d\t", v)
	}
	fmt.Print("]\n")
}
func (drill Drill) HasRepeatLayers() bool {
	seq := drill.Layers
	layerMap := make(map[int]int)
	for _, l := range seq {
		if _, ok := layerMap[l]; ok {
			return true
		} else {
			layerMap[l] = 1
		}
	}
	return false
}
func (drill Drill) StdSeq(stdSeq []int) Drill {
	var (
		seq []int     = []int{0}
		h   []float64 = []float64{drill.Z}
	)
	if !drill.IsValid() {
		log.SetFlags(log.Lshortfile)
		drill.Print()
		debug.PrintStack()
		log.Fatal("error")
	}

	layers := drill.Layers
	LayerFloorHeights := drill.LayerHeights
	var idx1, idx2 int = 1, 1
	for idx2 < len(stdSeq) && idx1 < len(layers) {
		if layers[idx1] == stdSeq[idx2] {
			h = append(h, LayerFloorHeights[idx1])
			seq = append(seq, layers[idx1])
			idx1++
			idx2++
		} else {
			h = append(h, LayerFloorHeights[idx1-1])
			seq = append(seq, stdSeq[idx2])
			idx2++
		}
	}
	for ; idx2 < len(stdSeq); idx2++ {
		h = append(h, LayerFloorHeights[len(LayerFloorHeights)-1])
		seq = append(seq, stdSeq[idx2])
	}
	drill.Layers = seq
	drill.LayerHeights = h
	return drill
}
func (drill Drill) UnStdSeq() Drill {
	var (
		seq []int     = []int{0}
		h   []float64 = []float64{drill.Z}
	)
	if len(drill.Layers) < 2 {
		log.SetFlags(log.Lshortfile)
		log.Fatal("error")
	}
	for idx := 1; idx < len(drill.LayerHeights); idx++ {
		if drill.LayerHeights[idx] == drill.LayerHeights[idx-1] {
			idx++
			continue
		} else {
			seq = append(seq, drill.Layers[idx])
			h = append(h, drill.LayerHeights[idx])
			idx++
		}
	}
	drill.Layers = seq
	drill.LayerHeights = h
	return drill
}
func (drill Drill) Round() (drill2 Drill) {
	drill2 = drill
	drill2.X = math.Ceil(drill.X)
	drill2.Y = math.Ceil(drill.Y)
	drill2.Z = math.Ceil(drill.Z)
	for idx, h := range drill2.LayerHeights {
		drill2.LayerHeights[idx] = math.Ceil(h)
	}
	return drill2
}
func (drill Drill) RoundDrills(drills []Drill) *[]Drill {
	var rDrills []Drill
	for _, d := range drills {
		rDrills = append(rDrills, d.Round())
	}
	return &rDrills
}
