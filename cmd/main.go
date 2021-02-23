package main

import (
	"log"
	"probabDrill-main/internal/constant"
	"probabDrill-main/internal/utils"
)

func main() {
	log.SetFlags(log.Lshortfile)
	drills := constant.DrillSet()
	virtualDrillsCrossGrid := utils.GetGridDrills(drills)
	utils.DisplayDrills(virtualDrillsCrossGrid)
	//for _, d := range drills {
	//	log.Println(d.X, "\t", d.Y)
	//}
}
