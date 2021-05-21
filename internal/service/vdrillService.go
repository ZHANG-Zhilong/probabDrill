package service

import (
	"fmt"
	"github.com/fogleman/poissondisc"
	"github.com/spf13/viper"
	"gonum.org/v1/gonum/mat"
	"log"
	"math"
	"math/rand"
	"probabDrill/apps/probDrill/model"
	"probabDrill/internal/constant"
	"probabDrill/internal/utils"
	"runtime/debug"
	"time"
)

func GetGridDrills(drillSet []model.Drill) (virtualDrills []model.Drill) {

	px, py := viper.GetFloat64("GridXY"), viper.GetFloat64("GridXY")
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
func GenVDrillFromRDrillsM1(rDrills []model.Drill, b []float64, p *mat.Dense, x, y float64) (vdrill model.Drill) {
	rDrills = constant.GetRealDrills()
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	vdrill = vdrill.MakeDrill(constant.GenVDrillName(), x, y, 0)
	nearDrills := vdrill.NearKDrills(rDrills, viper.GetInt("RadiusIn"))

	for _, d := range nearDrills {
		if math.Abs(x-d.X) < viper.GetFloat64("MinDrillDist") && math.Abs(y-d.Y) < viper.GetFloat64("MinDrillDist") {
			return d
		}
	}
	utils.SetClassicalIdwWeights(vdrill, nearDrills)
	model.SetLengthAndZ(&vdrill, nearDrills)
	blocks := constant.GetBlocksR()
	vdrill.LayerHeights = utils.InterceptBlocks(blocks, vdrill.Z, vdrill.BottomHeight())
	//virtual(name, x, y, z, length, heights, weight)  还差 layers,

	//p(layer|block)
	pBlocksWMat, _ := utils.ProbBlocksWMat(nearDrills, blocks)
	pLayersWithWeightMat, _ := utils.ProbLayerWMat(nearDrills, blocks)
	//pDrillSetBlockLayerMat, _ := utils.ProbBlockLayerMatG(rDrills, blocks)
	pDrillSetBlockLayerMat := constant.GetBlockLayerRDrillMat()
	var pLayerBlockMat = mat.NewDense(len(blocks), viper.GetInt("StdLen"), nil)
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
		for lidx := 1; lidx < viper.GetInt("StdLen"); lidx++ { //layer[0] is ground.
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

func GenVDrillM1(drills []model.Drill, blocks []float64, pBlockLayerMat *mat.Dense, x, y float64) model.Drill {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	//get blocks split of the study area.
	vdrill := drills[0].MakeDrill(constant.GenVDrillName(), x, y, 0)

	//get the drills near virtual drill.(target drill), if the dist < threshold the drill will be removed.
	nearDrills := vdrill.NearKDrills(drills, viper.GetInt("RadiusIn"))

	//filter the drills too near.
	for _, d := range nearDrills {
		dist := d.Distance(vdrill)
		if math.Abs(dist) < viper.GetFloat64("MinDrillDist") {
			d.Name = "real"
			return d
		}

	}

	//set the idw weight for near drills.
	utils.SetClassicalIdwWeights(vdrill, nearDrills)

	//set the length and drill.Z for virtual drill (target drill).
	model.SetLengthAndZ(&vdrill, nearDrills)

	//get the local for the virtual drill from the study area's blocks. local blocks is the subset of the blocks.
	heights := utils.InterceptBlocks(blocks, vdrill.Z, vdrill.Z-vdrill.GetLength())

	//get the Vector of P(blocks) in local area
	pBlocksWMat, _ := utils.ProbBlocksWMat(nearDrills, blocks)

	//get the Matrix of P(layers_ij) for local area. P(layer_ij) in different local blocks.
	pLayersWithWeightMat, _ := utils.ProbLayerWMat(nearDrills, blocks)

	//the result. P(layer|blocks)  in local area. the param P(blocks|layers) is input param.
	var pLayerBlockMat = mat.NewDense(len(blocks), viper.GetInt("StdLen"), nil)

	for bidx := 1; bidx < len(blocks); bidx++ {
		pBlock := pBlocksWMat.At(bidx, 0)

		//traverse general blocks, to generate Vector P(layer_ij|block_i)
		for lidx := 1; lidx < viper.GetInt("StdLen"); lidx++ { //layer[0] is ground.
			pBlockLayers := pBlockLayerMat.At(bidx, lidx)
			pLayerWithWeight := pLayersWithWeightMat.At(bidx, lidx)
			//ensure the validation of the rst.
			if pBlock >= 1e-7 {
				pLayerBlock := pBlockLayers * pLayerWithWeight / pBlock
				pLayerBlockMat.Set(bidx, lidx, pLayerBlock)
			}
		}
	}

	//print the P(layer_ij|block_i) formatted.
	//fa := mat.Formatted(pLayerBlockMat, mat.Prefix(""), mat.Squeeze())
	//fmt.Printf("with all values:\na = %v\n\n", fa)

	//generate virtual drill.
	var bidx int
	for idx := 1; idx < len(heights); idx++ {
		bidx, _ = utils.BlocksIndex(blocks, heights[idx-1], heights[idx])
		if bidx == -1 {
			_, err := utils.BlocksIndex(blocks, heights[idx-1], heights[idx])
			if err != nil {
				log.Print(err)
			}
			log.Println(blocks)
			log.Println(heights[idx-1], heights[idx])
			log.Println(heights)
			break
		}
		//row, _ := pLayerBlockMat.Dims()
		probs := pLayerBlockMat.RawRowView(bidx)
		layer, _ := utils.FindMaxFloat64s(probs)
		//if layer == 0 {
		//	log.Printf("err in p(layer|block) probs, %#v\n", probs)
		//	fa := mat.Formatted(pLayerBlockMat, mat.Prefix(""), mat.Squeeze())
		//	fmt.Printf("P(layer|block) with all values:\n%v\n\n", fa)
		//	os.Exit(-1)
		//}
		if err := vdrill.AddLayerWithHeight(layer, heights[idx]); err != nil {
			log.Println(":AddLayerWithHeight failed", err)
		}
	}
	//not die zhi. use die zhi principle when draw the slice figure, keep the origin of the interpolated vdrill data.

	//update drill data.
	vdrill.GetLength()
	vdrill.Merge()
	vdrill.UnBlock()
	return vdrill
}
func GenVDrillM1Second(drills []model.Drill, blocks []float64, pBlockLayerMat *mat.Dense, x, y float64) (vdrill model.Drill) {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	//get blocks split of the study area.
	vdrill = vdrill.MakeDrill(constant.GenVDrillName(), x, y, 0)

	//get the drills near virtual drill.(target drill)
	nearDrills := vdrill.NearKDrills(drills, viper.GetInt("RadiusIn"))

	//filter the drills too near.
	for _, d := range nearDrills {
		if math.Abs(x-d.X) < viper.GetFloat64("MinDrillDist") && math.Abs(y-d.Y) < viper.GetFloat64("MinDrillDist") {
			d.Name = "real"
			return d
		}
	}

	//set the idw weight for near drills.
	utils.SetClassicalIdwWeights(vdrill, nearDrills)

	//set the length and drill.Z for virtual drill (target drill).
	model.SetLengthAndZ(&vdrill, nearDrills)

	//get the local for the virtual drill from the study area's blocks. local blocks is the subset of the blocks.
	heights := utils.InterceptBlocks(blocks, vdrill.Z, vdrill.Z-vdrill.GetLength())

	//get the Vector of P(blocks) in local area
	pBlocksWMat, _ := utils.ProbBlocksWMat(nearDrills, blocks)

	//get the Matrix of P(layers_ij) for local area. P(layer_ij) in different local blocks.
	pLayersWithWeightMat, _ := utils.ProbLayerWMat(nearDrills, blocks)

	//the result. P(layer|blocks)  in local area. the param P(blocks|layers) is input param.
	var pLayerBlockMat = mat.NewDense(len(blocks), viper.GetInt("StdLen"), nil)

	for bidx := 1; bidx < len(blocks); bidx++ {
		pBlock := pBlocksWMat.At(bidx, 0)

		//traverse general blocks, to generate Vector P(layer_ij|block_i)
		for lidx := 1; lidx < viper.GetInt("StdLen"); lidx++ { //layer[0] is ground.
			pBlockLayers := pBlockLayerMat.At(bidx, lidx)
			pLayerWithWeight := pLayersWithWeightMat.At(bidx, lidx)
			//ensure the validation of the rst.
			if pBlock >= 1e-7 {
				pLayerBlock := pBlockLayers * pLayerWithWeight / pBlock
				pLayerBlockMat.Set(bidx, lidx, pLayerBlock)
			}
		}
	}

	//print the P(layer_ij|block_i) formatted.
	//fa := mat.Formatted(pLayerBlockMat, mat.Prefix(""), mat.Squeeze())
	//fmt.Printf("with all values:\na = %v\n\n", fa)

	//generate virtual drill.
	var bidx int
	for idx := 1; idx < len(heights); idx++ {
		bidx, _ = utils.BlocksIndex(blocks, heights[idx-1], heights[idx])
		if bidx == -1 {
			fmt.Println(blocks, heights[idx-1], heights[idx])
			fmt.Println(heights)
		}
		probs := pLayerBlockMat.RawRowView(bidx)
		layer, _ := utils.FindSecondMaxFloat64s(probs)
		if err := vdrill.AddLayerWithHeight(layer, heights[idx]); err != nil {
			log.Println(":AddLayerWithHeight failed", err)
		}
	}

	//not die zhi. use die zhi principle when draw the slice figure, keep the origin of the interpolated vdrill data.

	//update drill data.
	vdrill.GetLength()
	vdrill.Merge()
	vdrill.UnBlock()
	return
}
func GenVDrillsM1(drills []model.Drill, points []float64) (vdrills []model.Drill) {
	blocks := utils.MakeBlocks(drills, viper.GetFloat64("BlockResZ"))
	pBlockLayerMat, _ := utils.ProbBlockLayerMatG(drills, blocks)
	for idx := 1; idx < len(points); idx += 2 {
		vdrills = append(vdrills, GenVDrillM1(drills, blocks, pBlockLayerMat, points[idx-1], points[idx]))
	}
	return vdrills
}
func GenVDrillsM1Second(drills []model.Drill, points []float64) (vdrills []model.Drill) {
	blocks := utils.MakeBlocks(drills, viper.GetFloat64("BlockResZ"))
	pBlockLayerMat, _ := utils.ProbBlockLayerMatG(drills, blocks)
	for idx := 1; idx < len(points); idx += 2 {
		vdrills = append(vdrills, GenVDrillM1Second(drills, blocks, pBlockLayerMat, points[idx-1], points[idx]))
	}
	return vdrills
}
func GenVDrillFromHelpDrillsM1(helpDrills []model.Drill, b []float64, p *mat.Dense, x, y float64) (vdrill model.Drill) {
	helpDrills = constant.GetHelpDrills()
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	vdrill = vdrill.MakeDrill(constant.GenVDrillName(), x, y, 0)
	nearDrills := vdrill.NearKDrills(helpDrills,viper.GetInt("RadiusIn"))

	//if vdrill's position has already exist a drill, return directly
	for _, d := range nearDrills {
		if math.Abs(x-d.X) < viper.GetFloat64("MinDrillDist") && math.Abs(y-d.Y) < viper.GetFloat64("MinDrillDist") {
			return d
		}
	}

	//set idw weight for incident drills
	utils.SetClassicalIdwWeights(vdrill, nearDrills)
	model.SetLengthAndZ(&vdrill, nearDrills)

	blocks := constant.GetBlocksH()
	vdrill.LayerHeights = utils.InterceptBlocks(blocks, vdrill.Z, vdrill.BottomHeight())
	//virtual(name, x, y, z, length, heights, weight)  还差 layers,

	//p(layer|block)
	pDrillSetBlockLayerMat := constant.GetBlockLayerHDrillMat()
	pBlocksWMat, _ := utils.ProbBlocksWMat(nearDrills, blocks)
	pLayersWithWeightMat, _ := utils.ProbLayerWMat(nearDrills, blocks)
	var pLayerBlockMat = mat.NewDense(len(blocks), viper.GetInt("StdLen"), nil)
	viper.GetInt("StdLen")
	//p(layer|block0)
	for bidx := 1; bidx < len(blocks); bidx++ { //traverse general blocks.
		pBlock := pBlocksWMat.At(bidx, 0)
		for lidx := 1; lidx < viper.GetInt("StdLen"); lidx++ { //layer[0] is ground.
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

func GenVDrillIDW(drillSet []model.Drill, blocks []float64, p *mat.Dense, x, y float64) (vdrill model.Drill) {
	log.SetFlags(log.Lshortfile)
	vdrill = vdrill.MakeDrill(constant.GenVDrillName(), x, y, 0)
	nearDrills := vdrill.NearKDrills(drillSet, viper.GetInt("RadiusIn"))

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
	vdrill.Merge()
	if !vdrill.IsValid() {
		vdrill.Display()
		log.Fatal("invalid vdrill")
	}
	return vdrill
}
func GenVDrillsIDW(drills []model.Drill, points []float64) (vdrills []model.Drill) {
	for idx := 1; idx < len(points); idx += 2 {
		vdrills = append(vdrills, GenVDrillIDW(drills, nil, nil, points[idx-1], points[idx]))
	}
	return vdrills
}
func GenHelpDrills(drillSet []model.Drill) (hdrills []model.Drill) {
	var x0, y0, x1, y1 float64
	x0, y0, x1, y1 = drillSet[0].GetRec(drillSet)
	points := poissondisc.Sample(x0, y0, x1, y1, viper.GetFloat64("MinDistance"), viper.GetInt("MaxAttemptAdd"), nil)
	for _, p := range points {
		hdrills = append(hdrills, GenVDrillIDW(drillSet, nil, nil, p.X, p.Y))
	}
	return
}

type GenVDrills func([]model.Drill, []float64) []model.Drill

func GenVDrillsBetween(drillSet []model.Drill, drill1, drill2 model.Drill, n int, gen GenVDrills) (vDrills []model.Drill) {
	log.SetFlags(log.Lshortfile)
	vertices := utils.MiddleKPoints(drill1.X, drill1.Y, drill2.X, drill2.Y, n)
	vDrills = append(vDrills, gen(drillSet, vertices)...)
	return
}

func GetPeByDrill(drill1, drill2 model.Drill) (meanPe float64) {
	//calculate the pe of the drills' different layer interface.
	layerPEs := make([]float64, viper.GetInt("StdLen"), viper.GetInt("StdLen"))
	for layer := 1; layer < viper.GetInt("StdLen"); layer++ {
		var hasLayer int
		var observeValues, estimateValues []float64
		h1, e1 := drill1.LayerBottomHeight(layer)
		h2, e2 := drill2.LayerBottomHeight(layer)
		if e1 == nil && e2 == nil {
			hasLayer++
			observeValues = append(observeValues, h1[0])
			estimateValues = append(estimateValues, h2[0])
		}
		if pe, err3 := utils.PercentageError(observeValues, estimateValues); err3 == nil {
			layerPEs[layer] = pe
		} else {
			layerPEs[layer] = -1
		}
	}

	//filter pes
	for idx, val := range layerPEs {
		if val == -1 {
			layerPEs[idx] = 0
		}
	}
	meanPe, _ = utils.GetMean(layerPEs)
	return
}
func GetCompareDrills(drillSet []model.Drill, ratio int, genVDrills GenVDrills) (real, compare []model.Drill, err error) {

	rand.Seed(time.Now().UnixNano())
	var drillSet1, drillSet2 []model.Drill

	//drillSet1 used to compare, drillSet2 used to as drill dta.
	//分离用于对比和插值的钻孔数据集
	for _, d := range drillSet {
		if rand.Intn(100) < ratio {
			drillSet1 = append(drillSet1, d)
		} else {
			drillSet2 = append(drillSet2, d)
		}
	}

	//获取用于比较的钻孔数据集的位置，将其放入一维数组
	var vertices []float64
	for _, d := range drillSet1 {
		vertices = append(vertices, d.X, d.Y)
	}

	//在比较钻孔数据集的原地位置处，生成相应的虚拟钻孔数据
	estimateDrills := genVDrills(drillSet2, vertices)
	if len(drillSet1) != len(estimateDrills) {
		return nil, nil, fmt.Errorf(
			":GetPeAvgByLayer,len(drillSet1) != len(estimateDrills):p1,%v,p2,%v ", drillSet1, estimateDrills)
	}
	return drillSet1, estimateDrills, nil
}
func GetPeByArea(realDrillSet []model.Drill, ratio int, genVDrills GenVDrills, round int) (x, y, pe []float64, err error) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	for idx := 0; idx < round; idx++ {
		compareDrills, estimateDrills, _ := GetCompareDrills(realDrillSet, ratio, genVDrills)
		for idx2 := 0; idx2 < len(compareDrills); idx2++ {
			x = append(x, compareDrills[idx2].X)
			y = append(y, compareDrills[idx2].Y)
			//p := GetPeByDrill(compareDrills[idx2], estimateDrills[idx2])
			//if p > 1 {
			//	log.Printf("the avg of pe between the two drill is: %f\n", p)
			//	log.Printf("comp drill: %#v\n", compareDrills[idx2])
			//	log.Printf("esti drill: %#v\n", estimateDrills[idx2])
			//}
			pe = append(pe, GetPeByDrill(compareDrills[idx2], estimateDrills[idx2]))
		}
		//过滤一下？把pe>1的pe都除以最大的pe？
	}
	return
}
func GetPeAvgByLayer(realDrillSet []model.Drill, ratio int, genVDrills GenVDrills) (layerPEs []float64, err error) {
	//什么情况下会返回-1？？？为何pe会大于1？
	compareDrills, estimateDrills, _ := GetCompareDrills(realDrillSet, ratio, genVDrills)
	//for idx, _ := range compareDrills {
	//	log.Printf("%#v", compareDrills[idx])
	//	log.Printf("%#v\n", estimateDrills[idx])
	//}
	//分别计算钻孔不同地层界面之间的百分比误差
	layerPEs = make([]float64, viper.GetInt("StdLen"), viper.GetInt("StdLen"))
	for layer := 1; layer < viper.GetInt("StdLen"); layer++ {
		var hasLayer int
		var observeValues, estimateValues []float64
		for idx := 0; idx < len(compareDrills); idx++ {
			h1, e1 := compareDrills[idx].LayerBottomHeight(layer)
			h2, e2 := estimateDrills[idx].LayerBottomHeight(layer)
			if e1 == nil && e2 == nil {
				hasLayer++
				observeValues = append(observeValues, h1[0])
				estimateValues = append(estimateValues, h2[0])
			}
		}
		if pe, err3 := utils.PercentageError(observeValues, estimateValues); err3 == nil {
			layerPEs[layer] = pe
		} else {
			layerPEs[layer] = -1
		}
	}
	return layerPEs, nil
}
