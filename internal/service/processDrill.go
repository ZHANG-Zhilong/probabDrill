package service

import (
	"log"
	"math"
	"probabDrill/internal/constant"
	"probabDrill/internal/entity"
	"probabDrill/internal/utils"
	"sort"
)

func GetGridDrills(drills []entity.Drill) (virtualDrills []entity.Drill) {

	px, py := constant.GridXY, constant.GridXY
	l, r, t, b := getDrillsRecXOY(drills)

	//grid to interpolate
	gridx, gridy := utils.GetGrids(px, py, l, r, t, b)
	log.Println(gridx)
	log.Println(gridy)
	blocks := makeBlocks(drills, constant.BlockResZ)
	bx, by := constant.GetBoundary()
	var in, out int
	for idx := range gridx {
		for idy := range gridy {
			x := gridx[idx]
			y := gridy[idy]
			if utils.IsInPolygon(bx, by, x, y) {
				in++

				virtualDrills = append(virtualDrills, generateVirtualDrill(x, y, blocks))
			} else {
				out++
			}
		}
	}
	log.Println("drillIn:", in, " drillOut:", out)
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
	log.Println(virtualDrill.Name)
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
		var probBlockWithWeight = utils.StatProbBlockWithWeight(incidentDrills, ceil, floor)
		var probLayersWithWeight = make([]float64, constant.StdLen, constant.StdLen)
		var probBlockLayers = make([]float64, constant.StdLen, constant.StdLen)
		var probLayerBlock2s = make([]float64, constant.StdLen, constant.StdLen)

		for layerIdx := int64(1); layerIdx < constant.StdLen; layerIdx++ { //layer[0] is ground.
			//layerIdx = 26
			probLayersWithWeight[layerIdx] = utils.StatProbLayerWithWeight(incidentDrills, ceil, floor, layerIdx)
			probBlockLayers[layerIdx] = utils.StatProbBlockLayer(constant.DrillSet(), ceil, floor, layerIdx)

			if probBlockWithWeight >= 0.0000001 {
				probLayerBlock2s[layerIdx] = probBlockLayers[layerIdx] * probLayersWithWeight[layerIdx] / probBlockWithWeight
			}
			a, b, c := probLayersWithWeight[layerIdx], probBlockLayers[layerIdx], probLayerBlock2s[layerIdx]
			utils.Hole(a, b, c)
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
		layer, prob := utils.FindMaxFloat64s(probs)
		virtualDrill.Layers = append(virtualDrill.Layers, int64(layer))
		utils.Hole(prob)
		//log.Println(bidx, ceil, floor, layer, prob)
		//log.Println(probs)
	}

	//log.Println("before merged.")
	//virtualDrill.Print()
	virtualDrill.Merge()
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
func setClassicalIdwWeights(center entity.Drill, aroundDrills []entity.Drill) (weights []float64) {
	var (
		weightSum       float64
		hasZeroDistance bool
		zeroIdx         int
	)

	//get distance
	for idx, aroundDrill := range aroundDrills {
		dist := center.DistanceBetween(aroundDrill)
		weights = append(weights, dist)
		if dist < 0.0001 {
			hasZeroDistance = true
			zeroIdx = idx
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
		dists[i] = drill.DistanceBetween(d)
	}

	sort.Float64s(dists)
	radius := dists[includeNum-1]

	for _, d := range drillSet {
		if distance := drill.DistanceBetween(d); distance <= radius && d.Name != drill.Name {
			drills = append(drills, d)
		}
	}
	return drills
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
func makeBlocks(drillSet []entity.Drill, resz float64) (blocksHeight []float64) {
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
	return -1
}