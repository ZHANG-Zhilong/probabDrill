package utils

import (
	"fmt"
	"probabDrill/internal/constant"
	"testing"
)

func TestDrill_UnifyStratum(t *testing.T) {
	drills := constant.SimpleDrillSet()
	DisplayDrills(drills)
	uniLayers := UnifyStratum(&drills)
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
		newLayer = UnifySeq(seqs[idx], newLayer)
	}
	fmt.Println(newLayer)
}
func TestCheckSeq(t *testing.T) {
	seq1 := []int{0, 1, 2, 3, 4}
	seq2 := []int{0, 1, 2, 3, 2, 4}
	seq3 := []int{0, 1, 3, 2, 4}
	CheckSeq(&seq1)
	fmt.Println(seq1)
	CheckSeq(&seq2)
	fmt.Println(seq2)
	CheckSeq(&seq3)
	fmt.Println(seq3)
}
