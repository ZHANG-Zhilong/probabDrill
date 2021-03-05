package utils

import (
	"log"
	"math"
	"probabDrill/internal/entity"
)

func StatProbBlock(drills []entity.Drill, ceil, floor float64) (prob float64) {
	for _, drill := range drills {
		if drill.HasBlock(ceil, floor) {
			prob += 1.0
		}
	}
	prob = prob / float64(len(drills))
	return prob
}
func StatProbBlockWithWeight(drills []entity.Drill, blockCeil, blockFloor float64) (prob float64) {
	if len(drills) < 1 || blockCeil <= blockFloor {
		log.Fatal("error")
		return
	}
	for _, d := range drills {
		if d.HasBlock(blockCeil, blockFloor) {
			prob += d.GetWeight()
		}
	}
	if math.IsNaN(prob) || math.IsInf(prob, 0) {
		log.Fatal("error")
		return
	}
	return prob
}
func StatProbBlockAndLayer(drills []entity.Drill, blockCeil, blockFloor float64, layer int) (prob float64) {
	//p(blockAndLayer)
	log.SetFlags(log.Lshortfile)
	for _, drill := range drills {
		if seq, ok := drill.GetLayerSeq(blockCeil, blockFloor); ok && seq == layer {
			prob += 1.0
		}
	}
	prob = prob / float64(len(drills))
	if math.IsNaN(prob) || math.IsInf(prob, 0) {
		log.Fatal("error")
	}
	return prob
}
func StatProbBlockAndLayerWithWeight(drills []entity.Drill, blockCeil, blockFloor float64, layer int) (prob float64) {
	for _, drill := range drills {
		if seq, ok := drill.GetLayerSeq(blockCeil, blockFloor); ok && seq == layer {
			prob += drill.GetWeight()
		}
	}
	if math.IsNaN(prob) || math.IsInf(prob, 0) {
		log.SetFlags(log.Lshortfile)
		log.Fatal("error")
	}
	return prob
}
func StatProbBlockLayerWithWeight(drills []entity.Drill, blockCeil, blockFloor float64, layer int) (prob float64) {
	//p(block|layer)=p(blockAndLayer)/p(layer)= p(block ∩ layer)/p(layer)  //∩->\cap
	probLayer := StatProbLayerWithWeight(drills, blockCeil, blockFloor, layer)
	probLayerAndBlock := StatProbBlockAndLayerWithWeight(drills, blockCeil, blockFloor, layer)
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
func StatProbBlockLayer(drills []entity.Drill, blockCeil, blockFloor float64, layer int) (prob float64) {
	//p(block|layer) = p(blockAndLayer)/p(layer)
	probLayer := StatProbLayer(drills, blockCeil, blockFloor, layer)
	probBlockAndLayer := StatProbBlockAndLayer(drills, blockCeil, blockFloor, layer)
	if probLayer > 0 {
		prob = probBlockAndLayer / probLayer
	}
	return prob
}
func StatProbLayer(drills []entity.Drill, blockCeil, blockFloor float64, layer int) (prob float64) {
	log.SetFlags(log.Lshortfile)
	if len(drills) < 1 || blockCeil <= blockFloor {
		log.Fatal("error")
		return -1
	}

	//p(layer) = p(block)p(layer|block)+p(blank)p(blank|layer)
	//p(blank)+p(block)=1, p(layer1|block)+p(layer2|block)=1, p(layer1|blank)+p(layer2|blank)=1

	probBlock := StatProbBlock(drills, blockCeil, blockFloor)
	probBlank := 1 - probBlock

	//here is a transformation, that p(block|layer) is general
	//p(layer) means the drill has layer, p(block) means drill has block
	probLayerBlock := StatProbLayerBlock(drills, blockCeil, blockFloor, layer)
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
func StatProbLayerBlockWithWeight(drills []entity.Drill, blockCeil, blockFloor float64, layer int) (prob float64) {
	//p(layer|block) = p(layerAndBlock)/p(block)
	probBlockAndLayerWithWeight := StatProbBlockAndLayerWithWeight(drills, blockCeil, blockFloor, layer)
	probBlocksWithWeight := StatProbBlockWithWeight(drills, blockCeil, blockFloor)
	if probBlocksWithWeight > 0 {
		prob = probBlockAndLayerWithWeight / probBlocksWithWeight
	}
	return prob
}
func StatProbLayerBlock(drills []entity.Drill, blockCeil, blockFloor float64, layer int) (prob float64) {
	//p(layer|block) = p(layerAndBlock)/p(block)
	probBlockAndLayer := StatProbBlockAndLayer(drills, blockCeil, blockFloor, layer)
	probBlocks := StatProbBlock(drills, blockCeil, blockFloor)
	if probBlocks > 0 {
		prob = probBlockAndLayer / probBlocks
	}
	return prob
}
func StatProbLayerWithWeight(drills []entity.Drill, blockCeil, blockFloor float64, layer int) (prob float64) {
	log.SetFlags(log.Lshortfile)
	if len(drills) < 1 || blockCeil <= blockFloor {
		log.Fatal("error")
		return -1
	}

	//p(layer) = p(block)p(layer|block)+p(blank)p(blank|layer)
	//p(blank)+p(block)=1, p(layer1|block)+p(layer2|block)=1, p(layer1|blank)+p(layer2|blank)=1

	probBlockWithWeight := StatProbBlockWithWeight(drills, blockCeil, blockFloor)
	probBlankWithWeight := 1 - probBlockWithWeight

	//here is a transformation, that p(block|layer) is general
	//p(layer) means the drill has layer, p(block) means drill has block

	probLayerBlockWithWeight := StatProbLayerBlockWithWeight(drills, blockCeil, blockFloor, layer)
	probLayerBlankWithWeight := 0.0

	prob = probBlockWithWeight*probLayerBlockWithWeight + probBlankWithWeight*probLayerBlankWithWeight

	if math.IsNaN(prob) || math.IsInf(prob, 0) {
		return -1
		log.Fatal("error")
	}
	return prob
}

