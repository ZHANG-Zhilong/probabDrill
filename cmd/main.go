package main

import (
	"log"
	"probabDrill/internal/constant"
	"probabDrill/internal/service"
	"probabDrill/internal/utils"
)

func main() {
	log.SetFlags(log.Lshortfile)
	drillSet := constant.DrillSet()
	virtualDrillsCrossGrid := service.GetGridDrills(drillSet)
	utils.DisplayDrills(virtualDrillsCrossGrid)
}
