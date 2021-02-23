package main

import (
	"log"
	"probabDrill/internal/constant"
	"probabDrill/internal/utils"
)

func main() {
	log.SetFlags(log.Lshortfile)
	drills := constant.DrillSet()
	virtualDrillsCrossGrid := utils.GetGridDrills(drills)
	utils.DisplayDrills(virtualDrillsCrossGrid)
}
