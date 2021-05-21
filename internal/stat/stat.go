package stat

import (
	"github.com/spf13/viper"
	"gonum.org/v1/gonum/mat"
	"log"
	"probabDrill/apps/probDrill/model"
)

func ProbLayers(drills []model.Drill) (probsVec *mat.Dense) {
	probsVec = mat.NewDense(1, viper.GetInt("StdLen"), nil)
	for layer := 0; layer < viper.GetInt("StdLen"); layer++ {
		prob := probLayer(drills, layer)
		probsVec.Set(0, layer, prob)
	}
	return probsVec
}
func ProbBlocks(drills []model.Drill, blocksHeights *[]float64) (probsVec *mat.Dense) {
	probsVec = mat.NewDense(len(*blocksHeights), 1, nil)
	for idx := 1; idx < len(*blocksHeights); idx++ {
		prob := probBlock(drills, (*blocksHeights)[idx-1], (*blocksHeights)[idx])
		probsVec.Set(idx, 0, prob)
	}
	return
}
func ProbLBs(drills []model.Drill, blockHeights *[]float64) (probsMatrix *mat.Dense) {
	probsMatrix = mat.NewDense(len(*blockHeights), viper.GetInt("StdLen"), nil)
	for idx := 1; idx < len(*blockHeights); idx++ {
		ceil, floor := (*blockHeights)[idx-1], (*blockHeights)[idx]
		probsVec := probLB(drills, ceil, floor)
		probsMatrix.SetRow(idx, probsVec.RawVector().Data)
	}
	return probsMatrix
}
func ProbBLs(drillsPtr []model.Drill, blocksPtr []float64) (probMat *mat.Dense) {
	row, col := len(blocksPtr), viper.GetInt("StdLen")
	probMat = mat.NewDense(row, col, nil)
	for layer := 1; layer < viper.GetInt("StdLen"); layer++ {
		prob := probBL(drillsPtr, layer, blocksPtr)
		probMat.SetCol(layer, prob.RawVector().Data)
	}
	return
}
func bayesLBs(drill model.Drill, nearDrills *[]model.Drill) () {

}

func ProbLayersW(drills []model.Drill) (probsVec *mat.Dense) {
	probsVec = mat.NewDense(1, int(viper.GetInt("StdLen")), nil)
	for layer := 0; layer < int(viper.GetInt("StdLen")); layer++ {
		prob := probLayerW(drills, int(layer))
		probsVec.Set(0, layer, prob)
	}
	return probsVec
}
func ProbBlocksW(drills []model.Drill, blocksHeights []float64) (probsVec *mat.Dense) {
	probsVec = mat.NewDense(len(blocksHeights), 1, nil)
	for idx := 1; idx < len(blocksHeights); idx++ {
		prob := probBlockW(drills, blocksHeights[idx-1], blocksHeights[idx])
		probsVec.Set(idx, 0, prob)
	}
	return
}
func probLayerW(drills []model.Drill, layer int) (prob float64) {
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
func probBlockW(drills []model.Drill, blockCeil, blockFloor float64) (prob float64) {
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

func probLayer(drills []model.Drill, layer int) (prob float64) {
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
func probBlock(drills []model.Drill, ceil, floor float64) (prob float64) {
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
func probLB(drills []model.Drill, ceil, floor float64) (probsVec *mat.VecDense) {
	probsVec = mat.NewVecDense(viper.GetInt("StdLen"), nil)
	pblock := probBlock(drills, ceil, floor)
	for layer := 1; layer < viper.GetInt("StdLen"); layer++ {
		prob := probLAndB(drills, layer, ceil, floor)
		prob = prob / pblock
		probsVec.SetVec(layer, prob)
	}
	return probsVec
}
func probLAndB(drills []model.Drill, layer int, ceil, floor float64) (prob float64) {
	var total, has float64
	total = float64(len(drills))
	for _, d := range drills {
		if d.HasBlock(ceil, floor) && d.HasLayer(layer) > 0 {
			has += 1.0
		}
	}
	return has / total
}
func probBL(drillsPtr []model.Drill, layer int, blocks []float64) (probVec *mat.VecDense) {
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
