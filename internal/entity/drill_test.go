package entity

import (
	"fmt"
	"testing"
)

func TestDrill_NearDrills(t *testing.T) {
	drill := DrillSet()[0]
	gotNears := drill.NearDrills(DrillSet(), 3)
	for _, d := range gotNears {
		fmt.Println(d.Distance(drill), d)
	}
	for _, d := range DrillSet() {
		fmt.Println(d.Distance(drill), d)
	}
}
