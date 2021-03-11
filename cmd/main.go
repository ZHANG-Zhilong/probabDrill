package main

import (
	"log"
	"probabDrill/internal/constant"
	"probabDrill/internal/entity"
	"probabDrill/internal/service"
)

func main() {
	log.SetFlags(log.Lshortfile)
	drillSet := constant.GetDrillSet()
	virtualDrillsCrossGrid := service.GetGridDrills(drillSet)
	entity.DisplayDrills(virtualDrillsCrossGrid)

}

