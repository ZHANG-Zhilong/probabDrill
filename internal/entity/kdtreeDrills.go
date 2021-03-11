package entity

import (
	"github.com/fogleman/poissondisc"
	"github.com/kyroy/kdtree"
	"log"
	"math"
	"probabDrill"
	"sync"
)

var (
	onceKdDrill   sync.Once
	onceKdHelp    sync.Once
	kdtDrills     *kdtree.KDTree
	kdtHelpDrills *kdtree.KDTree
)

func GetNearKDrills(vdrill Drill, k int) (drills []Drill) {
	onceKdDrill.Do(initKdDrills)
	rst := kdtDrills.KNN(&vdrill, k)
	for _, d := range rst {
		if ds, ok := d.(*Drill); ok {
			drills = append(drills, *ds)
		}
	}
	return
}
func GetNearHelpDrills(vdrill Drill, k int) (drills []Drill) {
	onceKdHelp.Do(initKdHelpDrills)
	rst := kdtHelpDrills.KNN(&vdrill, k)
	for _, d := range rst {
		if ds, ok := d.(*Drill); ok {
			drills = append(drills, *ds)
		}
	}
	return
}
func initKdHelpDrills() () {
	var (
		x0, y0, x1, y1 float64
		helpDrills     []kdtree.Point
	)
	rdrills := DrillSet()
	x0, y0, x1, y1 = drillSet[0].GetRec(rdrills)
	points := poissondisc.Sample(x0, y0, x1, y1, probabDrill.MinDistance, probabDrill.MaxAttemptAdd, nil)
	for _, p := range points {
		idwDrill := genIDWDrills(drillSet, p.X, p.Y)
		helpDrills = append(helpDrills, &idwDrill)
	}
	kdtHelpDrills = kdtree.New(helpDrills)
}
func initKdDrills() () {
	rdrills := DrillSet()
	var kdDrills []kdtree.Point
	for _, d := range rdrills {
		kdDrills = append(kdDrills, &d)
	}
	kdtDrills = kdtree.New(kdDrills)
}

func genIDWDrills(drills []Drill, x, y float64) (vdrill Drill) {
	log.SetFlags(log.Lshortfile)
	vdrill = vdrill.MakeDrill(GenVDrillName(), x, y, 0)
	//nearDrills := vdrill.NearDrills(drills, probabDrill.RadiusIn)
	nearDrills := GetNearKDrills(vdrill, probabDrill.RadiusIn)
	for _, d := range nearDrills { // if the position of the vdrill is just at a real drill's position
		if math.Abs(x-d.X) < 0.001 && math.Abs(y-d.Y) < 0.001 {
			return d
		}
	}
	setClassicalIdwWeights(vdrill, nearDrills)
	nearDrills = UnifyDrillsStrata(nearDrills, CheckSeqMinNeg)
	vdrill.Layers = nearDrills[0].Layers
	var vHeights = make([]float64, len(vdrill.Layers), len(vdrill.Layers))
	for lidx, _ := range vdrill.Layers {
		for _, d := range nearDrills {
			vHeights[lidx] += decimal(d.GetWeight() * d.LayerHeights[lidx])
		}
		if lidx-1 >= 0 && math.Abs(vHeights[lidx]-vHeights[lidx-1]) < 10e-5 {
			vHeights[lidx] = vHeights[lidx-1]
		}
	}
	vdrill.LayerHeights = vHeights
	vdrill.Z = vHeights[0]
	vdrill.GetLength()
	vdrill.UnStdSeq()
	if !vdrill.IsValid() {
		vdrill.Print()
		log.Fatal("invalid vdrill.\n")
	}
	return vdrill
}
func setClassicalIdwWeights(center Drill, aroundDrills []Drill) (weights []float64) {
	log.SetFlags(log.Lshortfile)
	var (
		weightSum       float64
		hasZeroDistance bool
		zeroIdx         int
	)

	//get distance
	for idx, aroundDrill := range aroundDrills {
		dist := center.Distance(aroundDrill)
		weights = append(weights, dist) //as distance
		if dist < 1e-1 {
			hasZeroDistance = true
			zeroIdx = idx
		}
	}
	if hasZeroDistance {
		for idx, _ := range weights {
			weights[idx] = 0
		}
		if weights != nil && zeroIdx >= 0 && zeroIdx < len(weights) {
			weights[zeroIdx] = 1
		}
	} else {
		for idx, _ := range weights { //cal weight
			weights[idx] = 1e7 / math.Pow(weights[idx], probabDrill.IdwPow)
			weightSum += weights[idx]
		}
		for idx, _ := range weights { //归一化, and set int the drill.
			weights[idx] = decimal(weights[idx] / weightSum)
			aroundDrills[idx].SetWeight(weights[idx])
		}
		weightSum = 0
		for _, w := range weights {
			weightSum += w
		}
		weightSum = decimal(weightSum)
	}

	if math.Abs(weightSum-1) > 1e-1 {
		log.Println(weights)
		log.Fatalf("error, total weight:%f\n", weightSum)
	}
	return weights
}
