package utils

import (
	"fmt"
	"gonum.org/v1/gonum/mat"
	"log"
	"math"
	probabDrill "probabDrill/conf"
	"probabDrill/internal/entity"
	"runtime/debug"
)

func ProbBlocksWMat(drills []entity.Drill, blocks []float64) (probs *mat.Dense, err error) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	probs = mat.NewDense(len(blocks), 1, nil)
	for bidx := 1; bidx < len(blocks); bidx++ {
		prob := ProbBlockW(drills, blocks[bidx-1], blocks[bidx])
		if math.IsInf(prob, 10) || math.IsNaN(prob) {
			return nil, fmt.Errorf("has invalid prob, %#v", probs)
		}
		probs.Set(bidx, 0, prob)
	}
	return
}
func ProbBlocks(drills []entity.Drill, blocks []float64) (probs *mat.Dense, err error) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	probs = mat.NewDense(len(blocks), 1, nil)
	for bidx := 1; bidx < len(blocks); bidx++ {
		prob := probBlock(drills, blocks[bidx-1], blocks[bidx])
		if math.IsInf(prob, 10) || math.IsNaN(prob) {
			return nil, fmt.Errorf("has invalid prob, %#v", probs)
		}
		probs.Set(bidx, 0, Decimal(prob))
	}
	return
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
func ProbBlockW(drills []entity.Drill, ceil, floor float64) (prob float64) {
	if len(drills) < 1 || ceil <= floor {
		debug.PrintStack()
		log.Println(drills)
		log.Println(len(drills))
		log.Fatal("error")
	}
	for _, d := range drills {
		if d.HasBlock(ceil, floor) {
			prob += d.GetWeight()
		}
	}
	if math.IsNaN(prob) || math.IsInf(prob, 0) {
		log.Fatal("error")
		return
	}
	return prob
}

func ProbLayerWMat(drills []entity.Drill, blocks []float64) (probMat *mat.Dense, err error) {
	probMat = mat.NewDense(len(blocks), probabDrill.StdLen, nil)
	for bidx := 1; bidx < len(blocks); bidx++ {
		for lidx := 1; lidx < probabDrill.StdLen; lidx++ {
			prob := ProbLayerW(drills, blocks[bidx-1], blocks[bidx], lidx)
			probMat.Set(bidx, lidx, prob)
		}
	}
	return probMat, nil
}
func ProbLayerW(drills []entity.Drill, ceil, floor float64, layer int) (prob float64) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if len(drills) < 1 || ceil <= floor {
		log.Fatal("error")
		return -1
	}

	//p(layer) = p(block)p(layer|block)+p(blank)p(blank|layer)
	//p(blank)+p(block)=1, p(layer1|block)+p(layer2|block)=1, p(layer1|blank)+p(layer2|blank)=1

	probBlockWithWeight := ProbBlockW(drills, ceil, floor)
	probBlankWithWeight := 1 - probBlockWithWeight

	//here is a transformation, that p(block|layer) is general
	//p(layer) means the drill has layer, p(block) means drill has block

	probLayerBlockWithWeight := ProbLayerBlockW(drills, ceil, floor, layer)
	probLayerBlankWithWeight := 0.0

	prob = probBlockWithWeight*probLayerBlockWithWeight + probBlankWithWeight*probLayerBlankWithWeight

	if math.IsNaN(prob) || math.IsInf(prob, 0) {
		log.Fatal("error")
	}
	return prob
}
func probLayer(drills []entity.Drill, ceil, floor float64, layer int) (prob float64) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if len(drills) == 0 {
		log.Fatal("input drill is empty.\n")
	} else if ceil <= floor {
		log.Fatalf("input param is invalid ceil<= floor, ceil = %f, floor = %f.", ceil, floor)
	}

	//p(layer) = p(block)p(layer|block)+p(blank)p(blank|layer)
	//p(blank)+p(block)=1, p(layer1|block)+p(layer2|block)=1, p(layer1|blank)+p(layer2|blank)=1

	probBlock := probBlock(drills, ceil, floor)
	probBlank := 1 - probBlock

	//here is a transformation, that p(block|layer) is general
	//p(layer) means the drill has layer, p(block) means drill has block
	probLayerBlock := ProbLayerBlock(drills, ceil, floor, layer)
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

func ProbBlockLayerMatG(drillSet []entity.Drill, blocks []float64) (probsMat *mat.Dense, err error) {
	probsMat = mat.NewDense(len(blocks), probabDrill.StdLen, nil)
	for bidx := 1; bidx < len(blocks); bidx++ {
		for lidx := 1; lidx < probabDrill.StdLen; lidx++ {
			prob, _ := ProbBlockLayerG(drillSet, blocks[bidx-1], blocks[bidx], lidx)
			probsMat.Set(bidx, lidx, Decimal(prob))
		}
	}
	return probsMat, nil
}
func ProbBlockLayerG(drillSet []entity.Drill, ceil, floor float64, layer int) (prob float64, err error) {
	//p(block|layer) = p(blockAndLayer)/p(layer)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	pLayer := probLayer(drillSet, ceil, floor, layer)
	pBlockAndLayer := probBlockAndLayer(drillSet, ceil, floor, layer)
	if pLayer > 0 {
		prob = pBlockAndLayer / pLayer
	}
	return prob, nil
}

func ProbLayerBlockMat(drillSet []entity.Drill, blocks []float64) (probsMat *mat.Dense, err error) {
	probsMat = mat.NewDense(len(blocks), probabDrill.StdLen, nil)
	for bidx := 1; bidx < len(blocks); bidx++ {
		for lidx := 1; lidx < probabDrill.StdLen; lidx++ {
			prob := ProbLayerBlock(drillSet, blocks[bidx-1], blocks[bidx], lidx)
			probsMat.Set(bidx, lidx, Decimal(prob))
		}
	}
	return probsMat, nil
}
func ProbLayerBlock(drills []entity.Drill, ceil, floor float64, layer int) (prob float64) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	//p(layer|block) = p(layerAndBlock)/p(block)
	pBlockAndLayer := probBlockAndLayer(drills, ceil, floor, layer)
	probBlocks := probBlock(drills, ceil, floor)
	if probBlocks > 0 {
		prob = pBlockAndLayer / probBlocks
	}
	return prob
}
func ProbLayerBlockW(drills []entity.Drill, ceil, floor float64, layer int) (prob float64) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	//p(layer|block) = p(layerAndBlock)/p(block)
	probBlockAndLayerWithWeight := probBlockAndLayerW(drills, ceil, floor, layer)
	probBlocksWithWeight := ProbBlockW(drills, ceil, floor)
	if probBlocksWithWeight > 0 {
		prob = probBlockAndLayerWithWeight / probBlocksWithWeight
	}
	return prob
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
func probBlockAndLayerW(drills []entity.Drill, ceil, floor float64, layer int) (prob float64) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	for _, drill := range drills {
		if seq, ok := drill.GetLayerSeq(ceil, floor); ok && seq == layer {
			prob += drill.GetWeight()
		}
	}
	if math.IsNaN(prob) || math.IsInf(prob, 0) {
		log.SetFlags(log.Lshortfile)
		log.Fatal("error")
	}
	return prob
}
