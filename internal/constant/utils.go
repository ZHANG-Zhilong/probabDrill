package constant

import (
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"math"
	"os"
	"probabDrill/apps/probDrill/model"
	"strconv"
)

func genIDWDrill(drills []model.Drill, x, y float64) model.Drill {
	x, y = decimal(x), decimal(y)
	log.SetFlags(log.Lshortfile)
	vdrill := drills[0].MakeDrill(GenVDrillName(), x, y, 0)

	nearDrills := vdrill.NearKDrills(drills, viper.GetInt("RadiusIn"))
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
			h := d.LayerHeights[lidx]
			w := d.GetWeight()
			vHeights[lidx] += w * h
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
func setClassicalIdwWeights(center model.Drill, nearDrills []model.Drill) (weights []float64) {
	log.SetFlags(log.Lshortfile)
	var (
		weightSum       float64
		hasZeroDistance bool
		zeroIdx         int
	)

	//get distance
	for idx, aroundDrill := range nearDrills {
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
			weights[idx] = 1e7 / math.Pow(weights[idx], viper.GetFloat64("IdwPow"))
			weightSum += weights[idx]
		}
		for idx, _ := range weights { //归一化, and set int the drill.
			weights[idx] = decimal(weights[idx] / weightSum)
			nearDrills[idx].SetWeight(weights[idx])
		}
		weightSum = 0
		for _, w := range weights {
			weightSum += w
		}
		weightSum = decimal(weightSum)
	}

	if math.Abs(weightSum-1) > 1e-1 {
		log.Printf("center drill: %#v\n", center)
		for _, d := range nearDrills {
			log.Printf("dist:%f, weight:%f, near drill: %#v\n", center.Distance(d), center.GetWeight(), d)
		}
		log.Fatalf("error, total weight is not 1, the total weight is:%f, weights is :%#v\n", weightSum, weights)
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

func SimpleDrillSet() (drills []model.Drill) {
	var drill1, drill2, drill3, drill4 model.Drill
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

	drills = []model.Drill{drill1, drill2, drill3, drill4}
	return
}
func DisplayDrills(drills []model.Drill) {
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
