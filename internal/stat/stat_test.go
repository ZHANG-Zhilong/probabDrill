package stat

import (
	"fmt"
	"gonum.org/v1/gonum/mat"
	"probabDrill/internal/constant"
	"probabDrill/internal/entity"
	"probabDrill/internal/utils"
	"testing"
)

func TestProbLayers(t *testing.T) {
	drillSet := constant.GetDrillSet()
	probs := ProbLayers(&drillSet)
	fmt.Println(probs)
}
func TestProbBlocks(t *testing.T) {
	drills := constant.GetDrillSet()
	blocks := utils.MakeBlocks(drills, 1)
	probs := ProbBlocks(&drills, &blocks)
	fmt.Println(probs)
}
func TestProbLBs(t *testing.T) {
	//drills := constant.GetDrillSet()
	//blockHeights := service.MakeBlocks(drills, 1)
	//mat := ProbLBs(&drills, blockHeights)
	//fmt.Println(mat)
	drills := constant.SimpleDrillSet()
	entity.DisplayDrills(drills)
	blockHeights := utils.MakeBlocks(drills, 0.5)

	pb := ProbBlocks(&drills, &blockHeights)
	//fa1 := mat.Formatted(pb, mat.Prefix("    "), mat.Squeeze())
	//fmt.Printf("with only non-zero values:\na = % v\n\n", fa1)
	fmt.Println("p(blocks)", pb)

	pl := ProbLayers(&drills)
	//fa2 := mat.Formatted(pl, mat.Prefix("    "), mat.Squeeze())
	//fmt.Printf("with only non-zero values:\na = % v\n\n", fa2)
	fmt.Println("p(layers)", pl)

	plbs := ProbLBs(&drills, &blockHeights)
	fa3 := mat.Formatted(plbs, mat.Prefix("    "), mat.Squeeze())
	fmt.Printf("with only non-zero values:\na = % v\n\n", fa3)

	pbls := ProbBLs(&drills, &blockHeights)
	fa4 := mat.Formatted(pbls, mat.Prefix("    "), mat.Squeeze())
	fmt.Printf("with only non-zero values:\na = % v\n\n", fa4)

	utils.DrawDrills(drills)
}


