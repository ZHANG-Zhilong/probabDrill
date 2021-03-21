package service

import (
	"fmt"
	"github.com/fogleman/poissondisc"
	"gonum.org/v1/gonum/mat"
	"log"
	"math"
	"probabDrill"
	"probabDrill/internal/constant"
	"probabDrill/internal/entity"
	"probabDrill/internal/utils"
	"runtime/debug"
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
				virtualDrills = append(virtualDrills, GenVDrillFromRDrillsM1(drillSet, nil, nil, x, y))
			} else {
				out++
			}
		}
	}
	log.Println("drillIn:", in, " drillOut:", out)
	for idx, _ := range bx {
		x, y := bx[idx], by[idx]
		virtualDrills = append(virtualDrills, GenVDrillFromRDrillsM1(drillSet, nil, nil, x, y))
	}
	return
}
func GenVDrillFromRDrillsM1(rDrills []entity.Drill, b []float64, p *mat.Dense, x, y float64) (vdrill entity.Drill) {
	rDrills = constant.GetRealDrills()
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	vdrill = vdrill.MakeDrill(constant.GenVDrillName(), x, y, 0)
	nearDrills := vdrill.NearKDrills(rDrills, probabDrill.RadiusIn)

	for _, d := range nearDrills {
		if math.Abs(x-d.X) < probabDrill.MinDrillDist && math.Abs(y-d.Y) < probabDrill.MinDrillDist {
			return d
		}
	}
	utils.SetClassicalIdwWeights(vdrill, nearDrills)
	entity.SetLengthAndZ(&vdrill, nearDrills)
	blocks := constant.GetBlocksR()
	vdrill.LayerHeights = utils.InterceptBlocks(blocks, vdrill.Z, vdrill.BottomHeight())
	//virtual(name, x, y, z, length, heights, weight)  还差 layers,

	//p(layer|block)
	pBlocksWMat, _ := utils.ProbBlocksWMat(nearDrills, blocks)
	pLayersWithWeightMat, _ := utils.ProbLayerWMat(nearDrills, blocks)
	//pDrillSetBlockLayerMat, _ := utils.ProbBlockLayerMatG(rDrills, blocks)
	pDrillSetBlockLayerMat := constant.GetBlockLayerRDrillMat()
	var pLayerBlockMat = mat.NewDense(len(blocks), probabDrill.StdLen, nil)
	for bidx := 1; bidx < len(blocks); bidx++ { //traverse general blocks.
		//p(layer|block0)
		ceil, floor := blocks[bidx-1], blocks[bidx]
		if ceil <= floor {
			debug.PrintStack()
			log.Println(blocks)
			log.Printf("len(blocks)%d, bidx:%d, ceil:%f, floor %f\n", len(blocks), bidx, ceil, floor)
			log.Fatal("error")
		}
		pBlock := pBlocksWMat.At(bidx, 0)
		for lidx := 1; lidx < probabDrill.StdLen; lidx++ { //layer[0] is ground.
			pBlockLayers := pDrillSetBlockLayerMat.At(bidx, lidx)
			pLayerWithWeight := pLayersWithWeightMat.At(bidx, lidx)
			if pBlock >= 1e-7 {
				pLayerBlock := pBlockLayers * pLayerWithWeight / pBlock
				pLayerBlockMat.Set(bidx, lidx, pLayerBlock)
			}
		}
	}

	for idx := 1; idx < len(vdrill.LayerHeights); idx++ {
		ceil, floor := vdrill.LayerHeights[idx-1], vdrill.LayerHeights[idx]
		bidx := constant.BlockIndexR(ceil, floor)
		probs := pLayerBlockMat.RawRowView(bidx)
		layer, _ := utils.FindMaxFloat64s(probs)
		vdrill.Layers = append(vdrill.Layers, int(layer))
	}
	vdrill.Merge()
	return
}

func GenVDrillM1(drills []entity.Drill, blocks []float64, pBlockLayerMat *mat.Dense, x, y float64) (vdrill entity.Drill) {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	vdrill = vdrill.MakeDrill(constant.GenVDrillName(), x, y, 0)
	nearDrills := vdrill.NearKDrills(drills, probabDrill.RadiusIn)
	for _, d := range nearDrills {
		if math.Abs(x-d.X) < probabDrill.MinDrillDist && math.Abs(y-d.Y) < probabDrill.MinDrillDist {
			d.Name = "real"
			return d
		}
	}
	utils.SetClassicalIdwWeights(vdrill, nearDrills)
	entity.SetLengthAndZ(&vdrill, nearDrills)
	heights := utils.InterceptBlocks(blocks, vdrill.Z, vdrill.Z-vdrill.GetLength())

	pBlocksWMat, _ := utils.ProbBlocksWMat(nearDrills, blocks)
	pLayersWithWeightMat, _ := utils.ProbLayerWMat(nearDrills, blocks)
	var pLayerBlockMat = mat.NewDense(len(blocks), probabDrill.StdLen, nil)
	for bidx := 1; bidx < len(blocks); bidx++ { //traverse general blocks.
		pBlock := pBlocksWMat.At(bidx, 0)
		for lidx := 1; lidx < probabDrill.StdLen; lidx++ { //layer[0] is ground.
			pBlockLayers := pBlockLayerMat.At(bidx, lidx)
			pLayerWithWeight := pLayersWithWeightMat.At(bidx, lidx)
			if pBlock >= 1e-7 {
				pLayerBlock := pBlockLayers * pLayerWithWeight / pBlock
				pLayerBlockMat.Set(bidx, lidx, pLayerBlock)
			}
		}
	}
	var bidx int
	for idx := 1; idx < len(heights); idx++ {
		bidx, _ = utils.BlocksIndex(blocks, heights[idx-1], heights[idx])
		if bidx == -1 {
			fmt.Println(blocks, heights[idx-1], heights[idx])
			fmt.Println(heights)
		}
		probs := pLayerBlockMat.RawRowView(bidx)
		layer, _ := utils.FindMaxFloat64s(probs)
		if err := vdrill.AddLayerWithHeight(layer, heights[idx]); err != nil {
			log.Println("add layer failed.")
		}
	}
	//add other layers? after interpolate the bottom of the virtual drill?
	//for {
	//	bidx++
	//	if bidx >= len(blocks) {
	//		break
	//	}
	//	probs := pLayerBlockMat.RawRowView(bidx)
	//	layer, _ := utils.FindMaxFloat64s(probs)
	//	if layer == 0 {
	//		break
	//	}
	//	vdrill.AddLayerWithHeight(layer, blocks[bidx])
	//}
	vdrill.GetLength() //modify drill.
	vdrill.Merge()
	vdrill.UnBlock()
	return
}
func GenVDrillsM1(drills []entity.Drill, points []float64) (vdrills []entity.Drill) {
	blocks := utils.MakeBlocks(drills, probabDrill.BlockResZ)
	pBlockLayerMat, _ := utils.ProbBlockLayerMatG(drills, blocks)
	for idx := 1; idx < len(points); idx += 2 {
		vdrills = append(vdrills, GenVDrillM1(drills, blocks, pBlockLayerMat, points[idx-1], points[idx]))
	}
	return vdrills
}
func GenVDrillFromHelpDrillsM1(helpDrills []entity.Drill, b []float64, p *mat.Dense, x, y float64) (vdrill entity.Drill) {
	helpDrills = constant.GetHelpDrills()
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	vdrill = vdrill.MakeDrill(constant.GenVDrillName(), x, y, 0)
	nearDrills := vdrill.NearKDrills(helpDrills, probabDrill.RadiusIn)

	//if vdrill's position has already exist a drill, return directly
	for _, d := range nearDrills {
		if math.Abs(x-d.X) < probabDrill.MinDrillDist && math.Abs(y-d.Y) < probabDrill.MinDrillDist {
			return d
		}
	}

	//set idw weight for incident drills
	utils.SetClassicalIdwWeights(vdrill, nearDrills)
	entity.SetLengthAndZ(&vdrill, nearDrills)

	blocks := constant.GetBlocksH()
	vdrill.LayerHeights = utils.InterceptBlocks(blocks, vdrill.Z, vdrill.BottomHeight())
	//virtual(name, x, y, z, length, heights, weight)  还差 layers,

	//p(layer|block)
	pDrillSetBlockLayerMat := constant.GetBlockLayerHDrillMat()
	pBlocksWMat, _ := utils.ProbBlocksWMat(nearDrills, blocks)
	pLayersWithWeightMat, _ := utils.ProbLayerWMat(nearDrills, blocks)
	var pLayerBlockMat = mat.NewDense(len(blocks), probabDrill.StdLen, nil)
	//p(layer|block0)
	for bidx := 1; bidx < len(blocks); bidx++ { //traverse general blocks.
		pBlock := pBlocksWMat.At(bidx, 0)
		for lidx := 1; lidx < probabDrill.StdLen; lidx++ { //layer[0] is ground.
			pBlockLayers := pDrillSetBlockLayerMat.At(bidx, lidx)
			pLayerW := pLayersWithWeightMat.At(bidx, lidx)
			if pBlock >= 1e-7 {
				pLayerBlock := pBlockLayers * pLayerW / pBlock
				pLayerBlockMat.Set(bidx, lidx, pLayerBlock)
			}
		}
	}

	for idx := 1; idx < len(vdrill.LayerHeights); idx++ {
		ceil, floor := vdrill.LayerHeights[idx-1], vdrill.LayerHeights[idx]
		bidx := constant.BlockIndexH(ceil, floor)
		probs := pLayerBlockMat.RawRowView(bidx)
		layer, _ := utils.FindMaxFloat64s(probs)
		vdrill.Layers = append(vdrill.Layers, int(layer))
	}
	vdrill.Merge()
	vdrill.UnBlock()
	return
}

func GenVDrillIDW(drillSet []entity.Drill, blocks []float64, p *mat.Dense, x, y float64) (vdrill entity.Drill) {
	log.SetFlags(log.Lshortfile)
	vdrill = vdrill.MakeDrill(constant.GenVDrillName(), x, y, 0)
	nearDrills := vdrill.NearKDrills(drillSet, probabDrill.RadiusIn)

	for _, d := range nearDrills { // if the position of the vdrill is just at a real drill's position
		if math.Abs(x-d.X) < 0.001 && math.Abs(y-d.Y) < 0.001 {
			return d
		}
	}
	utils.SetClassicalIdwWeights(vdrill, nearDrills)
	nearDrills = constant.UnifyDrillsSeq(nearDrills, constant.CheckSeqMinNeg)
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
	vdrill.UnBlock()
	if !vdrill.IsValid() {
		vdrill.Display()
		log.Fatal("invalid vdrill")
	}
	return vdrill
}
func GenVDrillsIDW(drills []entity.Drill, points []float64) (vdrills []entity.Drill) {
	for idx := 1; idx < len(points); idx += 2 {
		vdrills = append(vdrills, GenVDrillIDW(drills, nil, nil, points[idx-1], points[idx]))
	}
	return vdrills
}
func GenHelpDrills() (hdrills []entity.Drill) {
	drillSet := constant.GetRealDrills()
	var x0, y0, x1, y1, r float64
	r = 800 // min distance between points
	k := 10 // max attempts to add neighboring point
	x0, y0, x1, y1 = drillSet[0].GetRec(drillSet)
	points := poissondisc.Sample(x0, y0, x1, y1, r, k, nil)
	for _, p := range points {
		hdrills = append(hdrills, GenVDrillIDW(drillSet, nil, nil, p.X, p.Y))
	}
	return
}

type GenVDrills func([]entity.Drill, []float64) []entity.Drill

func GenVDrillsBetween(drillSet []entity.Drill, drill1, drill2 entity.Drill, n int, gen GenVDrills) (vDrills []entity.Drill) {
	log.SetFlags(log.Lshortfile)
	vertices := utils.MiddleKPoints(drill1.X, drill1.Y, drill2.X, drill2.Y, n)
	vDrills = append(vDrills, gen(drillSet, vertices)...)
	return
}
