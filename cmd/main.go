<<<<<<< HEAD
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
=======
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
>>>>>>> 34e6069b56f966af97c9e8c24edc3db14aa285d2
