package utils

import (
	"fmt"
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
func InterceptBlocks(blocks []float64, top, bottom float64) (heights []float64) {
	//find the biggest element in blocks, which is smaller than top
	//find the smallest element in blocks, which is smaller than bottom
	idxa, idxb := 0, 0
	for idx, h := range blocks {
		if h < top && idxa == 0 {
			idxa = idx
		}
		if h < bottom && idxb == 0 {
			idxb = idx
			break
		}
	}
	heights = append(heights, Decimal(top))
	heights = append(heights, blocks[idxa:idxb]...)
	heights = append(heights, Decimal(bottom))

	return
}
func BlocksIndex(blocks []float64, ceil, floor float64) (index int, err error) {
	log.SetFlags(log.Lshortfile)
	if ceil > blocks[0] || floor < blocks[len(blocks)-1] {
		return 0, nil
	}
	for idx := 1; idx < len(blocks); idx++ {
		if ceil <= blocks[idx-1] && floor >= blocks[idx] {
			return idx, nil
		}
	}
	return -1, fmt.Errorf(":invalid blocks index found")
}
