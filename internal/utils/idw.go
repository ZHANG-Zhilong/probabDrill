package utils

import (
	"log"
	"math"
	probabDrill "probabDrill/conf"
	"probabDrill/internal/entity"
)

func SetClassicalIdwWeights(cdrill entity.Drill, nearDrills []entity.Drill) (weights []float64) {
	log.SetFlags(log.Lshortfile)
	var (
		weightSum       float64
		hasZeroDistance bool
		zeroIdx         int
	)
	if len(nearDrills) == 0 {
		log.Fatal("invalid input, that near drills is empty.\n")
	}
	//过滤掉过近的钻孔数据
	//var nearDrills2 []entity.Drill
	//for _, d := range nearDrills {
	//	if cdrill.Distance(d) > probabDrill.MinDrillDist {
	//		nearDrills2 = append(nearDrills2, d)
	//	}
	//}
	//nearDrills = nearDrills2

	//get distance
	for idx, aroundDrill := range nearDrills {
		dist := cdrill.Distance(aroundDrill)
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
		log.Println("weights", weights)
		log.Printf("center drill, %#v\n", cdrill)
		for _, d := range nearDrills {
			log.Printf("dist: %f, neardrill:%#v\n", cdrill.Distance(d), d)
		}
		log.Fatalf("error, total weight:%f\n", weightSum)
	}
	return weights
}
