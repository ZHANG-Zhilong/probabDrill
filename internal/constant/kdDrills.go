package constant

import (
	"fmt"
	"github.com/kyroy/kdtree"
	"probabDrill/internal/entity"
	"sync"
)

var (
	onceKdDrill sync.Once
	kdtDrills   *kdtree.KDTree
)

func NearKDrills(drill entity.Drill, k int) (drills []entity.Drill) {
	onceKdDrill.Do(initKdDrills)
	rst := kdtDrills.KNN(&drill, k)
	for _, d := range rst {
		if ds, ok := d.(*entity.Drill); ok {
			drills = append(drills, *ds)
		}
	}
	return drills
}
func initKdDrills() () {
	rdrills := GetDrillSet()
	var kdDrills []kdtree.Point
	for _, d := range rdrills {
		kdDrills = append(kdDrills, &d)
	}
	kdtDrills = kdtree.New(kdDrills)
}

var (
	onceKdHelp    sync.Once
	kdtHelpDrills *kdtree.KDTree
)

func NearHelpDrills(drill entity.Drill, k int) (drills []entity.Drill) {
	onceKdHelp.Do(initKdHelpDrills)
	rst := kdtHelpDrills.KNN(&drill, k)
	for _, d := range rst {
		if ds, ok := d.(*entity.Drill); ok {
			drills = append(drills, *ds)
		}
	}
	return
}
func initKdHelpDrills() () {
	var helpDrills []kdtree.Point
	var helpDrillSet = GetHelpDrillSet()
	for _, d := range helpDrillSet {
		helpDrills = append(helpDrills, &d)
	}
	kdtHelpDrills = kdtree.New(helpDrills)
}

func Demo() {
	rdrills := GetDrillSet()[1:7]
	var kdDrills []kdtree.Point
	for _, d := range rdrills {
		kdDrills = append(kdDrills, &d)
	}
	kdtDrills = kdtree.New(kdDrills)

	var n1 []entity.Drill
	drill := entity.Drill{X: 1, Y: 2, Z: 3}
	var tar kdtree.Point = &drill
	rst := kdtDrills.KNN(tar, 2)
	kdtDrills.Balance()
	for _, d := range rst {
		if ds, ok := d.(*entity.Drill); ok {
			n1 = append(n1, *ds)
		}
	}
	n2 := drill.NearDrills(rdrills, 2)
	fmt.Println(n1)
	fmt.Println(n2)
}
