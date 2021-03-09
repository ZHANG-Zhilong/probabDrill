package service

import (
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"
	"log"
	"math"
	"probabDrill/internal/constant"
	"probabDrill/internal/entity"
	"probabDrill/internal/stat"
	"probabDrill/internal/utils"
	"runtime/debug"
	"sort"
)

func GetGridDrills(drillSet []entity.Drill) (virtualDrills []entity.Drill) {

	px, py := constant.GridXY, constant.GridXY
	l, r, t, b := getDrillsRecXOY(drillSet)

	//grid to interpolate
	gridx, gridy := utils.GetGrids(px, py, l, r, t, b)
	log.Println(gridx)
	log.Println(gridy)
	bx, by := constant.GetBoundary()
	var in, out int
	for idx := range gridx {
		for idy := range gridy {
			x := gridx[idx]
			y := gridy[idy]
			if utils.IsInPolygon(bx, by, x, y) {
				in++
				//virtualDrills = append(virtualDrills, GenVDrillM1(drillSet, x, y, blocks))
				virtualDrills = append(virtualDrills, GenVDrillM2(&drillSet, x, y))
			} else {
				out++
			}
		}
	}
	log.Println("drillIn:", in, " drillOut:", out)
	for idx, _ := range bx {
		x, y := bx[idx], by[idx]
		virtualDrills = append(virtualDrills, GenVDrillM1(&drillSet, x, y))
	}
	return
}
func setLengthAndZ(drill *entity.Drill, incidentDrills []entity.Drill) {
	log.SetFlags(log.Lshortfile)
	var length, z, bottom float64 = 0, 0, 0
	for _, d := range incidentDrills {
		if d.GetBottomHeight() < bottom {
			bottom = d.GetBottomHeight()
		}
	}
	for idx := 0; idx < len(incidentDrills); idx++ {
		length += incidentDrills[idx].GetLength() * incidentDrills[idx].GetWeight()
		z += incidentDrills[idx].Z * incidentDrills[idx].GetWeight()
	}
	if length > drill.Z-bottom {
		length = drill.Z - bottom
	}
	drill.SetLength(length)
	drill.SetZ(z)
	if drill.GetBottomHeight() < bottom {
		drill.SetLength(drill.Z - bottom)
	}
	if drill.Z <= drill.GetBottomHeight() {
		debug.PrintStack()
		drill.Print()
		log.Fatal("error")
	}
}
func GenVDrillM1(drillSet *[]entity.Drill, x, y float64) (virtualDrill entity.Drill) {
	log.SetFlags(log.Lshortfile)
	virtualDrill = virtualDrill.MakeDrill(constant.GenVDrillName(), x, y, 0)
	nearDrills := obtainNearDrills(drillSet, virtualDrill, constant.RadiusIn)
	for _, d := range *nearDrills {
		if math.Abs(x-d.X) < 0.001 && math.Abs(y-d.Y) < 0.001 {
			return d
		}
	}
	setClassicalIdwWeights(virtualDrill, *nearDrills)
	setLengthAndZ(&virtualDrill, *nearDrills)
	blocks := MakeBlocks(*drillSet, 0.02)
	virtualDrill.LayerHeights = explodedHeights(blocks, virtualDrill.Z, virtualDrill.GetBottomHeight())
	//virtual(name, x, y, z, length, heights, weight)  还差 layers,

	//p(layer|block)
	var probLayerBlocks3s [][]float64
	probLayerBlocks3s = append(probLayerBlocks3s, make([]float64, constant.StdLen, constant.StdLen))

	for bidx := 1; bidx < len(blocks); bidx++ { //traverse general blocks.
		//p(layer|block0)
		//bidx = 107
		ceil, floor := blocks[bidx-1], blocks[bidx]
		var probBlock = utils.StatProbBlockWithWeight(*nearDrills, ceil, floor)
		//var probBlock = utils.StatProbBlock(*nearDrills, ceil, floor)
		var probLayers = make([]float64, constant.StdLen, constant.StdLen)
		var probBlockLayers = make([]float64, constant.StdLen, constant.StdLen)
		var probLayerBlock2s = make([]float64, constant.StdLen, constant.StdLen)

		for layerIdx := int(1); layerIdx < constant.StdLen; layerIdx++ { //layer[0] is ground.
			//layerIdx = 26
			probLayers[layerIdx] = utils.StatProbLayerWithWeight(*nearDrills, ceil, floor, layerIdx)
			//probLayers[layerIdx] = utils.StatProbLayer(*nearDrills, ceil, floor, layerIdx)
			probBlockLayers[layerIdx] = utils.StatProbBlockLayer(*drillSet, ceil, floor, layerIdx)

			if probBlock >= 0.0000001 {
				probLayerBlock2s[layerIdx] = probBlockLayers[layerIdx] * probLayers[layerIdx] / probBlock
			}
		}
		probLayerBlocks3s = append(probLayerBlocks3s, probLayerBlock2s)

	}

	for idx := 1; idx < len(virtualDrill.LayerHeights); idx++ {
		ceil, floor := virtualDrill.LayerHeights[idx-1], virtualDrill.LayerHeights[idx]
		bidx := blocksIndex(blocks, ceil, floor)
		probs := probLayerBlocks3s[bidx]
		layer, prob := utils.FindMaxFloat64s(probs)
		virtualDrill.Layers = append(virtualDrill.Layers, int(layer))
		utils.Hole(prob)
	}
	virtualDrill.Merge()
	return
}
func GenVDrillM2(drillSet *[]entity.Drill, x, y float64) (vdrill entity.Drill) {
	log.SetFlags(log.Lshortfile)
	vdrill = vdrill.MakeDrill(constant.GenVDrillName(), x, y, 0)
	nearDrills := obtainNearDrills(drillSet, vdrill, constant.RadiusIn)
	for _, d := range *nearDrills {
		if math.Abs(x-d.X) < 0.001 && math.Abs(y-d.Y) < 0.001 {
			return d
		}
	}
	setClassicalIdwWeights(vdrill, *nearDrills)
	setLengthAndZ(&vdrill, *nearDrills)
	blocks := MakeBlocks(*drillSet, constant.BlockResZ)
	vdrill.LayerHeights = explodedHeights(blocks, vdrill.Z, vdrill.GetBottomHeight())
	vdBlocks := vdrill.LayerHeights
	//virtual(name, x, y, z, length, heights, weight)  还差 layers,

	//p(layer|block)
	var probBW *mat.Dense = stat.ProbBlocksW(nearDrills, &vdBlocks)
	var probLW *mat.Dense = stat.ProbLayersW(nearDrills)
	var probBLs *mat.Dense = stat.ProbBLs(drillSet, &blocks)
	var probLBs mat.Dense

	probBW.Apply(func(i, j int, val float64) float64 {
		return 1 / val
	}, probBW)
	probLBs.Mul(probBW, probLW)

	for idx := 1; idx < len(vdrill.LayerHeights); idx++ {
		ceil, floor := vdrill.LayerHeights[idx-1], vdrill.LayerHeights[idx]
		bidx := blocksIndex(blocks, ceil, floor)
		if bidx == -1 {
			log.Printf("drill name: %s, idx:%d\n", vdrill.Name, idx)
			log.Printf("ceil:%f, floor:%f\n", ceil, floor)
			log.Println(vdrill.LayerHeights)
			log.Fatalln(blocks)
		}
		s1 := probLBs.RawRowView(idx)
		s2 := probBLs.RawRowView(bidx)
		dst := make([]float64, len(s1))
		floats.AddTo(dst, s1, s2)
		probLBs.SetRow(idx, dst)
		layer, prob := utils.FindMaxFloat64s(dst)
		vdrill.Layers = append(vdrill.Layers, int(layer))
		utils.Hole(prob)
	}
	vdrill.Merge()
	return
}
func GenVDrillIDW(drillSet *[]entity.Drill, x, y float64) (vdrill entity.Drill) {
	log.SetFlags(log.Lshortfile)
	vdrill = vdrill.MakeDrill(constant.GenVDrillName(), x, y, 0)
	nearDrills := obtainNearDrills(drillSet, vdrill, constant.RadiusIn)
	for _, d := range *nearDrills { // if the position of the vdrill is just at a real drill's position
		if math.Abs(x-d.X) < 0.001 && math.Abs(y-d.Y) < 0.001 {
			return d
		}
	}
	setClassicalIdwWeights(vdrill, *nearDrills)
	nearDrills = utils.UnifyDrillsStrata(nearDrills, utils.CheckSeqMinNeg)
	vdrill.Layers = (*nearDrills)[0].Layers
	var vHeights []float64 = make([]float64, len(vdrill.Layers), len(vdrill.Layers))
	for lidx, _ := range vdrill.Layers {
		for _, d := range *nearDrills {
			vHeights[lidx] += d.GetWeight() * d.LayerHeights[lidx]
		}
		if lidx-1 >= 0 && math.Abs(vHeights[lidx]-vHeights[lidx-1]) < 10e-5 {
			vHeights[lidx] = vHeights[lidx-1]
		}
	}
	vdrill.LayerHeights = vHeights
	vdrill.Z = vHeights[0]
	vdrill.SetLength(vdrill.GetLength())
	vdrill.UnStdSeq()
	if !vdrill.IsValid() {
		vdrill.Print()
		log.Fatal("invalid vdirll.\n")
	}
	return vdrill
}
func setClassicalIdwWeights(center entity.Drill, aroundDrills []entity.Drill) (weights []float64) {
	var (
		weightSum       float64
		hasZeroDistance bool
		zeroIdx         int
	)

	//get distance
	for idx, aroundDrill := range aroundDrills {
		dist := center.DistanceBetween(aroundDrill)
		weights = append(weights, dist) //as distance
		if dist < 0.0001 {
			hasZeroDistance = true
			zeroIdx = idx
		}
	}
	if hasZeroDistance {
		for idx, _ := range weights {
			weights[idx] = 0
		}
		if weights != nil && zeroIdx >= 0 && zeroIdx < len(weights) {
			weights[zeroIdx] = 1
		}
	} else {
		for idx, _ := range weights { //cal weight
			weights[idx] = 1 / math.Pow(weights[idx], constant.IdwPow)
			weightSum += weights[idx]
		}
		for idx, _ := range weights { //归一化, and set int the drill.
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
func obtainNearDrills(drillSet *[]entity.Drill, drill entity.Drill, includeNum int) *[]entity.Drill {
	var drills []entity.Drill
	if includeNum > len(*drillSet) {
		includeNum = len(*drillSet)
		return drillSet
	}
	dists := make([]float64, len(*drillSet), len(*drillSet))
	for i, d := range *drillSet {
		dists[i] = drill.DistanceBetween(d)
	}

	sort.Float64s(dists)
	radius := dists[includeNum-1]

	for _, d := range *drillSet {
		if distance := drill.DistanceBetween(d); distance <= radius && d.Name != drill.Name {
			drills = append(drills, d)
		}
	}
	return &drills
}
func heightRange(drills []entity.Drill) (ceil float64, floor float64) {
	ceil, floor = -math.MaxFloat64, math.MaxFloat64
	for _, d := range drills {
		if d.Z > ceil {
			ceil = d.Z
		}
		if d.LayerHeights[len(d.LayerHeights)-1] < floor {
			floor = d.LayerHeights[len(d.LayerHeights)-1]
		}
	}
	return ceil, floor
}
func MakeBlocks(drillSet []entity.Drill, res float64) (blocksHeight []float64) {
	drillsCeil, drillsFloor := -math.MaxFloat64, math.MaxFloat64
	for _, d := range drillSet {
		if d.Z > drillsCeil {
			drillsCeil = d.Z
		}
		if d.LayerHeights[len(d.LayerHeights)-1] < drillsFloor {
			drillsFloor = d.LayerHeights[len(d.LayerHeights)-1]
		}
	}

	blocksHeight = append(blocksHeight, drillsCeil)

	for drillsCeil-res > drillsFloor {
		blocksHeight = append(blocksHeight, drillsCeil-res)
		drillsCeil = drillsCeil - res
	}

	//the last block may be un-standard block length, whose length may less than res
	blocksHeight = append(blocksHeight, drillsFloor)

	return
}
func getDrillsRecXOY(drills []entity.Drill) (l, r, t, b float64) {
	log.SetFlags(log.Lshortfile)
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
	//log.Println("rec t,b,l,r is: ", t, b, l, r)
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
func blocksIndex(blocks []float64, ceil, floor float64) (index int) {
	log.SetFlags(log.Lshortfile)
	for idx := 1; idx < len(blocks); idx++ {
		if ceil <= blocks[idx-1] && floor >= blocks[idx] {
			return idx
		}
	}
	debug.PrintStack()
	log.Fatal("error")
	return -1
}

func GetVirtualDrillsBetween(drill1, drill2 entity.Drill, n int,
	gen func(*[]entity.Drill, float64, float64) entity.Drill) (
	virtualDrills []entity.Drill) {
	log.SetFlags(log.Lshortfile)
	drillSet := constant.DrillSet()
	x1, y1 := drill1.X, drill1.Y
	x2, y2 := drill2.X, drill2.Y
	vertices := utils.SplitSegment(x1, y1, x2, y2, n)
	for idx := 1; idx < len(vertices); idx += 2 {
		virtualDrills = append(virtualDrills, gen(&drillSet, vertices[idx-1], vertices[idx]))
	}
	virtualDrills = append(virtualDrills, drill2)
	return
}
