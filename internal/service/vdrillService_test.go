package service

import (
	"fmt"
	"probabDrill"
	"probabDrill/internal/constant"
	"probabDrill/internal/entity"
	"probabDrill/internal/utils"
	"testing"
)

func TestGenVDrillsM1(t *testing.T) {
	drillNames := []string{"TZZK92", "TZJT31", "TZZK40", "TZJT28", "TZZK69", "TZZK70", "TZZK72"}
	var realDrills []entity.Drill
	for _, name := range drillNames {
		if drill, ok := constant.GetRealDrillByName(name); ok {
			realDrills = append(realDrills, drill)
		}
	}
	drillSet := constant.GetHelpDrills()
	blocks := utils.MakeBlocks(drillSet, probabDrill.BlockResZ)
	pBlockLayerMat, _ := utils.ProbBlockLayerMatG(drillSet, blocks)
	var vdrills []entity.Drill
	for idx := 1; idx < len(realDrills); idx++ {
		middleVertices := utils.MiddleKPoints(realDrills[idx-1].X, realDrills[idx-1].Y, realDrills[idx].X, realDrills[idx].Y, 3)
		for idx2 := 1; idx2 < len(middleVertices); idx2 += 2 {
			vdrills = append(vdrills, GenVDrillM1(drillSet, blocks, pBlockLayerMat, middleVertices[idx2-1], middleVertices[idx2]))
		}
		vdrills = append(vdrills, realDrills[idx])
	}
	//extend drills by die zhi principle.
	unifiedSeq := constant.GetUnifiedSeq(vdrills, constant.CheckSeqZiChun)
	vdrills = utils.ExtendDrills(unifiedSeq, vdrills)
	//trunc drills.
	var avgBot float64
	for _, d := range vdrills {
		avgBot += d.BottomHeight()
	}
	avgBot = avgBot / float64(len(vdrills))
	for idx, _ := range vdrills {
		vdrills[idx] = vdrills[idx].Trunc(avgBot)
	}
	utils.DrawDrills(vdrills, "./TestGenVDrillsM1.svg")
}
func TestGenVDrillsIDW(t *testing.T) {
	drillNames := []string{"TZZK92", "TZJT31", "TZZK40", "TZJT28", "TZZK69", "TZZK70", "TZZK72"}
	var sampleDrills []entity.Drill
	for _, name := range drillNames {
		if drill, ok := constant.GetRealDrillByName(name); ok {
			sampleDrills = append(sampleDrills, drill)
		}
	}
	var vdrills []entity.Drill
	realDrills := constant.GetHelpDrills()
	//realDrills := constant.GetRealDrills()
	for idx := 1; idx < len(sampleDrills); idx++ {
		middleDrills := GenVDrillsBetween(realDrills, sampleDrills[idx-1], sampleDrills[idx], 3, GenVDrillsIDW)
		vdrills = append(vdrills, middleDrills...)
		vdrills = append(vdrills, sampleDrills[idx])
	}

	//extend drills by die zhi principle.
	unifiedSeq := constant.GetUnifiedSeq(vdrills, constant.CheckSeqZiChun)
	vdrills = utils.ExtendDrills(unifiedSeq, vdrills)
	//trunc drills.
	var avgBot float64
	for _, d := range vdrills {
		avgBot += d.BottomHeight()
	}
	avgBot = avgBot / float64(len(vdrills))
	for idx, _ := range vdrills {
		vdrills[idx] = vdrills[idx].Trunc(avgBot)
	}
	utils.DrawDrills(vdrills, "./TestGenVDrillsIDW.svg")
	//entity.DisplayDrills(vdrills)
}
func TestGenVDrillsBetween(t *testing.T) {
	drill1 := constant.GetRealDrills()[0]
	drill2 := constant.GetRealDrills()[1]
	drillSet := constant.GetRealDrills()
	vdrills := GenVDrillsBetween(drillSet, drill1, drill2, 5, GenVDrillsM1)
	utils.DrawDrills(vdrills, "./between.svg")
}
func TestDrawDrills(t *testing.T) {
	drillNames := []string{"TZZK92", "TZJT31", "TZZK40", "TZJT28", "TZZK69", "TZZK70", "TZZK72"}
	var drills []entity.Drill
	for _, name := range drillNames {
		if drill, ok := constant.GetRealDrillByName(name); ok {
			drills = append(drills, drill)
		}
	}
	utils.DrawDrills(drills, "./realDrill.svg")
	constant.DisplayDrills(drills)
}
func TestGenHelpDrills(t *testing.T) {
	drills := GenHelpDrills()
	fmt.Println(len(drills))
}
func BenchmarkGenHelpDrills(b *testing.B) {
	drills := GenHelpDrills()
	fmt.Println(len(drills))
}
