package entity

import (
	"fmt"
	"log"
	"math"
	"probabDrill"
	"runtime/debug"
	"sort"
	"strconv"
)

const Ground int = 0

//Drill data
type Drill struct {
	Name         string
	X, Y, Z      float64
	length       float64
	Layers       []int     //layers' seq id.
	LayerHeights []float64 //layer's bottom height.
	weight       float64
}

func NewBasicDrill(name string, x, y, z float64) *Drill {
	return &Drill{
		Name: name,
		X:    x,
		Y:    y,
		Z:    z,
	}
}

func (drill Drill) Display() {
	log.SetFlags(log.Lshortfile)
	//debug.PrintStack()
	log.Printf("Name: %s\nPosition:%.2f, %.2f, %.2f\nLength:%.2f\n",
		drill.Name, drill.X, drill.Y, drill.Z, drill.GetLength())
	fmt.Print("Layers: ")
	printSliceInt(drill.Layers)

	fmt.Print("Heights:")
	printSliceFloat64(drill.LayerHeights)
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
func (drill *Drill) AddLayer(layer int) (err error) {
	log.SetFlags(log.Lshortfile)
	if len(drill.LayerHeights) > len(drill.Layers) {
		drill.Layers = append(drill.Layers, layer)
		return nil
	}
	return fmt.Errorf(":add failed")
}
func (drill *Drill) AddLayerWithHeight(layer int, layerDepthHeight float64) {
	log.SetFlags(log.Lshortfile)
	drill.Layers = append(drill.Layers, layer)
	if layerDepthHeight > drill.LayerHeights[len(drill.LayerHeights)-1] {
		log.Fatal("error")
	}
	drill.LayerHeights = append(drill.LayerHeights, layerDepthHeight)
}
func (drill *Drill) SetZ(z float64) {
	drill.Z = z
}
func (drill *Drill) SetWeight(weight float64) {
	drill.weight = decimal(weight)
}
func (drill Drill) GetWeight() float64 {
	log.SetFlags(log.Lshortfile)
	if math.IsInf(drill.weight, 10) || math.IsNaN(drill.weight) {
		debug.PrintStack()
		drill.Display()
		log.Fatal("invalid drill weight.\n")
	}
	return drill.weight
}
func (drill Drill) GetLayerBottomHeight(layer int) (height []float64) {
	for idx := 0; idx < len(drill.Layers); idx++ {
		if layer == drill.Layers[idx] {
			height = append(height, drill.LayerHeights[idx])
		}
	}
	return
}
func (drill *Drill) GetLength() (length float64) {
	if math.Abs(drill.length-0) < 1e-7 && drill.Z > drill.BottomHeight() {
		drill.length = math.Abs(drill.Z - drill.LayerHeights[len(drill.LayerHeights)-1])
	}
	return drill.length
}
func (drill *Drill) SetLength(length float64) {
	drill.length = length
}
func (drill Drill) BottomHeight() (bottom float64) {
	//bottom = drill.Z - drill.GetLength()
	bottom = drill.LayerHeights[len(drill.LayerHeights)-1]
	return
}
func (drill Drill) IsValid() (valid bool) {
	log.SetFlags(log.Lshortfile)
	if math.Abs(drill.GetLength()-drill.Z-drill.BottomHeight()) < 1e-2 {
		log.Fatal("math.Abs(drill.GetLength() - drill.Z-drill.BottomHeight() ) < 1e-2")
	}
	if len(drill.LayerHeights) != len(drill.Layers) {
		log.Fatal("len(drill.LayerHeights) != len(drill.Layers)")
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
func (drill Drill) LayerThickness(layer int) (thickness float64, ok bool) {
	//only return first layer's thickness
	if drill.HasLayer(layer) > 0 {
		for idx, _ := range drill.LayerHeights {
			if layer == drill.Layers[idx] && idx >= 1 {
				return drill.LayerHeights[idx-1] - drill.LayerHeights[idx], true
			}
		}
		return
	} else {
		return -1, false
	}
}
func (drill Drill) HasBlock(ceil, floor float64) (has bool) {
	if ceil <= floor {
		return false
	}
	drillCeil, drillFloor := drill.Z, drill.BottomHeight()
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

	if drill.LayerHeights[len(drill.LayerHeights)-1] != drill.BottomHeight() {
		log.Fatal("error: drill.LayerHeights[len(drill.LayerHeights)-1] != drill.BottomHeight()")
	}
	if len(drill.LayerHeights) != len(drill.Layers) {
		drill.Display()
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
func (drill Drill) Distance(drill2 Drill) (dist float64) {
	return math.Hypot(drill.X-drill2.X, drill.Y-drill2.Y)
}
func (drill Drill) Explode(blocks []float64) (scattered Drill) {
	if blocks == nil {
		return
	}
	if blocks[0] < drill.Z {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Println("abnormal Drills")
		drill.Display()
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
		//here suppose that resolution z < min layer thickness,
		//so there are 2 layers in the block at most.
		var bidx []int
		var thickness []float64

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
		thickness = append(thickness, ceil-drill.LayerHeights[bidx[0]])
		for idx := 1; idx < l; idx++ {
			thickness = append(thickness,
				drill.LayerHeights[bidx[idx]]-drill.LayerHeights[bidx[idx-1]])
		}

		//!!
		bidx = append(bidx, bidx[l-1]+1)
		thickness = append(thickness, drill.LayerHeights[bidx[l-1]]-floor)
		if len(bidx) > 2 {
			log.SetFlags(log.Lshortfile | log.LstdFlags)
			log.Println("Warning, the resolution z is too large!")
			log.Printf("param: ceil %.2f, floor %.2f, block %.2f", ceil, floor, ceil-floor)
			log.Println(drill)
		}

		var maxThick float64 = -math.MaxFloat64
		var maxIndex int = 0
		for idx, thick := range thickness {
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
	if ceil > drill.BottomHeight() && floor < drill.BottomHeight() {
		return drill.GetLayerSeq(ceil, drill.BottomHeight())
	}
	return -1, false
}
func printSliceFloat64(s []float64) () {
	fmt.Print("[")
	for _, v := range s {
		fmt.Printf("%+2.2f\t", v)
	}
	fmt.Print("]\n")
}
func printSliceInt(s []int) () {
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
func (drill Drill) UnifySeq(stdSeq []int) Drill {
	var (
		seq []int     = []int{0}
		h   []float64 = []float64{drill.Z}
	)
	if !drill.IsValid() {
		log.SetFlags(log.Lshortfile)
		drill.Display()
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
func (drill *Drill) UnBlock() {
	log.SetFlags(log.Lshortfile)
	if len(drill.Layers) < 2 {
		log.Fatal("error")
	}

	layers := []int{0}
	h := []float64{drill.Z}

	for idx := 1; idx < len(drill.LayerHeights); idx++ {
		if math.Abs(drill.LayerHeights[idx]-drill.LayerHeights[idx-1]) < probabDrill.MinThicknessInDrill {
			continue
		} else {
			layers = append(layers, drill.Layers[idx])
			h = append(h, decimal(drill.LayerHeights[idx]))
		}
	}
	drill.Layers = layers
	drill.LayerHeights = h
}
func (drill *Drill) Update() {
	drill.Z = drill.LayerHeights[0]
	drill.length = drill.LayerHeights[0] - drill.LayerHeights[len(drill.LayerHeights)-1]
	drill.IsValid()
}
func (drill Drill) GetRec(drills []Drill) (x1, y1, x2, y2 float64) {
	x1, y1, x2, y2 = math.MaxFloat64, math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64
	for _, d := range drills {
		x1 = math.Min(x1, d.X)
		y1 = math.Min(y1, d.Y)
		x2 = math.Max(x2, d.X)
		y2 = math.Max(y2, d.Y)
	}
	return
}
func decimal(value float64) float64 {
	value = math.Trunc(value*1e2+0.5) * 1e-2
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}
func (drill Drill) NearKDrills(drillSet []Drill, includeNum int) (nears []Drill) {

	if includeNum > len(drillSet)&& includeNum > 1 {
		includeNum = len(drillSet) - 1
		return drill.NearKDrills(drillSet, includeNum)
	}

	sort.Slice(drillSet, func(i, j int) bool {
		dist1 := drill.Distance(drillSet[i])
		dist2 := drill.Distance(drillSet[j])
		return dist1 < dist2
	})
	var start int
	if drillSet[0].Distance(drill) < 1e-1 {
	}
	for _, d := range drillSet {
		if d.Distance(drill) < 1e-1 {
			start++
		}
	}
	nears = make([]Drill, includeNum)
	copy(nears, drillSet[start:start+includeNum])
	return nears
}
func (drill Drill) NearDrills(drills []Drill, dist float64) (near []Drill, err error) {
	sort.Slice(drills, func(i, j int) bool {
		dist1 := drill.Distance(drills[i])
		dist2 := drill.Distance(drills[j])
		return dist1 < dist2
	})
	for _, d := range drills {
		dis := d.Distance(drill)
		if dis > 10e-1 && dis < dist {
			near = append(near, d)
		} else {
			break
		}
	}
	if len(near) == 0 {
		log.SetFlags(log.Lshortfile)
		return nil, fmt.Errorf(":search radius is too small and there is no drill in the search radius")
	}
	return near, nil
}
func SetLengthAndZ(drill *Drill, incidentDrills []Drill) {
	log.SetFlags(log.Lshortfile)
	var length, z, bottom float64 = 0, 0, 0
	for _, d := range incidentDrills {
		if d.BottomHeight() < bottom {
			bottom = d.BottomHeight()
		}
	}
	for idx := 0; idx < len(incidentDrills); idx++ {
		length += incidentDrills[idx].GetLength() * incidentDrills[idx].GetWeight()
		z += incidentDrills[idx].Z * incidentDrills[idx].GetWeight()
	}
	if length > drill.Z-bottom {
		length = drill.Z - bottom
	}
	drill.SetLength(decimal(length))
	drill.SetZ(decimal(z))
	if drill.BottomHeight() < bottom {
		drill.SetLength(decimal(drill.Z - bottom))
	}
	if drill.Z <= drill.BottomHeight() {
		debug.PrintStack()
		drill.Display()
		log.Fatal("error")
	}
}
func (drill Drill) Trunc(depth float64) (drill2 Drill) {
	var layers []int
	var heights []float64
	for idx, h := range drill.LayerHeights {
		if h > depth {
			layers = append(layers, drill.Layers[idx])
			heights = append(heights, drill.LayerHeights[idx])
		} else if h == depth {
			layers = append(layers, drill.Layers[idx])
			heights = append(heights, drill.LayerHeights[idx])
			break
		} else {
			layers = append(layers, drill.Layers[idx])
			heights = append(heights, depth)
			break
		}
	}
	drill2 = drill
	drill2.LayerHeights = heights
	drill2.Layers = layers
	return drill2
}
