package utils

import (
	"log"
	"math"
	"probabDrill"
	"probabDrill/internal/entity"
)

func SetClassicalIdwWeights(center entity.Drill, nearDrills []entity.Drill) (weights []float64) {
	log.SetFlags(log.Lshortfile)
	var (
		weightSum       float64
		hasZeroDistance bool
		zeroIdx         int
	)

	//get distance
	for idx, aroundDrill := range nearDrills {
		dist := center.Distance(aroundDrill)
		weights = append(weights, dist) //as distance
		if dist < 1e-1 {
			hasZeroDistance = true
			zeroIdx = idx
		}
	}
	if hasZeroDistance {
		for idx, _ := range weights {
			weights[idx] = 0
		}
		if weights != nil && zeroIdx >= 0 && zeroIdx < len(weights) {
			weights[zeroIdx] = 1
		}
	} else {
		for idx, _ := range weights { //cal weight
			weights[idx] = 1e7 / math.Pow(weights[idx], probabDrill.IdwPow)
			weightSum += weights[idx]
		}
		for idx, _ := range weights { //归一化, and set int the drill.
			weights[idx] = Decimal(weights[idx] / weightSum)
			nearDrills[idx].SetWeight(weights[idx])
		}
		weightSum = 0
		for _, w := range weights {
			weightSum += w
		}
		weightSum = Decimal(weightSum)
	}

	if math.Abs(weightSum-1) > 1e-1 {
		log.Println(weights)
		log.Fatalf("error, total weight:%f\n", weightSum)
	}
	return weights
}
