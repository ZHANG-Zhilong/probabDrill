package main

import (
	"awesome/internal/constant"
	"awesome/internal/utils"
)

func main() {
	drills := constant.DrillSet()
	virtualDrillsCrossGrid := utils.GetGridDrills(drills)
	utils.DisplayDrills(virtualDrillsCrossGrid)
}
