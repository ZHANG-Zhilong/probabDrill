<<<<<<< HEAD
package utils

import (
	"awesome/internal/constant"
	"awesome/internal/entity"
	"fmt"
	"log"
	"math"
	"sort"
)

func DisplayDrills(drills []entity.Drill) {
	for _, d := range drills {
		d.Print()
	}
	fmt.Printf("total %d drills.", len(drills))
}
func GetGridDrills(drills []entity.Drill) (virtualDrills []entity.Drill) {

	px, py := constant.ResXY, constant.ResXY
	l, r, t, b := getDrillsRecXOY(drills)

	//grid to interpolate
	gridx, gridy := getGrid(px, py, l, r, t, b)
	fmt.Println(gridx, gridy)
	blocks := getBlockHeights(drills, constant.ResZ)
	bx, by := constant.GetBoundary()

	for idx := range gridx {
		for idy := range gridy {
			x := gridx[idx]
			y := gridy[idy]
			if isInPolygon(bx, by, x, y) {
				virtualDrills = append(virtualDrills, generateVirtualDrill(x, y, blocks))
			}
		}
	}
	for idx, _ := range bx {
		x, y := bx[idx], by[idx]
		virtualDrills = append(virtualDrills, generateVirtualDrill(x, y, blocks))
	}
	return
}
func setLengthAndZ(drill *entity.Drill, incidentDrills []entity.Drill) {
	var length, z float64 = 0, 0
	for idx := 0; idx < len(incidentDrills); idx++ {
		length += incidentDrills[idx].GetLength() * incidentDrills[idx].GetWeight()
		z += incidentDrills[idx].Z * incidentDrills[idx].GetWeight()
	}
	drill.SetLength(length)
	drill.SetZ(z)
}
func generateVirtualDrill(x, y float64, blocks []float64) (virtualDrill entity.Drill) {
	log.SetFlags(log.Lshortfile)
	virtualDrill = virtualDrill.MakeDrill(constant.GenVirtualDrillName(), x, y, 0)
	incidentDrills := obtainIncidentDrills(virtualDrill, constant.RadiusIn)
	setClassicalIdwWeights(virtualDrill, incidentDrills)
	setLengthAndZ(&virtualDrill, incidentDrills)
	virtualDrill.LayerFloorHeights = explodedHeights(blocks, virtualDrill.Z, virtualDrill.GetBottomHeight())
	//virtual(name, x, y, z, length, heights, weight)  还差 layers,

	//p(layer|block)
	var probLayerBlocks3s [][]float64
	probLayerBlocks3s = append(probLayerBlocks3s, make([]float64, constant.StdLen, constant.StdLen))

	for bidx := 1; bidx < len(blocks); bidx++ { //traverse general blocks.
		//p(layer|block0)
		//bidx = 107
		ceil, floor := blocks[bidx-1], blocks[bidx]
		var probBlockWithWeight = statProbBlockWithWeight(incidentDrills, ceil, floor)
		var probLayersWithWeight = make([]float64, constant.StdLen, constant.StdLen)
		var probBlockLayers = make([]float64, constant.StdLen, constant.StdLen)
		var probLayerBlock2s = make([]float64, constant.StdLen, constant.StdLen)

		for layerIdx := int64(1); layerIdx < constant.StdLen; layerIdx++ { //layer[0] is ground.
			//layerIdx = 26
			probLayersWithWeight[layerIdx] = statProbLayerWithWeight(incidentDrills, ceil, floor, layerIdx)
			probBlockLayers[layerIdx] = statProbBlockLayer(constant.DrillSet(), ceil, floor, layerIdx)

			if probBlockWithWeight >= 0.0000001 {
				probLayerBlock2s[layerIdx] = probBlockLayers[layerIdx] * probLayersWithWeight[layerIdx] / probBlockWithWeight
			}
			a, b, c := probLayersWithWeight[layerIdx], probBlockLayers[layerIdx], probLayerBlock2s[layerIdx]
			hole(a, b, c)
		}
		probLayerBlocks3s = append(probLayerBlocks3s, probLayerBlock2s)

		//log.Printf("blocks ceil:%.2f, floor:%.2f, p(block):%f\n", ceil, floor, probBlockWithWeight)
		//printFloat64s(probLayersWithWeight)
		//printFloat64s(probBlockLayers)
		//printFloat64s(probLayerBlock2s)

		//idx, val := findMaxFloat64s(probLayerBlock2s)
		//log.Println(idx, val)
		//fmt.Println("======")
	}

	//log.Println("initial drill.")
	//virtualDrill.Print()

	for idx := 1; idx < len(virtualDrill.LayerFloorHeights); idx++ {
		ceil, floor := virtualDrill.LayerFloorHeights[idx-1], virtualDrill.LayerFloorHeights[idx]
		bidx := blocksIndex(blocks, ceil, floor)
		probs := probLayerBlocks3s[bidx]
		layer, prob := findMaxFloat64s(probs)
		virtualDrill.Layers = append(virtualDrill.Layers, int64(layer))
		hole(prob)
		//log.Println(bidx, ceil, floor, layer, prob)
		//log.Println(probs)
	}

	//log.Println("before merged.")
	//virtualDrill.Print()
	virtualDrill = mergeDrill(virtualDrill)
	//log.Println("after merged.")
	//if !virtualDrill.IsValid() {
	//	virtualDrill.Print()
	//	log.Fatal("invalid drill")
	//}
	//virtualDrill.Print()

	//log.Println("p(layers_j|block_i)")
	//for i, p := range probLayerBlocks3s {
	//	idx, val := findMaxFloat64s(p)
	//	fmt.Printf("%d: ", i)
	//	fmt.Printf("[%d, %.2f] ", idx, val)
	//	printFloat64s(p)
	//}

	//log.Println("incident drills")
	//for _, d := range incidentDrills {
	//	d.Print()
	//	//log.Printf("%+v", d)
	//}
	return
}
func getGrid(px, py, l, r, t, b float64) (gridx, gridy []float64) {
	gridx = append(gridx, l)
	gridy = append(gridy, b)
	for (l + px) < r {
		gridx = append(gridx, l+px)
		l = l + px
	}
	for (b + py) < t {
		gridy = append(gridy, b+py)
		b = b + py
	}
	gridx = append(gridx, r)
	gridy = append(gridy, t)
	return
}
func findMaxFloat64s(float64s []float64) (idx int, val float64) {
	if len(float64s) < 1 {
		return 0, 0
	}
	idx, val = -math.MaxInt64, -math.MaxFloat64
	for id, va := range float64s {
		if va > val {
			idx, val = id, va
		}
	}
	return idx, val
}
func setClassicalIdwWeights(center entity.Drill, aroundDrills []entity.Drill) (weights []float64) {
	var (
		weightSum       float64
		hasZeroDistance bool
		zeroIdx         int
	)

	//get distance
	for idx, aroundDrill := range aroundDrills {
		if dist, ok := drillDist(center, aroundDrill); ok {
			weights = append(weights, dist)
			if dist < 0.0001 {
				hasZeroDistance = true
				zeroIdx = idx
			}
		}
	}
	if hasZeroDistance {
		for idx, _ := range weights {
			weights[idx] = 0
		}
		weights[zeroIdx] = 1
	} else {
		for idx, _ := range weights {
			weights[idx] = 1 / math.Pow(weights[idx], constant.IdwPow)
			weightSum += weights[idx]
		}
		for idx, _ := range weights {
			weights[idx] = weights[idx] / weightSum
			aroundDrills[idx].SetWeight(weights[idx])
		}
	}

	if weightSum > 1+0.0000001 && weightSum < 1-0.0000001 {
		log.SetFlags(log.Lshortfile)
		log.Println(weights)
		log.Fatalf("error: total weight:%f\n", weightSum)
	}
	return weights
}
func obtainIncidentDrills(drill entity.Drill, includeNum int) (drills []entity.Drill) {

	drillSet := constant.DrillSet()

	dists := make([]float64, len(drillSet), len(drillSet))
	for i, d := range drillSet {
		dists[i], _ = drillDist(drill, d)
	}

	sort.Float64s(dists)
	radius := dists[includeNum-1]

	for _, d := range drillSet {
		if distance, ok := drillDist(drill, d); ok && distance <= radius && d.Name != drill.Name {
			drills = append(drills, d)
		}
	}
	return drills
}
func hasLayer(drill entity.Drill, layer int64) (num int) {
	for _, l := range drill.Layers {
		if l == layer {
			num++
		}
	}
	return num
}
func hasBlock(drill entity.Drill, ceil, floor float64) (has bool) {
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
func heightRange(drills []entity.Drill) (ceil float64, floor float64) {
	ceil, floor = -math.MaxFloat64, math.MaxFloat64
	for _, d := range drills {
		if d.Z > ceil {
			ceil = d.Z
		}
		if d.LayerFloorHeights[len(d.LayerFloorHeights)-1] < floor {
			floor = d.LayerFloorHeights[len(d.LayerFloorHeights)-1]
		}
	}
	return ceil, floor
}
func getBlockHeights(drillSet []entity.Drill, resz float64) (blocksHeight []float64) {
	drillsCeil, drillsFloor := -math.MaxFloat64, math.MaxFloat64
	for _, d := range drillSet {
		if d.Z > drillsCeil {
			drillsCeil = d.Z
		}
		if d.LayerFloorHeights[len(d.LayerFloorHeights)-1] < drillsFloor {
			drillsFloor = d.LayerFloorHeights[len(d.LayerFloorHeights)-1]
		}
	}

	blocksHeight = append(blocksHeight, drillsCeil)

	for drillsCeil-resz > drillsFloor {
		blocksHeight = append(blocksHeight, drillsCeil-resz)
		drillsCeil = drillsCeil - resz
	}

	//the last block may be un-standard block length, whose length may less than resz
	blocksHeight = append(blocksHeight, drillsFloor)

	return
}
func getLayerSeq(drill entity.Drill, ceil, floor float64) (seq int64, ok bool) {
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
		return getLayerSeq(drill, drill.Z, floor)
	}

	//case 3.2
	if ceil > drill.GetBottomHeight() && floor < drill.GetWeight() {
		return getLayerSeq(drill, ceil, drill.GetBottomHeight())
	}
	return -1, false
}
func mergeDrill(drill entity.Drill) entity.Drill {
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
	return drill
}
func explodeDrill(drill entity.Drill, blocks []float64) (scatteredDrill entity.Drill) {
	if blocks == nil {
		return
	}
	if blocks[0] < drill.Z {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Println("abnormal Drills")
		drill.Print()
	}
	scatteredDrill = scatteredDrill.MakeDrill(drill.Name, drill.X, drill.Y, drill.Z)

	var idxa int
	for idx, h := range blocks {
		if h < drill.Z {
			idxa = idx
			break
		}
	}
	var drillBlocks []float64 = []float64{drill.Z}
	var drillLayers []int64 = []int64{entity.Ground}

	for idx := idxa; idx < len(blocks); idx++ {
		if blocks[idx] <= drill.Z && blocks[idx] >= drill.Z-drill.GetLength() {
			drillBlocks = append(drillBlocks, blocks[idx])
			if seq, ok := getLayerSeq(
				drill, drillBlocks[len(drillBlocks)-2], drillBlocks[len(drillBlocks)-1]); ok {
				drillLayers = append(drillLayers, seq)
			}
		}
	}
	if drillBlocks[len(drillBlocks)-1] != drill.Z-drill.GetLength() {
		drillBlocks = append(drillBlocks, drill.Z-drill.GetLength())
		if seq, ok := getLayerSeq(
			drill, drillBlocks[len(drillBlocks)-2], drillBlocks[len(drillBlocks)-1]); ok {
			drillLayers = append(drillLayers, seq)
		}
	}
	scatteredDrill.LayerFloorHeights = drillBlocks
	scatteredDrill.Layers = drillLayers
	scatteredDrill.SetWeight(drill.GetWeight())
	return
}
func drillDist(drill1, drill2 entity.Drill) (dist float64, ok bool) {
	x1, y1, x2, y2 := drill1.X, drill1.Y, drill2.X, drill2.Y
	dist = math.Sqrt((x2-x1)*(x2-x1) + (y2-y1)*(y2-y1))
	if math.IsNaN(dist) || math.IsInf(dist, 0) || dist < 0 {
		return -1, false
	}
	return dist, true
}
func getDrillsRecXOY(drills []entity.Drill) (l, r, t, b float64) {
	l, b = math.MaxFloat64, math.MaxFloat64
	r, t = -math.MaxFloat64, -math.MaxFloat64
	for _, drill := range drills {
		if drill.X < l {
			l = drill.X
		}
		if drill.X > r {
			r = drill.X
		}
		if drill.Y < b {
			b = drill.Y
		}
		if drill.Y > t {
			t = drill.Y
		}
	}
	return
}
func explodedHeights(blocks []float64, ceil, floor float64) (heights []float64) {
	idxa := int(0)
	for idx, h := range blocks {
		if h < ceil {
			idxa = idx
			break
		}
	}
	heights = append(heights, ceil)

	for idx := idxa; idx < len(blocks); idx++ {
		if blocks[idx] <= ceil && blocks[idx] >= floor {
			heights = append(heights, blocks[idx])
		}
	}
	if heights[len(heights)-1] > floor {
		heights = append(heights, floor)
	}
	return
}
func statProbBlock(drills []entity.Drill, ceil, floor float64) (prob float64) {
	for _, drill := range drills {
		if hasBlock(drill, ceil, floor) {
			prob += 1.0
		}
	}
	prob = prob / float64(len(drills))
	return prob
}
func statProbBlockWithWeight(drills []entity.Drill, blockCeil, blockFloor float64) (prob float64) {
	if len(drills) < 1 || blockCeil <= blockFloor {
		log.Fatal("error")
		return
	}
	for _, d := range drills {
		if hasBlock(d, blockCeil, blockFloor) {
			prob += d.GetWeight()
		}
	}
	if math.IsNaN(prob) || math.IsInf(prob, 0) {
		log.Fatal("error")
		return
	}
	return prob
}
func statProbBlockAndLayer(drills []entity.Drill, blockCeil, blockFloor float64, layer int64) (prob float64) {
	//p(blockAndLayer)
	log.SetFlags(log.Lshortfile)
	for _, drill := range drills {
		if seq, ok := getLayerSeq(drill, blockCeil, blockFloor); ok && seq == layer {
			prob += 1.0
		}
	}
	prob = prob / float64(len(drills))
	if math.IsNaN(prob) || math.IsInf(prob, 0) {
		log.Fatal("error")
	}
	return prob
}
func statProbBlockAndLayerWithWeight(drills []entity.Drill, blockCeil, blockFloor float64, layer int64) (prob float64) {
	for _, drill := range drills {
		if seq, ok := getLayerSeq(drill, blockCeil, blockFloor); ok && seq == layer {
			prob += drill.GetWeight()
		}
	}
	if math.IsNaN(prob) || math.IsInf(prob, 0) {
		log.SetFlags(log.Lshortfile)
		log.Fatal("error")
	}
	return prob
}
func statProbBlockLayerWithWeight(drills []entity.Drill, blockCeil, blockFloor float64, layer int64) (prob float64) {
	//p(block|layer)=p(blockAndLayer)/p(layer)= p(block ∩ layer)/p(layer)  //∩->\cap
	probLayer := statProbLayerWithWeight(drills, blockCeil, blockFloor, layer)
	probLayerAndBlock := statProbBlockAndLayerWithWeight(drills, blockCeil, blockFloor, layer)
	prob = probLayerAndBlock / probLayer
	if math.IsNaN(prob) || math.IsInf(prob, 0) {
		log.Fatal("error")
		return
	}
	if prob == 0 {
		log.SetFlags(log.LstdFlags | log.LstdFlags)
		log.Printf("p(layer) %f, p(layerAndblock) %f\n", probLayer, probLayerAndBlock)
	}
	return prob
}
func statProbBlockLayer(drills []entity.Drill, blockCeil, blockFloor float64, layer int64) (prob float64) {
	//p(block|layer) = p(blockAndLayer)/p(layer)
	probLayer := statProbLayer(drills, blockCeil, blockFloor, layer)
	probBlockAndLayer := statProbBlockAndLayer(drills, blockCeil, blockFloor, layer)
	if probLayer > 0 {
		prob = probBlockAndLayer / probLayer
	}
	return prob
}
func statProbLayer(drills []entity.Drill, blockCeil, blockFloor float64, layer int64) (prob float64) {
	log.SetFlags(log.Lshortfile)
	if len(drills) < 1 || blockCeil <= blockFloor {
		log.Fatal("error")
		return -1
	}

	//p(layer) = p(block)p(layer|block)+p(blank)p(blank|layer)
	//p(blank)+p(block)=1, p(layer1|block)+p(layer2|block)=1, p(layer1|blank)+p(layer2|blank)=1

	probBlock := statProbBlock(drills, blockCeil, blockFloor)
	probBlank := 1 - probBlock

	//here is a transformation, that p(block|layer) is general
	//p(layer) means the drill has layer, p(block) means drill has block
	probLayerBlock := statProbLayerBlock(drills, blockCeil, blockFloor, layer)
	probLayerBlank := 0.0

	prob = probBlock*probLayerBlock + probBlank*probLayerBlank

	if math.IsNaN(prob) || math.IsInf(prob, 0) {
		return
		log.Fatal("error")
	}
	if prob > 1 {
		prob = 1
	}
	return prob
}
func statProbLayerBlockWithWeight(drills []entity.Drill, blockCeil, blockFloor float64, layer int64) (prob float64) {
	//p(layer|block) = p(layerAndBlock)/p(block)
	probBlockAndLayerWithWeight := statProbBlockAndLayerWithWeight(drills, blockCeil, blockFloor, layer)
	probBlocksWithWeight := statProbBlockWithWeight(drills, blockCeil, blockFloor)
	if probBlocksWithWeight > 0 {
		prob = probBlockAndLayerWithWeight / probBlocksWithWeight
	}
	return prob
}
func statProbLayerBlock(drills []entity.Drill, blockCeil, blockFloor float64, layer int64) (prob float64) {
	//p(layer|block) = p(layerAndBlock)/p(block)
	probBlockAndLayer := statProbBlockAndLayer(drills, blockCeil, blockFloor, layer)
	probBlocks := statProbBlock(drills, blockCeil, blockFloor)
	if probBlocks > 0 {
		prob = probBlockAndLayer / probBlocks
	}
	return prob
}
func statProbLayerWithWeight(drills []entity.Drill, blockCeil, blockFloor float64, layer int64) (prob float64) {
	log.SetFlags(log.Lshortfile)
	if len(drills) < 1 || blockCeil <= blockFloor {
		log.Fatal("error")
		return -1
	}

	//p(layer) = p(block)p(layer|block)+p(blank)p(blank|layer)
	//p(blank)+p(block)=1, p(layer1|block)+p(layer2|block)=1, p(layer1|blank)+p(layer2|blank)=1

	probBlockWithWeight := statProbBlockWithWeight(drills, blockCeil, blockFloor)
	probBlankWithWeight := 1 - probBlockWithWeight

	//here is a transformation, that p(block|layer) is general
	//p(layer) means the drill has layer, p(block) means drill has block

	probLayerBlockWithWeight := statProbLayerBlockWithWeight(drills, blockCeil, blockFloor, layer)
	probLayerBlankWithWeight := 0.0

	prob = probBlockWithWeight*probLayerBlockWithWeight + probBlankWithWeight*probLayerBlankWithWeight

	if math.IsNaN(prob) || math.IsInf(prob, 0) {
		return -1
		log.Fatal("error")
	}
	return prob
}
func printFloat64s(s []float64) () {
	fmt.Print("[")
	for _, v := range s {
		if v > 0 {
			fmt.Printf("%.3f ", v)
		} else {
			fmt.Printf("%.3f ", v)
		}
	}
	fmt.Print("]\n")
}
func blocksIndex(blocks []float64, ceil, floor float64) (index int) {
	log.SetFlags(log.Lshortfile)
	for idx := 1; idx < len(blocks); idx++ {
		if ceil <= blocks[idx-1] && floor >= blocks[idx] {
			return idx
		}
	}
	return -1
}
func hole(vals ...float64) () {
	return
}
func isInPolygon(x, y []float64, x0, y0 float64) (isIn bool) {

	//vert[0], vert[last]
	var i, j int = 0, len(x) - 1
	if (y[i] >= y0) != (y[j] > y0) &&
		(y0 <= y[i] && y0 <= y[j] ||
			x0 <= (y0-y[i])*(x[j]-x[i])/(y[j]-y[i])+x[i]) {
		isIn = !isIn
	}

	//y0 is among y1 and y2, ray x0
	//if k=inf -> y1==y2  y0<=y1&&y0<y2 cross
	//if k< inf	x0<x1+k(y0-y1) cross
	for i := 1; i < len(x); i++ {
		if (y[i] >= y0) != (y[j] > y0) &&
			(y0 <= y[i] && y0 <= y[j] ||
				x0 <= (y0-y[i])*(x[j]-x[i])/(y[j]-y[i])+x[i]) {
			isIn = !isIn
		}
	}

	return isIn
}
=======
package utils

import (
	"fmt"
	"log"
	"math"
	"probabDrill-main/internal/constant"
	"probabDrill-main/internal/entity"
	"sort"
)

func DisplayDrills(drills []entity.Drill) {
	for _, d := range drills {
		d.Print()
	}
	fmt.Printf("total %d drills.", len(drills))
}
func GetGridDrills(drills []entity.Drill) (virtualDrills []entity.Drill) {

	px, py := constant.ResXY, constant.ResXY
	l, r, t, b := getDrillsRecXOY(drills)

	//grid to interpolate
	gridx, gridy := getGrid(px, py, l, r, t, b)
	fmt.Println(gridx, gridy)
	blocks := getBlockHeights(drills, constant.ResZ)
	bx, by := constant.GetBoundary()

	for idx := range gridx {
		for idy := range gridy {
			x := gridx[idx]
			y := gridy[idy]
			if isInPolygon(bx, by, x, y) {
				virtualDrills = append(virtualDrills, generateVirtualDrill(x, y, blocks))
			}
		}
	}
	//for idx, _ := range bx {
	//	x, y := bx[idx], by[idx]
	//	virtualDrills = append(virtualDrills, generateVirtualDrill(x, y, blocks))
	//}
	return
}
func setLengthAndZ(drill *entity.Drill, incidentDrills []entity.Drill) {
	var length, z float64 = 0, 0
	for idx := 0; idx < len(incidentDrills); idx++ {
		length += incidentDrills[idx].GetLength() * incidentDrills[idx].GetWeight()
		z += incidentDrills[idx].Z * incidentDrills[idx].GetWeight()
	}
	drill.SetLength(length)
	drill.SetZ(z)
}
func generateVirtualDrill(x, y float64, blocks []float64) (virtualDrill entity.Drill) {
	log.SetFlags(log.Lshortfile)
	virtualDrill = virtualDrill.MakeDrill(constant.GenVirtualDrillName(), x, y, 0)
	incidentDrills := obtainIncidentDrills(virtualDrill, constant.RadiusIn)
	setClassicalIdwWeights(virtualDrill, incidentDrills)
	setLengthAndZ(&virtualDrill, incidentDrills)
	virtualDrill.LayerFloorHeights = explodedHeights(blocks, virtualDrill.Z, virtualDrill.GetBottomHeight())
	//virtual(name, x, y, z, length, heights, weight)  还差 layers,

	//p(layer|block)
	var probLayerBlocks3s [][]float64
	probLayerBlocks3s = append(probLayerBlocks3s, make([]float64, constant.StdLen, constant.StdLen))

	for bidx := 1; bidx < len(blocks); bidx++ { //traverse general blocks.
		//p(layer|block0)
		//bidx = 107
		ceil, floor := blocks[bidx-1], blocks[bidx]
		var probBlockWithWeight = statProbBlockWithWeight(incidentDrills, ceil, floor)
		var probLayersWithWeight = make([]float64, constant.StdLen, constant.StdLen)
		var probBlockLayers = make([]float64, constant.StdLen, constant.StdLen)
		var probLayerBlock2s = make([]float64, constant.StdLen, constant.StdLen)

		for layerIdx := int64(1); layerIdx < constant.StdLen; layerIdx++ { //layer[0] is ground.
			//layerIdx = 26
			probLayersWithWeight[layerIdx] = statProbLayerWithWeight(incidentDrills, ceil, floor, layerIdx)
			probBlockLayers[layerIdx] = statProbBlockLayer(constant.DrillSet(), ceil, floor, layerIdx)

			if probBlockWithWeight >= 0.0000001 {
				probLayerBlock2s[layerIdx] = probBlockLayers[layerIdx] * probLayersWithWeight[layerIdx] / probBlockWithWeight
			}
			a, b, c := probLayersWithWeight[layerIdx], probBlockLayers[layerIdx], probLayerBlock2s[layerIdx]
			hole(a, b, c)
		}
		probLayerBlocks3s = append(probLayerBlocks3s, probLayerBlock2s)

		//log.Printf("blocks ceil:%.2f, floor:%.2f, p(block):%f\n", ceil, floor, probBlockWithWeight)
		//printFloat64s(probLayersWithWeight)
		//printFloat64s(probBlockLayers)
		//printFloat64s(probLayerBlock2s)

		//idx, val := findMaxFloat64s(probLayerBlock2s)
		//log.Println(idx, val)
		//fmt.Println("======")
	}

	//log.Println("initial drill.")
	//virtualDrill.Print()

	for idx := 1; idx < len(virtualDrill.LayerFloorHeights); idx++ {
		ceil, floor := virtualDrill.LayerFloorHeights[idx-1], virtualDrill.LayerFloorHeights[idx]
		bidx := blocksIndex(blocks, ceil, floor)
		probs := probLayerBlocks3s[bidx]
		layer, prob := findMaxFloat64s(probs)
		virtualDrill.Layers = append(virtualDrill.Layers, int64(layer))
		hole(prob)
		//log.Println(bidx, ceil, floor, layer, prob)
		//log.Println(probs)
	}

	//log.Println("before merged.")
	//virtualDrill.Print()
	virtualDrill = mergeDrill(virtualDrill)
	//log.Println("after merged.")
	//if !virtualDrill.IsValid() {
	//	virtualDrill.Print()
	//	log.Fatal("invalid drill")
	//}
	//virtualDrill.Print()

	//log.Println("p(layers_j|block_i)")
	//for i, p := range probLayerBlocks3s {
	//	idx, val := findMaxFloat64s(p)
	//	fmt.Printf("%d: ", i)
	//	fmt.Printf("[%d, %.2f] ", idx, val)
	//	printFloat64s(p)
	//}

	//log.Println("incident drills")
	//for _, d := range incidentDrills {
	//	d.Print()
	//	//log.Printf("%+v", d)
	//}
	return
}
func getGrid(px, py, l, r, t, b float64) (gridx, gridy []float64) {
	gridx = append(gridx, l)
	gridy = append(gridy, b)
	for (l + px) < r {
		gridx = append(gridx, l+px)
		l = l + px
	}
	for (b + py) < t {
		gridy = append(gridy, b+py)
		b = b + py
	}
	gridx = append(gridx, r)
	gridy = append(gridy, t)
	return
}
func findMaxFloat64s(float64s []float64) (idx int, val float64) {
	if len(float64s) < 1 {
		return 0, 0
	}
	idx, val = -math.MaxInt64, -math.MaxFloat64
	for id, va := range float64s {
		if va > val {
			idx, val = id, va
		}
	}
	return idx, val
}
func setClassicalIdwWeights(center entity.Drill, aroundDrills []entity.Drill) (weights []float64) {
	var (
		weightSum       float64
		hasZeroDistance bool
		zeroIdx         int
	)

	//get distance
	for idx, aroundDrill := range aroundDrills {
		if dist, ok := drillDist(center, aroundDrill); ok {
			weights = append(weights, dist)
			if dist < 0.0001 {
				hasZeroDistance = true
				zeroIdx = idx
			}
		}
	}
	if hasZeroDistance {
		for idx, _ := range weights {
			weights[idx] = 0
		}
		weights[zeroIdx] = 1
	} else {
		for idx, _ := range weights {
			weights[idx] = 1 / math.Pow(weights[idx], constant.IdwPow)
			weightSum += weights[idx]
		}
		for idx, _ := range weights {
			weights[idx] = weights[idx] / weightSum
			aroundDrills[idx].SetWeight(weights[idx])
		}
	}

	if weightSum > 1+0.0000001 && weightSum < 1-0.0000001 {
		log.SetFlags(log.Lshortfile)
		log.Println(weights)
		log.Fatalf("error: total weight:%f\n", weightSum)
	}
	return weights
}
func obtainIncidentDrills(drill entity.Drill, includeNum int) (drills []entity.Drill) {

	drillSet := constant.DrillSet()

	dists := make([]float64, len(drillSet), len(drillSet))
	for i, d := range drillSet {
		dists[i], _ = drillDist(drill, d)
	}

	sort.Float64s(dists)
	radius := dists[includeNum-1]

	for _, d := range drillSet {
		if distance, ok := drillDist(drill, d); ok && distance <= radius && d.Name != drill.Name {
			drills = append(drills, d)
		}
	}
	return drills
}
func hasLayer(drill entity.Drill, layer int64) (num int) {
	for _, l := range drill.Layers {
		if l == layer {
			num++
		}
	}
	return num
}
func hasBlock(drill entity.Drill, ceil, floor float64) (has bool) {
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
func heightRange(drills []entity.Drill) (ceil float64, floor float64) {
	ceil, floor = -math.MaxFloat64, math.MaxFloat64
	for _, d := range drills {
		if d.Z > ceil {
			ceil = d.Z
		}
		if d.LayerFloorHeights[len(d.LayerFloorHeights)-1] < floor {
			floor = d.LayerFloorHeights[len(d.LayerFloorHeights)-1]
		}
	}
	return ceil, floor
}
func getBlockHeights(drillSet []entity.Drill, resz float64) (blocksHeight []float64) {
	drillsCeil, drillsFloor := -math.MaxFloat64, math.MaxFloat64
	for _, d := range drillSet {
		if d.Z > drillsCeil {
			drillsCeil = d.Z
		}
		if d.LayerFloorHeights[len(d.LayerFloorHeights)-1] < drillsFloor {
			drillsFloor = d.LayerFloorHeights[len(d.LayerFloorHeights)-1]
		}
	}

	blocksHeight = append(blocksHeight, drillsCeil)

	for drillsCeil-resz > drillsFloor {
		blocksHeight = append(blocksHeight, drillsCeil-resz)
		drillsCeil = drillsCeil - resz
	}

	//the last block may be un-standard block length, whose length may less than resz
	blocksHeight = append(blocksHeight, drillsFloor)

	return
}
func getLayerSeq(drill entity.Drill, ceil, floor float64) (seq int64, ok bool) {
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
		return getLayerSeq(drill, drill.Z, floor)
	}

	//case 3.2
	if ceil > drill.GetBottomHeight() && floor < drill.GetWeight() {
		return getLayerSeq(drill, ceil, drill.GetBottomHeight())
	}
	return -1, false
}
func mergeDrill(drill entity.Drill) entity.Drill {
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
	return drill
}
func explodeDrill(drill entity.Drill, blocks []float64) (scatteredDrill entity.Drill) {
	if blocks == nil {
		return
	}
	if blocks[0] < drill.Z {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Println("abnormal Drills")
		drill.Print()
	}
	scatteredDrill = scatteredDrill.MakeDrill(drill.Name, drill.X, drill.Y, drill.Z)

	var idxa int
	for idx, h := range blocks {
		if h < drill.Z {
			idxa = idx
			break
		}
	}
	var drillBlocks []float64 = []float64{drill.Z}
	var drillLayers []int64 = []int64{entity.Ground}

	for idx := idxa; idx < len(blocks); idx++ {
		if blocks[idx] <= drill.Z && blocks[idx] >= drill.Z-drill.GetLength() {
			drillBlocks = append(drillBlocks, blocks[idx])
			if seq, ok := getLayerSeq(
				drill, drillBlocks[len(drillBlocks)-2], drillBlocks[len(drillBlocks)-1]); ok {
				drillLayers = append(drillLayers, seq)
			}
		}
	}
	if drillBlocks[len(drillBlocks)-1] != drill.Z-drill.GetLength() {
		drillBlocks = append(drillBlocks, drill.Z-drill.GetLength())
		if seq, ok := getLayerSeq(
			drill, drillBlocks[len(drillBlocks)-2], drillBlocks[len(drillBlocks)-1]); ok {
			drillLayers = append(drillLayers, seq)
		}
	}
	scatteredDrill.LayerFloorHeights = drillBlocks
	scatteredDrill.Layers = drillLayers
	scatteredDrill.SetWeight(drill.GetWeight())
	return
}
func drillDist(drill1, drill2 entity.Drill) (dist float64, ok bool) {
	x1, y1, x2, y2 := drill1.X, drill1.Y, drill2.X, drill2.Y
	dist = math.Sqrt((x2-x1)*(x2-x1) + (y2-y1)*(y2-y1))
	if math.IsNaN(dist) || math.IsInf(dist, 0) || dist < 0 {
		return -1, false
	}
	return dist, true
}
func getDrillsRecXOY(drills []entity.Drill) (l, r, t, b float64) {
	l, b = math.MaxFloat64, math.MaxFloat64
	r, t = -math.MaxFloat64, -math.MaxFloat64
	for _, drill := range drills {
		if drill.X < l {
			l = drill.X
		}
		if drill.X > r {
			r = drill.X
		}
		if drill.Y < b {
			b = drill.Y
		}
		if drill.Y > t {
			t = drill.Y
		}
	}
	return
}
func explodedHeights(blocks []float64, ceil, floor float64) (heights []float64) {
	idxa := int(0)
	for idx, h := range blocks {
		if h < ceil {
			idxa = idx
			break
		}
	}
	heights = append(heights, ceil)

	for idx := idxa; idx < len(blocks); idx++ {
		if blocks[idx] <= ceil && blocks[idx] >= floor {
			heights = append(heights, blocks[idx])
		}
	}
	if heights[len(heights)-1] > floor {
		heights = append(heights, floor)
	}
	return
}
func statProbBlock(drills []entity.Drill, ceil, floor float64) (prob float64) {
	for _, drill := range drills {
		if hasBlock(drill, ceil, floor) {
			prob += 1.0
		}
	}
	prob = prob / float64(len(drills))
	return prob
}
func statProbBlockWithWeight(drills []entity.Drill, blockCeil, blockFloor float64) (prob float64) {
	if len(drills) < 1 || blockCeil <= blockFloor {
		log.Fatal("error")
		return
	}
	for _, d := range drills {
		if hasBlock(d, blockCeil, blockFloor) {
			prob += d.GetWeight()
		}
	}
	if math.IsNaN(prob) || math.IsInf(prob, 0) {
		log.Fatal("error")
		return
	}
	return prob
}
func statProbBlockAndLayer(drills []entity.Drill, blockCeil, blockFloor float64, layer int64) (prob float64) {
	//p(blockAndLayer)
	log.SetFlags(log.Lshortfile)
	for _, drill := range drills {
		if seq, ok := getLayerSeq(drill, blockCeil, blockFloor); ok && seq == layer {
			prob += 1.0
		}
	}
	prob = prob / float64(len(drills))
	if math.IsNaN(prob) || math.IsInf(prob, 0) {
		log.Fatal("error")
	}
	return prob
}
func statProbBlockAndLayerWithWeight(drills []entity.Drill, blockCeil, blockFloor float64, layer int64) (prob float64) {
	for _, drill := range drills {
		if seq, ok := getLayerSeq(drill, blockCeil, blockFloor); ok && seq == layer {
			prob += drill.GetWeight()
		}
	}
	if math.IsNaN(prob) || math.IsInf(prob, 0) {
		log.SetFlags(log.Lshortfile)
		log.Fatal("error")
	}
	return prob
}
func statProbBlockLayerWithWeight(drills []entity.Drill, blockCeil, blockFloor float64, layer int64) (prob float64) {
	//p(block|layer)=p(blockAndLayer)/p(layer)= p(block ∩ layer)/p(layer)  //∩->\cap
	probLayer := statProbLayerWithWeight(drills, blockCeil, blockFloor, layer)
	probLayerAndBlock := statProbBlockAndLayerWithWeight(drills, blockCeil, blockFloor, layer)
	prob = probLayerAndBlock / probLayer
	if math.IsNaN(prob) || math.IsInf(prob, 0) {
		log.Fatal("error")
		return
	}
	if prob == 0 {
		log.SetFlags(log.LstdFlags | log.LstdFlags)
		log.Printf("p(layer) %f, p(layerAndblock) %f\n", probLayer, probLayerAndBlock)
	}
	return prob
}
func statProbBlockLayer(drills []entity.Drill, blockCeil, blockFloor float64, layer int64) (prob float64) {
	//p(block|layer) = p(blockAndLayer)/p(layer)
	probLayer := statProbLayer(drills, blockCeil, blockFloor, layer)
	probBlockAndLayer := statProbBlockAndLayer(drills, blockCeil, blockFloor, layer)
	if probLayer > 0 {
		prob = probBlockAndLayer / probLayer
	}
	return prob
}
func statProbLayer(drills []entity.Drill, blockCeil, blockFloor float64, layer int64) (prob float64) {
	log.SetFlags(log.Lshortfile)
	if len(drills) < 1 || blockCeil <= blockFloor {
		log.Fatal("error")
		return -1
	}

	//p(layer) = p(block)p(layer|block)+p(blank)p(blank|layer)
	//p(blank)+p(block)=1, p(layer1|block)+p(layer2|block)=1, p(layer1|blank)+p(layer2|blank)=1

	probBlock := statProbBlock(drills, blockCeil, blockFloor)
	probBlank := 1 - probBlock

	//here is a transformation, that p(block|layer) is general
	//p(layer) means the drill has layer, p(block) means drill has block
	probLayerBlock := statProbLayerBlock(drills, blockCeil, blockFloor, layer)
	probLayerBlank := 0.0

	prob = probBlock*probLayerBlock + probBlank*probLayerBlank

	if math.IsNaN(prob) || math.IsInf(prob, 0) {
		return
		log.Fatal("error")
	}
	if prob > 1 {
		prob = 1
	}
	return prob
}
func statProbLayerBlockWithWeight(drills []entity.Drill, blockCeil, blockFloor float64, layer int64) (prob float64) {
	//p(layer|block) = p(layerAndBlock)/p(block)
	probBlockAndLayerWithWeight := statProbBlockAndLayerWithWeight(drills, blockCeil, blockFloor, layer)
	probBlocksWithWeight := statProbBlockWithWeight(drills, blockCeil, blockFloor)
	if probBlocksWithWeight > 0 {
		prob = probBlockAndLayerWithWeight / probBlocksWithWeight
	}
	return prob
}
func statProbLayerBlock(drills []entity.Drill, blockCeil, blockFloor float64, layer int64) (prob float64) {
	//p(layer|block) = p(layerAndBlock)/p(block)
	probBlockAndLayer := statProbBlockAndLayer(drills, blockCeil, blockFloor, layer)
	probBlocks := statProbBlock(drills, blockCeil, blockFloor)
	if probBlocks > 0 {
		prob = probBlockAndLayer / probBlocks
	}
	return prob
}
func statProbLayerWithWeight(drills []entity.Drill, blockCeil, blockFloor float64, layer int64) (prob float64) {
	log.SetFlags(log.Lshortfile)
	if len(drills) < 1 || blockCeil <= blockFloor {
		log.Fatal("error")
		return -1
	}

	//p(layer) = p(block)p(layer|block)+p(blank)p(blank|layer)
	//p(blank)+p(block)=1, p(layer1|block)+p(layer2|block)=1, p(layer1|blank)+p(layer2|blank)=1

	probBlockWithWeight := statProbBlockWithWeight(drills, blockCeil, blockFloor)
	probBlankWithWeight := 1 - probBlockWithWeight

	//here is a transformation, that p(block|layer) is general
	//p(layer) means the drill has layer, p(block) means drill has block

	probLayerBlockWithWeight := statProbLayerBlockWithWeight(drills, blockCeil, blockFloor, layer)
	probLayerBlankWithWeight := 0.0

	prob = probBlockWithWeight*probLayerBlockWithWeight + probBlankWithWeight*probLayerBlankWithWeight

	if math.IsNaN(prob) || math.IsInf(prob, 0) {
		return -1
		log.Fatal("error")
	}
	return prob
}
func printFloat64s(s []float64) () {
	fmt.Print("[")
	for _, v := range s {
		if v > 0 {
			fmt.Printf("%.3f ", v)
		} else {
			fmt.Printf("%.3f ", v)
		}
	}
	fmt.Print("]\n")
}
func blocksIndex(blocks []float64, ceil, floor float64) (index int) {
	log.SetFlags(log.Lshortfile)
	for idx := 1; idx < len(blocks); idx++ {
		if ceil <= blocks[idx-1] && floor >= blocks[idx] {
			return idx
		}
	}
	return -1
}
func hole(vals ...float64) () {
	return
}
func isInPolygon(x, y []float64, x0, y0 float64) (isIn bool) {

	//vert[0], vert[last]
	var i, j int = 0, len(x) - 1
	if (y[i] >= y0) != (y[j] > y0) &&
		(y0 <= y[i] && y0 <= y[j] ||
			x0 <= (y0-y[i])*(x[j]-x[i])/(y[j]-y[i])+x[i]) {
		isIn = !isIn
	}

	//y0 is among y1 and y2, ray x0
	//if k=inf -> y1==y2  y0<=y1&&y0<y2 cross
	//if k< inf	x0<x1+k(y0-y1) cross
	for i := 1; i < len(x); i++ {
		if (y[i] >= y0) != (y[j] > y0) &&
			(y0 <= y[i] && y0 <= y[j] ||
				x0 <= (y0-y[i])*(x[j]-x[i])/(y[j]-y[i])+x[i]) {
			isIn = !isIn
		}
	}

	return isIn
}
>>>>>>> 34e6069b56f966af97c9e8c24edc3db14aa285d2
