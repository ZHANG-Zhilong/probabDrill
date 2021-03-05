package stat

import (
	"gonum.org/v1/gonum/mat"
	"log"
	"probabDrill/internal/constant"
	"probabDrill/internal/entity"
)

func ProbLayers(drills *[]entity.Drill) (probsVec *mat.VecDense) {
	probsVec = mat.NewVecDense(int(constant.StdLen), nil)
	for layer := 0; layer < int(constant.StdLen); layer++ {
		prob := probLayer(drills, int(layer))
		probsVec.SetVec(layer, prob)
	}
	return probsVec
}
func ProbBlocks(drills *[]entity.Drill, blocksHeights *[]float64) (probsVec *mat.VecDense) {
	probsVec = mat.NewVecDense(len(*blocksHeights), nil)
	for idx := 1; idx < len(*blocksHeights); idx++ {
		prob := probBlock(drills, (*blocksHeights)[idx-1], (*blocksHeights)[idx])
		probsVec.SetVec(idx, prob)
	}
	return
}
func ProbLBs(drills *[]entity.Drill, blockHeights *[]float64) (probsMatrix *mat.Dense) {
	probsMatrix = mat.NewDense(len(*blockHeights), constant.StdLen, nil)
	for idx := 1; idx < len(*blockHeights); idx++ {
		ceil, floor := (*blockHeights)[idx-1], (*blockHeights)[idx]
		probsVec := probLB(drills, ceil, floor)
		probsMatrix.SetRow(idx, probsVec.RawVector().Data)
	}
	return probsMatrix
}
func ProbBLs(drillsPtr *[]entity.Drill, blocksPtr *[]float64) (probMat *mat.Dense) {
	row, col := len(*blocksPtr), constant.StdLen
	probMat = mat.NewDense(row, col, nil)
	for layer := 1; layer < constant.StdLen; layer++ {
		prob := probBL(drillsPtr, layer, blocksPtr)
		probMat.SetCol(layer, prob.RawVector().Data)
	}
	return
}

func probLayer(drills *[]entity.Drill, layer int) (prob float64) {
	log.SetFlags(log.Lshortfile)
	if len(*drills) == 0 {
		log.Fatal("error")
	}
	var total, has float64
	total = float64(len(*drills))
	for _, d := range *drills {
		if d.HasLayer(layer) > 0 {
			has++
		}
	}
	return has / total
}
func probBlock(drills *[]entity.Drill, ceil, floor float64) (prob float64) {
	log.SetFlags(log.Lshortfile)
	if len(*drills) == 0 {
		log.Fatal("error")
	}
	var has, total float64
	total = float64(len(*drills))
	for _, d := range *drills {
		if d.HasBlock(ceil, floor) {
			has += 1.0
		}
	}
	return has / total
}
func probLB(drills *[]entity.Drill, ceil, floor float64) (probsVec *mat.VecDense) {
	probsVec = mat.NewVecDense(constant.StdLen, nil)
	pblock := probBlock(drills, ceil, floor)
	for layer := 1; layer < constant.StdLen; layer++ {
		prob := probLAndB(drills, layer, ceil, floor)
		prob = prob / pblock
		probsVec.SetVec(layer, prob)
	}
	return probsVec
}
func probLAndB(drills *[]entity.Drill, layer int, ceil, floor float64) (prob float64) {
	var total, has float64
	total = float64(len(*drills))
	for _, d := range *drills {
		if d.HasBlock(ceil, floor) && d.HasLayer(layer) > 0 {
			has += 1.0
		}
	}
	return has / total
}
func probBL(drillsPtr *[]entity.Drill, layer int, blocksPtr *[]float64) (probVec *mat.VecDense) {
	probVec = mat.NewVecDense(len(*blocksPtr), nil)
	pl := probLayer(drillsPtr, layer)
	if pl == 0 {
		return probVec
	}
	for idx := 1; idx < len(*blocksPtr); idx++ {
		plAndb := probLAndB(drillsPtr, layer, (*blocksPtr)[idx-1], (*blocksPtr)[idx])
		probVec.SetVec(idx, plAndb/pl)
	}
	return
}
