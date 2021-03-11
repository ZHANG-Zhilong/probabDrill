package constant

import (
	"gonum.org/v1/gonum/mat"
	"log"
	"math"
	"probabDrill"
	"probabDrill/internal/entity"
	"sync"
)

var (
	blockLayerRDrillMatOnce sync.Once
	blockLayerRDrillMat     *mat.Dense
)

func GetBlockLayerRDrillMat() (probs *mat.Dense) {
	blockLayerRDrillMatOnce.Do(initBlockLayerRDrillMat)
	return blockLayerRDrillMat
}
func initBlockLayerRDrillMat() {
	drills := GetDrillSet()
	blocks := GetBlocksR()
	blockLayerRDrillMat, _ = probBlockLayerMatG(drills, blocks)
}

var (
	blockLayerHDrillMatOnce sync.Once
	blockLayerHDrillMat     *mat.Dense
)

func GetBlockLayerHDrillMat() (probs *mat.Dense) {
	blockLayerHDrillMatOnce.Do(initBlockLayerHDrillMat)
	return blockLayerHDrillMat
}
func initBlockLayerHDrillMat() {
	drills := GetHelpDrillSet()
	blocks := GetBlocksH()
	blockLayerHDrillMat, _ = probBlockLayerMatG(drills, blocks)
}

func probBlockLayerMatG(drillSet []entity.Drill, blocks []float64) (probsMat *mat.Dense, err error) {
	probsMat = mat.NewDense(len(blocks), probabDrill.StdLen, nil)
	for bidx := 1; bidx < len(blocks); bidx++ {
		for lidx := 1; lidx < probabDrill.StdLen; lidx++ {
			prob, _ := probBlockLayerG(drillSet, blocks[bidx-1], blocks[bidx], lidx)
			probsMat.Set(bidx, lidx, prob)
		}
	}
	return probsMat, nil
}
func probBlockLayerG(drillSet []entity.Drill, ceil, floor float64, layer int) (prob float64, err error) {
	//p(block|layer) = p(blockAndLayer)/p(layer)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	pLayer := probLayer(drillSet, ceil, floor, layer)
	pBlockAndLayer := probBlockAndLayer(drillSet, ceil, floor, layer)
	if pLayer > 0 {
		prob = pBlockAndLayer / pLayer
	}
	return prob, nil
}

func probBlockAndLayer(drills []entity.Drill, ceil, floor float64, layer int) (prob float64) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	//p(blockAndLayer)
	for _, drill := range drills {
		if seq, ok := drill.GetLayerSeq(ceil, floor); ok && seq == layer {
			prob += 1.0
		}
	}
	prob = prob / float64(len(drills))
	if math.IsNaN(prob) || math.IsInf(prob, 0) {
		log.Fatal("error")
	}
	return prob
}
func probLayer(drills []entity.Drill, ceil, floor float64, layer int) (prob float64) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if len(drills) < 1 || ceil <= floor {
		log.Fatal("error")
		return -1
	}

	//p(layer) = p(block)p(layer|block)+p(blank)p(blank|layer)
	//p(blank)+p(block)=1, p(layer1|block)+p(layer2|block)=1, p(layer1|blank)+p(layer2|blank)=1

	probBlock := probBlock(drills, ceil, floor)
	probBlank := 1 - probBlock

	//here is a transformation, that p(block|layer) is general
	//p(layer) means the drill has layer, p(block) means drill has block
	probLayerBlock := probLayerBlock(drills, ceil, floor, layer)
	probLayerBlank := 0.0

	prob = probBlock*probLayerBlock + probBlank*probLayerBlank

	if math.IsNaN(prob) || math.IsInf(prob, 0) {
		log.Fatal("error")
	}
	if prob > 1 {
		prob = 1
	}
	return prob
}
func probBlock(drills []entity.Drill, ceil, floor float64) (prob float64) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	for _, drill := range drills {
		if drill.HasBlock(ceil, floor) {
			prob += 1.0
		}
	}
	prob = prob / float64(len(drills))
	return prob
}
func probLayerBlock(drills []entity.Drill, ceil, floor float64, layer int) (prob float64) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	//p(layer|block) = p(layerAndBlock)/p(block)
	pBlockAndLayer := probBlockAndLayer(drills, ceil, floor, layer)
	probBlocks := probBlock(drills, ceil, floor)
	if probBlocks > 0 {
		prob = pBlockAndLayer / probBlocks
	}
	return prob
}
