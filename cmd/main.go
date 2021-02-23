package main

import (
	"probabDrill-main/internal/constant"
	"probabDrill-main/internal/utils"
)
func main() {
	drills := constant.DrillSet()
	virtualDrillsCrossGrid := utils.GetGridDrills(drills)
	utils.DisplayDrills(virtualDrillsCrossGrid)
}
