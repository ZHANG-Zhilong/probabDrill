package utils

import (
	"log"
	"math"
	"probabDrill/internal/entity"
)

func statProbBlock(drills []entity.Drill, ceil, floor float64) (prob float64) {
	for _, drill := range drills {
		if hasBlock(drill, ceil, floor) {
			prob += 1.0
		}
	}
	prob = prob / float64(len(drills))
	return prob
}
func statProbBlockWithWeight(drills []entity.Drill, blockCeil, blockFloor float64) (prob float64) {
	if len(drills) < 1 || blockCeil <= blockFloor {
		log.Fatal("error")
		return
	}
	for _, d := range drills {
		if hasBlock(d, blockCeil, blockFloor) {
			prob += d.GetWeight()
		}
	}
	if math.IsNaN(prob) || math.IsInf(prob, 0) {
		log.Fatal("error")
		return
	}
	return prob
}
func statProbBlockAndLayer(drills []entity.Drill, blockCeil, blockFloor float64, layer int64) (prob float64) {
	//p(blockAndLayer)
	log.SetFlags(log.Lshortfile)
	for _, drill := range drills {
		if seq, ok := getLayerSeq(drill, blockCeil, blockFloor); ok && seq == layer {
			prob += 1.0
		}
	}
	prob = prob / float64(len(drills))
	if math.IsNaN(prob) || math.IsInf(prob, 0) {
		log.Fatal("error")
	}
	return prob
}
func statProbBlockAndLayerWithWeight(drills []entity.Drill, blockCeil, blockFloor float64, layer int64) (prob float64) {
	for _, drill := range drills {
		if seq, ok := getLayerSeq(drill, blockCeil, blockFloor); ok && seq == layer {
			prob += drill.GetWeight()
		}
	}
	if math.IsNaN(prob) || math.IsInf(prob, 0) {
		log.SetFlags(log.Lshortfile)
		log.Fatal("error")
	}
	return prob
}
func statProbBlockLayerWithWeight(drills []entity.Drill, blockCeil, blockFloor float64, layer int64) (prob float64) {
	//p(block|layer)=p(blockAndLayer)/p(layer)= p(block ∩ layer)/p(layer)  //∩->\cap
	probLayer := statProbLayerWithWeight(drills, blockCeil, blockFloor, layer)
	probLayerAndBlock := statProbBlockAndLayerWithWeight(drills, blockCeil, blockFloor, layer)
	prob = probLayerAndBlock / probLayer
	if math.IsNaN(prob) || math.IsInf(prob, 0) {
		log.Fatal("error")
		return
	}
	if prob == 0 {
		log.SetFlags(log.LstdFlags | log.LstdFlags)
		log.Printf("p(layer) %f, p(layerAndblock) %f\n", probLayer, probLayerAndBlock)
	}
	return prob
}
func statProbBlockLayer(drills []entity.Drill, blockCeil, blockFloor float64, layer int64) (prob float64) {
	//p(block|layer) = p(blockAndLayer)/p(layer)
	probLayer := statProbLayer(drills, blockCeil, blockFloor, layer)
	probBlockAndLayer := statProbBlockAndLayer(drills, blockCeil, blockFloor, layer)
	if probLayer > 0 {
		prob = probBlockAndLayer / probLayer
	}
	return prob
}
func statProbLayer(drills []entity.Drill, blockCeil, blockFloor float64, layer int64) (prob float64) {
	log.SetFlags(log.Lshortfile)
	if len(drills) < 1 || blockCeil <= blockFloor {
		log.Fatal("error")
		return -1
	}

	//p(layer) = p(block)p(layer|block)+p(blank)p(blank|layer)
	//p(blank)+p(block)=1, p(layer1|block)+p(layer2|block)=1, p(layer1|blank)+p(layer2|blank)=1

	probBlock := statProbBlock(drills, blockCeil, blockFloor)
	probBlank := 1 - probBlock

	//here is a transformation, that p(block|layer) is general
	//p(layer) means the drill has layer, p(block) means drill has block
	probLayerBlock := statProbLayerBlock(drills, blockCeil, blockFloor, layer)
	probLayerBlank := 0.0

	prob = probBlock*probLayerBlock + probBlank*probLayerBlank

	if math.IsNaN(prob) || math.IsInf(prob, 0) {
		return
		log.Fatal("error")
	}
	if prob > 1 {
		prob = 1
	}
	return prob
}
func statProbLayerBlockWithWeight(drills []entity.Drill, blockCeil, blockFloor float64, layer int64) (prob float64) {
	//p(layer|block) = p(layerAndBlock)/p(block)
	probBlockAndLayerWithWeight := statProbBlockAndLayerWithWeight(drills, blockCeil, blockFloor, layer)
	probBlocksWithWeight := statProbBlockWithWeight(drills, blockCeil, blockFloor)
	if probBlocksWithWeight > 0 {
		prob = probBlockAndLayerWithWeight / probBlocksWithWeight
	}
	return prob
}
func statProbLayerBlock(drills []entity.Drill, blockCeil, blockFloor float64, layer int64) (prob float64) {
	//p(layer|block) = p(layerAndBlock)/p(block)
	probBlockAndLayer := statProbBlockAndLayer(drills, blockCeil, blockFloor, layer)
	probBlocks := statProbBlock(drills, blockCeil, blockFloor)
	if probBlocks > 0 {
		prob = probBlockAndLayer / probBlocks
	}
	return prob
}
func statProbLayerWithWeight(drills []entity.Drill, blockCeil, blockFloor float64, layer int64) (prob float64) {
	log.SetFlags(log.Lshortfile)
	if len(drills) < 1 || blockCeil <= blockFloor {
		log.Fatal("error")
		return -1
	}

	//p(layer) = p(block)p(layer|block)+p(blank)p(blank|layer)
	//p(blank)+p(block)=1, p(layer1|block)+p(layer2|block)=1, p(layer1|blank)+p(layer2|blank)=1

	probBlockWithWeight := statProbBlockWithWeight(drills, blockCeil, blockFloor)
	probBlankWithWeight := 1 - probBlockWithWeight

	//here is a transformation, that p(block|layer) is general
	//p(layer) means the drill has layer, p(block) means drill has block

	probLayerBlockWithWeight := statProbLayerBlockWithWeight(drills, blockCeil, blockFloor, layer)
	probLayerBlankWithWeight := 0.0

	prob = probBlockWithWeight*probLayerBlockWithWeight + probBlankWithWeight*probLayerBlankWithWeight

	if math.IsNaN(prob) || math.IsInf(prob, 0) {
		return -1
		log.Fatal("error")
	}
	return prob
}
