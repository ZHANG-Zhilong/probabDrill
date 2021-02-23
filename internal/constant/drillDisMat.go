package constant

import (
	"math"
	"probabDrill/internal/entity"
	"sort"
	"sync"
)

//the order of disMat is the same as drillSet
var disMat [][]float64
var nameIdxMap map[string]int

var once3 sync.Once

//get the drills around the param drill, not include drill itself, include boundary
func GetDrillsIn(drill entity.Drill, distance float64) (drills []entity.Drill) {
	init3()
	drillSet := DrillSet()
	drillIdx, ok := nameIdxMap[drill.Name]
	if ok {
		dists := disMat[drillIdx]
		for idx, dist := range dists {
			if dist <= distance && dist >= 0.001 {
				drills = append(drills, drillSet[idx])
			}
		}
		return drills
	}
	return drills
}

//get the radius that there are only the specified include num drills around the param drill.
func GetRadiusInclude(drill entity.Drill, include int) (radius float64, ok bool) {
	init3()
	if idx, ok := nameIdxMap[drill.Name]; ok {
		var dists []float64 = make([]float64, len(disMat[idx]))
		copy(dists, disMat[idx])
		sort.Float64s(dists)
		for idx2, d := range dists {
			if include == idx2 {
				radius = d
				return radius, true
			}
		}
	}
	return radius, false
}

func getDist(x1, y1, x2, y2 float64) (dist float64, ok bool) {
	dist = math.Sqrt((x2-x1)*(x2-x1) + (y2-y1)*(y2-y1))
	if math.IsNaN(dist) || math.IsInf(dist, 0) || dist < 0 {
		return -1, false
	}
	return dist, true
}

func init3() {
	once3.Do(func() {

		drills := DrillSet()
		nameIdxMap = make(map[string]int)

		for idx1 := 0; idx1 < len(drills); idx1++ {

			var dists []float64
			nameIdxMap[drills[idx1].Name] = idx1

			for idx2 := 0; idx2 < len(drills); idx2++ {
				if dist, ok := getDist(drills[idx1].X, drills[idx1].Y, drills[idx2].X, drills[idx2].Y); ok {
					dists = append(dists, dist)
				}
			}

			disMat = append(disMat, dists)
		}
	})
}
