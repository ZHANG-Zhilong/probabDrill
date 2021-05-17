package stat

import (
	"gonum.org/v1/gonum/mat"
	"log"
	"probabDrill"
	"probabDrill/internal/entity"
)

func ProbLayers(drills []entity.Drill) (probsVec *mat.Dense) {
	probsVec = mat.NewDense(1, int(probabDrill.StdLen), nil)
	for layer := 0; layer < int(probabDrill.StdLen); layer++ {
		prob := probLayer(drills, int(layer))
		probsVec.Set(0, layer, prob)
	}
	return probsVec
}
func ProbBlocks(drills []entity.Drill, blocksHeights *[]float64) (probsVec *mat.Dense) {
	probsVec = mat.NewDense(len(*blocksHeights), 1, nil)
	for idx := 1; idx < len(*blocksHeights); idx++ {
		prob := probBlock(drills, (*blocksHeights)[idx-1], (*blocksHeights)[idx])
		probsVec.Set(idx, 0, prob)
	}
	return
}
func ProbLBs(drills []entity.Drill, blockHeights *[]float64) (probsMatrix *mat.Dense) {
	probsMatrix = mat.NewDense(len(*blockHeights), probabDrill.StdLen, nil)
	for idx := 1; idx < len(*blockHeights); idx++ {
		ceil, floor := (*blockHeights)[idx-1], (*blockHeights)[idx]
		probsVec := probLB(drills, ceil, floor)
		probsMatrix.SetRow(idx, probsVec.RawVector().Data)
	}
	return probsMatrix
}
func ProbBLs(drillsPtr []entity.Drill, blocksPtr []float64) (probMat *mat.Dense) {
	row, col := len(blocksPtr), probabDrill.StdLen
	probMat = mat.NewDense(row, col, nil)
	for layer := 1; layer < probabDrill.StdLen; layer++ {
		prob := probBL(drillsPtr, layer, blocksPtr)
		probMat.SetCol(layer, prob.RawVector().Data)
	}
	return
}
func bayesLBs(drill entity.Drill, nearDrills *[]entity.Drill) () {

}

func ProbLayersW(drills []entity.Drill) (probsVec *mat.Dense) {
	probsVec = mat.NewDense(1, int(probabDrill.StdLen), nil)
	for layer := 0; layer < int(probabDrill.StdLen); layer++ {
		prob := probLayerW(drills, int(layer))
		probsVec.Set(0, layer, prob)
	}
	return probsVec
}
func ProbBlocksW(drills []entity.Drill, blocksHeights []float64) (probsVec *mat.Dense) {
	probsVec = mat.NewDense(len(blocksHeights), 1, nil)
	for idx := 1; idx < len(blocksHeights); idx++ {
		prob := probBlockW(drills, blocksHeights[idx-1], blocksHeights[idx])
		probsVec.Set(idx, 0, prob)
	}
	return
}
func probLayerW(drills []entity.Drill, layer int) (prob float64) {
	log.SetFlags(log.Lshortfile)
	if len(drills) == 0 {
		log.Fatal("error")
	}
	var has float64
	for _, d := range drills {
		if d.HasLayer(layer) > 0 {
			has += 1.0 * d.GetWeight()
		}
	}
	return has
}
func probBlockW(drills []entity.Drill, blockCeil, blockFloor float64) (prob float64) {
	if len(drills) < 1 || blockCeil <= blockFloor {
		log.Fatal("error")
		return
	}
	for _, d := range drills {
		if d.HasBlock(blockCeil, blockFloor) {
			prob += d.GetWeight()
		}
	}
	return prob
}

func probLayer(drills []entity.Drill, layer int) (prob float64) {
	log.SetFlags(log.Lshortfile)
	if len(drills) == 0 {
		log.Fatal("error")
	}
	var total, has float64
	total = float64(len(drills))
	for _, d := range drills {
		if d.HasLayer(layer) > 0 {
			has += 1.0
		}
	}
	return has / total
}
func probBlock(drills []entity.Drill, ceil, floor float64) (prob float64) {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	if len(drills) == 0 {
		log.Fatal("error, stat/stat")
	}
	var has, total float64
	total = float64(len(drills))
	for _, d := range drills {
		if d.HasBlock(ceil, floor) {
			has += 1.0
		}
	}
	return has / total
}
func probLB(drills []entity.Drill, ceil, floor float64) (probsVec *mat.VecDense) {
	probsVec = mat.NewVecDense(probabDrill.StdLen, nil)
	pblock := probBlock(drills, ceil, floor)
	for layer := 1; layer < probabDrill.StdLen; layer++ {
		prob := probLAndB(drills, layer, ceil, floor)
		prob = prob / pblock
		probsVec.SetVec(layer, prob)
	}
	return probsVec
}
func probLAndB(drills []entity.Drill, layer int, ceil, floor float64) (prob float64) {
	var total, has float64
	total = float64(len(drills))
	for _, d := range drills {
		if d.HasBlock(ceil, floor) && d.HasLayer(layer) > 0 {
			has += 1.0
		}
	}
	return has / total
}
func probBL(drillsPtr []entity.Drill, layer int, blocks []float64) (probVec *mat.VecDense) {
	probVec = mat.NewVecDense(len(blocks), nil)
	pl := probLayer(drillsPtr, layer)
	if pl == 0 {
		return probVec
	}
	for idx := 1; idx < len(blocks); idx++ {
		plAndb := probLAndB(drillsPtr, layer, (blocks)[idx-1], (blocks)[idx])
		probVec.SetVec(idx, plAndb/pl)
	}
	return
}
