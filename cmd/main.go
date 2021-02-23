package main

import (
	"log"
	"probabDrill/internal/constant"
	"probabDrill/internal/service"
	"probabDrill/internal/utils"
)

func main() {
	log.SetFlags(log.Lshortfile)
	drills := constant.DrillSet()
	virtualDrillsCrossGrid := service.GetGridDrills(drills)
	utils.DisplayDrills(virtualDrillsCrossGrid)
}
