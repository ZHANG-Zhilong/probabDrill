package service

import (
	"fmt"
	"log"
	"probabDrill/internal/constant"
	"probabDrill/internal/entity"
	"probabDrill/internal/utils"
	"testing"
)

func TestGenerateVirtualDrill(t *testing.T) {

	log.SetFlags(log.Lshortfile)
	drill := constant.DrillSet()[1]

	log.Println("real drill")
	drill.Print()

	drillSet := constant.DrillSet()
	virtualDrill := GenVDrillM1(drillSet, drill.X+1, drill.Y+1)

	log.Println("virtual drill")
	virtualDrill.Print()
}
func TestGenerateVirtualDrill2(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	drill1 := entity.Drill{
		Name: "1", X: 0, Y: 0, Z: 0,
		Layers:       []int{0, 1, 2, 3},
		LayerHeights: []float64{0, -1, -2, -3},
	}
	drill2 := entity.Drill{
		Name: "2", X: 1, Y: 0, Z: 0,
		Layers:       []int{0, 1, 2, 2},
		LayerHeights: []float64{0, -1, -2, -3},
	}
	drill3 := entity.Drill{
		Name: "3", X: 0, Y: 1, Z: 0,
		Layers:       []int{0, 1, 2, 3},
		LayerHeights: []float64{0, -1.2, -2.3, -3},
	}
	drill4 := entity.Drill{
		Name: "4", X: 1, Y: 1, Z: 0,
		Layers:       []int{0, 1, 2, 3},
		LayerHeights: []float64{0, -1.5, -2.3, -3},
	}
	drillSet := []entity.Drill{drill1, drill2, drill3, drill4}
	blocks := utils.MakeBlocks(drillSet, 0.02)
	fmt.Println(blocks)
	var virtualDrills []entity.Drill
	for x := 0.0; x < 1; x += 0.1 {
		virtualDrills = append(virtualDrills, GenVDrillM1(drillSet, x, 0.5))
	}
	for _, v := range virtualDrills {
		v.Print()
	}
	utils.DrawDrills(virtualDrills, "m1.svg")
}
func TestGenVDrillM1(t *testing.T) {
	drillNames := []string{"TZZK92", "TZJT31", "TZZK40", "TZJT28", "TZZK69", "TZZK70", "TZZK72"}
	var drills []entity.Drill
	for _, name := range drillNames {
		if drill, ok := constant.GetDrillByName(name); ok {
			drills = append(drills, drill)
		}
	}
	var virtualDrills []entity.Drill
	for idx := 1; idx < len(drills); idx++ {
		vdrills := GenVDrillsBetween(drills[idx-1], drills[idx], 3, GenVDrillM1)
		virtualDrills = append(virtualDrills, vdrills...)
	}
	utils.DrawDrills(virtualDrills, "./m1+idw.svg")
}
func TestGenVDrillIDW(t *testing.T) {
	drillNames := []string{"TZZK92", "TZJT31", "TZZK40", "TZJT28", "TZZK69", "TZZK70", "TZZK72"}
	var drills []entity.Drill
	for _, name := range drillNames {
		if drill, ok := constant.GetDrillByName(name); ok {
			drills = append(drills, drill)
		}
	}
	var vdrills []entity.Drill
	for idx := 1; idx < len(drills); idx++ {
		bdrills := GenVDrillsBetween(drills[idx-1], drills[idx], 3, GenVDrillIDW)
		vdrills = append(vdrills, bdrills...)
	}
	utils.DrawDrills(vdrills, "./idw2.svg")
	utils.DisplayDrills(vdrills)
}
func TestGenVDrillsBetween(t *testing.T) {
	drill1 := constant.DrillSet()[0]
	drill2 := constant.DrillSet()[1]
	vdrills := GenVDrillsBetween(drill1, drill2, 5, GenVDrillM1)
	utils.DrawDrills(vdrills, "./between.svg")
}
func TestDrawDrills(t *testing.T) {
	drillNames := []string{"TZZK92", "TZJT31", "TZZK40", "TZJT28", "TZZK69", "TZZK70", "TZZK72"}
	var drills []entity.Drill
	for _, name := range drillNames {
		if drill, ok := constant.GetDrillByName(name); ok {
			drills = append(drills, drill)
		}
	}
	utils.DrawDrills(drills, "./realDrill.svg")
	utils.DisplayDrills(drills)
}

func TestGenHelpDrills(t *testing.T) {
	drills := GenHelpDrills()
	fmt.Println(len(drills))
}
func BenchmarkGenHelpDrills(b *testing.B) {
	drills := GenHelpDrills()
	fmt.Println(len(drills))
}