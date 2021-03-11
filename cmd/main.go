package main

import (
	"log"
	"probabDrill/internal/entity"
	"probabDrill/internal/service"
)

func main() {
	log.SetFlags(log.Lshortfile)
	drillSet := entity.DrillSet()
	virtualDrillsCrossGrid := service.GetGridDrills(drillSet)
	entity.DisplayDrills(virtualDrillsCrossGrid)

}

