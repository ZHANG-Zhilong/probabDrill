package constant

import (
	"github.com/spf13/viper"
	"log"
	"sync"
)

var (
	blocksOnceR sync.Once
	blocksR     []float64
)

func BlockIndexR(ceil, floor float64) int {
	blocksOnceR.Do(initBlocksR)
	log.SetFlags(log.Lshortfile)
	if ceil > blocksR[0] || floor < blocksR[len(blocksR)-1] {
		return 0
	}
	for idx := 1; idx < len(blocksR); idx++ {
		if ceil <= blocksR[idx-1] && floor >= blocksR[idx] {
			return idx
		}
	}
	return 0
}
func GetBlocksR() []float64 {
	blocksOnceR.Do(initBlocksR)
	if blocksR != nil {
		return blocksR
	}
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.Fatal("error")
	return nil
}
func initBlocksR() () {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	drillsCeil, drillsFloor = GetRealDrillCF()
	blocksR = append(blocksR, decimal(drillsCeil))
	for drillsCeil-viper.GetFloat64("BlockResZ") > drillsFloor {
		blocksR = append(blocksR, decimal(drillsCeil-viper.GetFloat64("BlockResZ")))
		drillsCeil = drillsCeil - viper.GetFloat64("BlockResZ")
	}
	//the last block may be un-standard block length, whose length may less than res
	blocksR = append(blocksR, decimal(drillsFloor))
}

var (
	blocksOnceH sync.Once
	blocksH     []float64
)

func BlockIndexH(ceil, floor float64) int {
	blocksOnceH.Do(initBlocksH)
	log.SetFlags(log.Lshortfile)
	if ceil > blocksH[0] || floor < blocksH[len(blocksH)-1] {
		return 0
	}
	for idx := 1; idx < len(blocksH); idx++ {
		if ceil <= blocksH[idx-1] && floor >= blocksH[idx] {
			return idx
		}
	}
	return 0
}
func GetBlocksH() []float64 {
	blocksOnceH.Do(initBlocksH)
	if blocksH != nil {
		return blocksH
	}
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.Fatal("error")
	return nil
}
func initBlocksH() () {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	drillsCeil, drillsFloor = GetHelpDrillsCF()
	blocksH = append(blocksH, decimal(drillsCeil))
	for drillsCeil-viper.GetFloat64("BlockResZ") > drillsFloor {
		blocksH = append(blocksH, decimal(drillsCeil-viper.GetFloat64("BlockResZ")))
		drillsCeil = drillsCeil - viper.GetFloat64("BlockResZ")
	}
	//the last block may be un-standard block length, whose length may less than res
	blocksH = append(blocksH, decimal(drillsFloor))
}
