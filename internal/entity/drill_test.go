package entity

import (
	"fmt"
	"probabDrill/internal/constant"
	"testing"
)

func TestDrill_NearDrills(t *testing.T) {
	drill := constant.GetDrillSet()[0]
	gotNears := drill.NearDrills(constant.GetDrillSet(), 3)
	for _, d := range gotNears {
		fmt.Println(d.Distance(drill), d)
	}
	for _, d := range constant.GetDrillSet() {
		fmt.Println(d.Distance(drill), d)
	}
}
