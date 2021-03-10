package service

import (
	"github.com/fogleman/poissondisc"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"
	"log"
	"math"
	"probabDrill"
	"probabDrill/internal/constant"
	"probabDrill/internal/entity"
	"probabDrill/internal/stat"
	"probabDrill/internal/utils"
)

func GetGridDrills(drillSet []entity.Drill) (virtualDrills []entity.Drill) {

	px, py := probabDrill.GridXY, probabDrill.GridXY
	l, r, t, b := drillSet[0].GetRec(drillSet)

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
				virtualDrills = append(virtualDrills, GenVDrillM2(drillSet, x, y))
			} else {
				out++
			}
		}
	}
	log.Println("drillIn:", in, " drillOut:", out)
	for idx, _ := range bx {
		x, y := bx[idx], by[idx]
		virtualDrills = append(virtualDrills, GenVDrillM1(drillSet, x, y))
	}
	return
}
func GenVDrillM1(drillSet []entity.Drill, x, y float64) (virtualDrill entity.Drill) {
	log.SetFlags(log.Lshortfile)
	virtualDrill = virtualDrill.MakeDrill(constant.GenVDrillName(), x, y, 0)
	nearDrills := virtualDrill.NearDrills(drillSet, probabDrill.RadiusIn)
	for _, d := range nearDrills {
		if math.Abs(x-d.X) < 0.001 && math.Abs(y-d.Y) < 0.001 {
			return d
		}
	}
	utils.SetClassicalIdwWeights(virtualDrill, nearDrills)
	entity.SetLengthAndZ(&virtualDrill, nearDrills)
	blocks := utils.MakeBlocks(drillSet, 0.02)
	virtualDrill.LayerHeights = utils.ExplodedHeights(blocks, virtualDrill.Z, virtualDrill.GetBottomHeight())
	//virtual(name, x, y, z, length, heights, weight)  还差 layers,

	//p(layer|block)
	var probLayerBlocks3s [][]float64
	probLayerBlocks3s = append(probLayerBlocks3s, make([]float64, probabDrill.StdLen, probabDrill.StdLen))

	for bidx := 1; bidx < len(blocks); bidx++ { //traverse general blocks.
		//p(layer|block0)
		//bidx = 107
		ceil, floor := blocks[bidx-1], blocks[bidx]
		var probBlock = utils.StatProbBlockWithWeight(nearDrills, ceil, floor)
		//var probBlock = utils.StatProbBlock(*NearDrills, ceil, floor)
		var probLayers = make([]float64, probabDrill.StdLen, probabDrill.StdLen)
		var probBlockLayers = make([]float64, probabDrill.StdLen, probabDrill.StdLen)
		var probLayerBlock2s = make([]float64, probabDrill.StdLen, probabDrill.StdLen)

		for layerIdx := int(1); layerIdx < probabDrill.StdLen; layerIdx++ { //layer[0] is ground.
			//layerIdx = 26
			probLayers[layerIdx] = utils.StatProbLayerWithWeight(nearDrills, ceil, floor, layerIdx)
			//probLayers[layerIdx] = utils.StatProbLayer(*NearDrills, ceil, floor, layerIdx)
			probBlockLayers[layerIdx] = utils.StatProbBlockLayer(drillSet, ceil, floor, layerIdx)

			if probBlock >= 0.0000001 {
				probLayerBlock2s[layerIdx] = probBlockLayers[layerIdx] * probLayers[layerIdx] / probBlock
			}
		}
		probLayerBlocks3s = append(probLayerBlocks3s, probLayerBlock2s)

	}

	for idx := 1; idx < len(virtualDrill.LayerHeights); idx++ {
		ceil, floor := virtualDrill.LayerHeights[idx-1], virtualDrill.LayerHeights[idx]
		bidx := utils.BlocksIndex(blocks, ceil, floor)
		probs := probLayerBlocks3s[bidx]
		layer, prob := utils.FindMaxFloat64s(probs)
		virtualDrill.Layers = append(virtualDrill.Layers, int(layer))
		utils.Hole(prob)
	}
	virtualDrill.Merge()
	return
}
func GenVDrillM2(drillSet []entity.Drill, x, y float64) (vdrill entity.Drill) {
	log.SetFlags(log.Lshortfile)
	vdrill = vdrill.MakeDrill(constant.GenVDrillName(), x, y, 0)
	nearDrills := vdrill.NearDrills(drillSet, probabDrill.RadiusIn)
	for _, d := range nearDrills {
		if math.Abs(x-d.X) < 0.001 && math.Abs(y-d.Y) < 0.001 {
			return d
		}
	}
	utils.SetClassicalIdwWeights(vdrill, nearDrills)
	entity.SetLengthAndZ(&vdrill, nearDrills)
	blocks := utils.MakeBlocks(drillSet, probabDrill.BlockResZ)
	vdrill.LayerHeights = utils.ExplodedHeights(blocks, vdrill.Z, vdrill.GetBottomHeight())
	vdBlocks := vdrill.LayerHeights
	//virtual(name, x, y, z, length, heights, weight)  还差 layers,

	//p(layer|block)
	var probBW *mat.Dense = stat.ProbBlocksW(nearDrills, vdBlocks)
	var probLW *mat.Dense = stat.ProbLayersW(nearDrills)
	var probBLs *mat.Dense = stat.ProbBLs(drillSet, blocks)
	var probLBs mat.Dense

	probBW.Apply(func(i, j int, val float64) float64 {
		return 1 / val
	}, probBW)
	probLBs.Mul(probBW, probLW)

	for idx := 1; idx < len(vdrill.LayerHeights); idx++ {
		ceil, floor := vdrill.LayerHeights[idx-1], vdrill.LayerHeights[idx]
		bidx := utils.BlocksIndex(blocks, ceil, floor)
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
func GenVDrillIDW(drillSet []entity.Drill, x, y float64) (vdrill entity.Drill) {
	log.SetFlags(log.Lshortfile)
	vdrill = vdrill.MakeDrill(constant.GenVDrillName(), x, y, 0)
	nearDrills := vdrill.NearDrills(drillSet, probabDrill.RadiusIn)
	//if vdrill.name == "virtual10" {
	//	for _, d := range nearDrills {
	//		fmt.Printf("%#v\n", d)
	//		fmt.Println(vdrill.Distance(d))
	//	}
	//	var dists []float64
	//	for _, d := range drillSet {
	//		dists = append(dists, vdrill.Distance(d))
	//	}
	//	sort.Slice(dists, func(i, j int) bool {
	//		return dists[i] < dists[j]
	//	})
	//	utils.PrintFloat64s(dists)
	//}
	for _, d := range nearDrills { // if the position of the vdrill is just at a real drill's position
		if math.Abs(x-d.X) < 0.001 && math.Abs(y-d.Y) < 0.001 {
			return d
		}
	}
	utils.SetClassicalIdwWeights(vdrill, nearDrills)
	nearDrills = utils.UnifyDrillsStrata(nearDrills, utils.CheckSeqMinNeg)
	vdrill.Layers = nearDrills[0].Layers
	var vHeights = make([]float64, len(vdrill.Layers), len(vdrill.Layers))
	for lidx, _ := range vdrill.Layers {
		for _, d := range nearDrills {
			vHeights[lidx] += utils.Decimal(d.GetWeight() * d.LayerHeights[lidx])
		}
		if lidx-1 >= 0 && math.Abs(vHeights[lidx]-vHeights[lidx-1]) < 10e-5 {
			vHeights[lidx] = vHeights[lidx-1]
		}
	}
	vdrill.LayerHeights = vHeights
	vdrill.Z = vHeights[0]
	vdrill.GetLength()
	vdrill.UnStdSeq()
	if !vdrill.IsValid() {
		vdrill.Print()
		log.Fatal("invalid vdirll.\n")
	}
	return vdrill
}
func GenHelpDrills() (hdrills []entity.Drill) {
	drillSet := constant.DrillSet()
	var x0, y0, x1, y1, r float64
	r = 400 // min distance between points
	k := 10 // max attempts to add neighboring point
	x0, y0, x1, y1 = drillSet[0].GetRec(drillSet)
	points := poissondisc.Sample(x0, y0, x1, y1, r, k, nil)
	for _, p := range points {
		hdrills = append(hdrills, GenVDrillIDW(drillSet, p.X, p.Y))
	}
	return
}

type GenVDrills func([]entity.Drill, float64, float64) entity.Drill

func GenVDrillsBetween(drill1, drill2 entity.Drill, n int, gen GenVDrills) (vDrills []entity.Drill) {
	log.SetFlags(log.Lshortfile)
	//drillSet := constant.DrillSet()
	drillSet := GenHelpDrills()
	vertices := utils.SplitSegment(drill1.X, drill1.Y, drill2.X, drill2.Y, n)
	for idx := 1; idx < len(vertices); idx += 2 {
		vDrills = append(vDrills, gen(drillSet, vertices[idx-1], vertices[idx]))
	}
	vDrills = append(vDrills, drill2)
	return
}
