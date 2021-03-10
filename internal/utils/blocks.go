package utils

import (
	"log"
	"math"
	"probabDrill/internal/entity"
	"runtime/debug"
)

func MakeBlocks(drillSet []entity.Drill, res float64) (blocksHeight []float64) {
	drillsCeil, drillsFloor := -math.MaxFloat64, math.MaxFloat64
	for _, d := range drillSet {
		if d.Z > drillsCeil {
			drillsCeil = d.Z
		}
		if d.LayerHeights[len(d.LayerHeights)-1] < drillsFloor {
			drillsFloor = d.LayerHeights[len(d.LayerHeights)-1]
		}
	}

	blocksHeight = append(blocksHeight, drillsCeil)

	for drillsCeil-res > drillsFloor {
		blocksHeight = append(blocksHeight, drillsCeil-res)
		drillsCeil = drillsCeil - res
	}

	//the last block may be un-standard block length, whose length may less than res
	blocksHeight = append(blocksHeight, drillsFloor)

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
	for idx := 1; idx < len(blocks); idx++ {
		if ceil <= blocks[idx-1] && floor >= blocks[idx] {
			return idx
		}
	}
	debug.PrintStack()
	log.Fatal("error")
	return -1
}

