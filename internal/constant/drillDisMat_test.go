package constant

import (
	"fmt"
	"probabDrill/internal/entity"
	"testing"
)

func TestGetDrillsIn(t *testing.T) {
	drill200 := entity.NewBasicDrill("TZK10", 126.70100000000093, - 54.17859999999986, 4)
	drillIn := GetDrillsIn(*drill200, 2200)
	fmt.Println(len(drillIn))
	for _, d := range drillIn {
		fmt.Println(d)
		fmt.Println(getDist(d.X, d.Y, drill200.X, drill200.Y))
	}
}
func TestGetRadiusInclude(t *testing.T) {
	drill200 := entity.NewBasicDrill("TZK10", 126.70100000000093, - 54.17859999999986, 4)
	num := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for _, n := range num {
		if dist, ok := GetRadiusInclude(*drill200, n); ok {
			drillIn := GetDrillsIn(*drill200, dist)
			if n != len(drillIn) {
				t.Error(dist)
				fmt.Println(len(drillIn))
				for _, d := range drillIn {
					fmt.Println(d)
					fmt.Println(getDist(d.X, d.Y, drill200.X, drill200.Y))
				}
			}
		}
	}

}
