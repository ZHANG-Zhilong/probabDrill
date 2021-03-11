package play

import (
	"fmt"
	"github.com/kyroy/kdtree"
	"probabDrill/internal/constant"
	"probabDrill/internal/entity"
)

type Data struct {
	value string
}

func PlayKdTree() {
	drillSet := constant.GetDrillSet()
	var drills []kdtree.Point
	for _, d := range drillSet {
		drills = append(drills, &d)
	}

	tree := kdtree.New(drills)

	// Insert
	//tree.Insert(&points.Point2D{X: 1, Y: 8})
	//tree.Insert(&points.Point2D{X: 7, Y: 5})

	// KNN (k-nearest neighbor)
	var dp kdtree.Point
	dp = &drillSet[0]
	fmt.Println(dp)
	fmt.Println("--")
	rst := tree.KNN(dp, 3)
	for _, d := range rst {
		if ds, ok := d.(*entity.Drill); ok {
			ds.Print()
		}
	}
	// [{3.00 1.00} {5.00 0.00}]

	// RangeSearch
	//fmt.Println(tree.RangeSearch(kdrange.New(1, 8, 0, 2)))
	// [{5.00 0.00} {3.00 1.00}]

	// Points
	//fmt.Println(tree.Points())
	// [{3.00 1.00} {1.00 8.00} {5.00 0.00} {8.00 3.00} {7.00 5.00}]

	// Remove
	//fmt.Println(tree.Remove(&points.Point2D{X: 5, Y: 0}))
	// {5.00 0.00}

	// String
	//fmt.Println(tree)
	// [[{1.00 8.00} {3.00 1.00} [<nil> {8.00 3.00} {7.00 5.00}]]]

	// Balance
	//tree.Balance()
	//fmt.Println(tree)
	// [[[{3.00 1.00} {1.00 8.00} <nil>] {7.00 5.00} {8.00 3.00}]]
}
