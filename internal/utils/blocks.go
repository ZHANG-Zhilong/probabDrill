package utils

import (
	"log"
	"math"
	"probabDrill/internal/entity"
)

func MakeBlocks(drills []entity.Drill, res float64) (blocks []float64) {
	drillsCeil, drillsFloor := -math.MaxFloat64, math.MaxFloat64
	for _, d := range drills {
		drillsCeil = math.Max(d.Z, drillsCeil)
		drillsFloor = math.Min(drillsFloor, d.LayerHeights[len(d.LayerHeights)-1])
	}
	blocks = append(blocks, Decimal(drillsCeil))

	for drillsCeil-res > drillsFloor {
		blocks = append(blocks, Decimal(drillsCeil-res))
		drillsCeil = drillsCeil - res
	}

	//the last block may be un-standard block length, whose length may less than res
	blocks = append(blocks, Decimal(drillsFloor))
	return
}
func ExplodedHeights(blocks []float64, ceil, floor float64) (heights []float64) {
	idxa := int(0)
	for idx, h := range blocks {
		if h < ceil {
			idxa = idx
			break
		}
	}
	heights = append(heights, ceil)

	for idx := idxa; idx < len(blocks); idx++ {
		if blocks[idx] <= ceil && blocks[idx] >= floor {
			heights = append(heights, blocks[idx])
		}
	}
	if heights[len(heights)-1] > floor {
		heights = append(heights, floor)
	}
	return
}
func BlocksIndex(blocks []float64, ceil, floor float64) (index int) {
	log.SetFlags(log.Lshortfile)
	if ceil > blocks[0] || floor < blocks[len(blocks)-1] {
		return 0
	}
	for idx := 1; idx < len(blocks); idx++ {
		if ceil <= blocks[idx-1] && floor >= blocks[idx] {
			return idx
		}
	}
	return 0
}
