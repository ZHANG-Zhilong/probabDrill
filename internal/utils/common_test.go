package utils

import (
	"fmt"
	"probabDrill/internal/constant"
	"probabDrill/internal/entity"
	"testing"
)

func TestDrill_UnifyStratum(t *testing.T) {
	drills := constant.SimpleDrillSet()
	DisplayDrills(drills)
	uniLayers := UnifyDrillsStrata(&drills, CheckSeqZiChun)
	DisplayDrills(drills)
	fmt.Println(uniLayers)
}
func TestUnifySeq(t *testing.T) {
	//seq1 := []int{0, 1, 2, 3, 4}
	//seq2 := []int{0, 1, 2, 3, 2, 4}
	//seq3 := []int{0, 1, 3, 2, 4}

	seq1 := []int{0, 1, 3, 6}
	seq2 := []int{0, 2, 5, 3}
	seq3 := []int{0, 1, 5, 6}
	seqs := [][]int{seq1, seq2, seq3}
	newLayer := seq1
	for idx := 1; idx < len(seqs); idx++ {
		newLayer = getUnifiedSeq(seqs[idx], newLayer, CheckSeqZiChun)
	}
	fmt.Println(newLayer)
}
func TestCheckSeq(t *testing.T) {
	seq1 := []int{0, 1, 2, 3, 4}
	seq2 := []int{0, 1, 2, 3, 2, 4}
	seq3 := []int{0, 1, 3, 2, 4}
	seq1 = CheckSeqMinNeg(seq1)
	fmt.Println(seq1)
	seq2 = CheckSeqMinNeg(seq2)
	fmt.Println(seq2)
	seq3 = CheckSeqMinNeg(seq3)
	fmt.Println(seq3)
}

func TestUnifyStratum(t *testing.T) {
	drill1 := entity.Drill{
		Layers:       []int{0, 1, 2, 3, 4},
		LayerHeights: []float64{0, -1, -2, -3, -4},
	}
	drill2 := entity.Drill{
		Layers:       []int{0, 1, 3, 2, 4},
		LayerHeights: []float64{0, -1, -2, -3, -4},
	}
	drill3 := entity.Drill{
		Layers:       []int{0, 1, 2, 3, 4},
		LayerHeights: []float64{0, -1, -2, -3, -4},
	}
	drills1 := []entity.Drill{drill1, drill2, drill3}
	drills1 = *UnifyDrillsStrata(&drills1, CheckSeqZiChun)
	fmt.Println("=======")
	DisplayDrills(drills1)
	drills1 = *UnifyDrillsStrata(&drills1, CheckSeqMinNeg)
	fmt.Println("=======")
	DisplayDrills(drills1)

	drill4 := entity.Drill{
		Layers:       []int{0, 1, 2, 3, 4},
		LayerHeights: []float64{0, -1, -2, -3, -4},
	}
	drill5 := entity.Drill{
		Layers:       []int{0, 1, 2, 3, 2, 4},
		LayerHeights: []float64{0, -1, -2, -3, -4, -5},
	}
	drill6 := entity.Drill{
		Layers:       []int{0, 1, 3, 2, 4},
		LayerHeights: []float64{0, -1, -2, -3, -4},
	}
	drills2 := []entity.Drill{drill4, drill5, drill6}
	drills2 = *UnifyDrillsStrata(&drills2, CheckSeqZiChun)
	DisplayDrills(drills2)
	fmt.Println("=======")
	drills2 = *UnifyDrillsStrata(&drills2, CheckSeqMinNeg)
	DisplayDrills(drills2)
	fmt.Println("=======")
	DrawDrills([]entity.Drill{constant.DrillSet()[1], constant.DrillSet()[2]})
}
