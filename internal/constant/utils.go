package constant

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"probabDrill"
	"probabDrill/internal/entity"
	"strconv"
)

func genIDWDrill(drills []entity.Drill, x, y float64) (vdrill entity.Drill) {
	log.SetFlags(log.Lshortfile)
	vdrill = vdrill.MakeDrill(GenVDrillName(), x, y, 0)
	nearDrills := vdrill.NearKDrills(drills, probabDrill.RadiusIn)
	for _, d := range nearDrills { // if the position of the vdrill is just at a real drill's position
		if math.Abs(x-d.X) < 0.001 && math.Abs(y-d.Y) < 0.001 {
			return d
		}
	}
	setClassicalIdwWeights(vdrill, nearDrills)
	nearDrills = UnifyDrillsSeq(nearDrills, CheckSeqMinNeg)
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
	vdrill.UnBlock()
	if !vdrill.IsValid() {
		vdrill.Display()
		log.Fatal("invalid vdrill.\n")
	}
	return vdrill
}
func setClassicalIdwWeights(center entity.Drill, aroundDrills []entity.Drill) (weights []float64) {
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
func readFile(path string) string {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	return string(content)
}

func SimpleDrillSet() (drills []entity.Drill) {
	var drill1, drill2, drill3, drill4 entity.Drill
	drill1 = drill1.MakeDrill("1", 0, 0, 0)
	drill2 = drill1.MakeDrill("2", 1, 0, 0)
	drill3 = drill1.MakeDrill("3", 1, 1, 0)
	drill4 = drill1.MakeDrill("4", 0, 1, 0)

	drill1.AddLayerWithHeight(1, -1)
	drill1.AddLayerWithHeight(1, -2)
	drill1.AddLayerWithHeight(6, -3)
	drill1.AddLayerWithHeight(3, -4)

	drill2.AddLayerWithHeight(2, -1)
	drill2.AddLayerWithHeight(5, -2)
	drill2.AddLayerWithHeight(3, -3)
	drill2.AddLayerWithHeight(4, -4)

	drill3.AddLayerWithHeight(1, -1)
	drill3.AddLayerWithHeight(5, -2)
	drill3.AddLayerWithHeight(6, -3)
	drill3.AddLayerWithHeight(4, -4)

	drill4.AddLayerWithHeight(1, -1)
	drill4.AddLayerWithHeight(2, -2)
	drill4.AddLayerWithHeight(3, -3)
	drill4.AddLayerWithHeight(4, -4)

	drills = []entity.Drill{drill1, drill2, drill3, drill4}
	return
}
func DisplayDrills(drills []entity.Drill) {
	for _, d := range drills {
		d.Display()
	}
	fmt.Printf("total %d drills.", len(drills))
}
func decimal(value float64) float64 {
	value = math.Trunc(value*1e2+0.5) * 1e-2
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}
