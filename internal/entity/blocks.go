package entity

import (
	"log"
	"probabDrill"
	"sync"
)

var blocksOnce sync.Once
var blocks []float64

func BlockIndex(ceil, floor float64) (idx int) {
	blocksOnce.Do(initBlocks)
	log.SetFlags(log.Lshortfile)
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

func GetBlocks() []float64 {
	once.Do(initBlocks)
	return blocks
}
func initBlocks() () {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	drillsCeil, drillsFloor = GetDrillsCF()
	blocks = append(blocks, decimal(drillsCeil))
	for drillsCeil-probabDrill.BlockResZ > drillsFloor {
		blocks = append(blocks, decimal(drillsCeil-probabDrill.BlockResZ))
		drillsCeil = drillsCeil - probabDrill.BlockResZ
	}
	//the last block may be un-standard block length, whose length may less than res
	blocks = append(blocks, decimal(drillsFloor))
}
