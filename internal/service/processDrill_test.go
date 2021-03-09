package service

import (
	"fmt"
	"log"
	"probabDrill/internal/constant"
	"probabDrill/internal/entity"
	"probabDrill/internal/utils"
	"testing"
)

func TestDisplayDrills(t *testing.T) {
	drillSet := constant.DrillSet()
	utils.DisplayDrills(drillSet)
}
func TestStatProbBlockLayer(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	drills := constant.DrillSet()
	blocks := MakeBlocks(drills, constant.BlockResZ)
	for idx := int(1); idx < len(blocks); idx++ {
		p := utils.StatProbBlockLayer(drills, blocks[idx-1], blocks[idx], 2)
		log.Println(p)
	}
}
func TestGetGrid(t *testing.T) {
	var tests = []struct {
		px, py, l, r, t, b float64
	}{
		{1, 2, -5, 5, 5, -5},
		{0.2, 0.4, -1, 1, 1, -1},
		{0.23, 0.4, -1, 1, 1, -1},
	}
	for _, test := range tests {
		x, y := utils.GetGrids(test.px, test.py, test.l, test.r, test.t, test.b)
		t.Errorf("px %f, py %f, l %f, r %f, t %f, b %f", test.px, test.py, test.l, test.r, test.t, test.b)
		t.Errorf("gridx, %.2f~%.2f-%.2f,len %d", test.l, test.r, test.px, len(x))
		t.Errorf("%.2f\n", x)
		t.Errorf("gridy, %.2f~%.2f-%.2f,len %d", test.t, test.b, test.py, len(y))
		t.Errorf("%.2f\n", y)
	}
}
func TestGetRecByDrills(test *testing.T) {
	drills := constant.DrillSet()
	l, r, t, b := getDrillsRecXOY(drills)
	fmt.Printf("%.2f, %.2f, %.2f, %.2f", l, r, t, b)
}
func TestGetLayerSeq(t *testing.T) {
	drills := constant.DrillSet()
	blocks := MakeBlocks(drills, constant.BlockResZ)
	drill0 := drills[0]
	heights := explodedHeights(blocks, drill0.Z, drill0.GetBottomHeight())

	layers := []int{0}
	for idx := 1; idx < len(heights); idx++ {
		if seq, ok := drill0.GetLayerSeq(heights[idx-1], heights[idx]); ok {
			layers = append(layers, seq)
		}
	}
	fmt.Println(layers)
	utils.PrintFloat64s(heights)
	drill0.Print()
}
func TestExplodeDrill(t *testing.T) {
	drill := constant.DrillSet()[0]
	blocks := MakeBlocks(constant.DrillSet(), constant.BlockResZ)
	drill.Print()
	drill = drill.Explode(blocks)
	drill.Print()
}
func TestGenerateVirtualDrill(t *testing.T) {

	log.SetFlags(log.Lshortfile)
	drill := constant.DrillSet()[1]

	log.Println("real drill")
	drill.Print()

	blocks := MakeBlocks(constant.DrillSet(), constant.BlockResZ)
	virtualDrill := generateVirtualDrill(constant.DrillSet(), drill.X+1, drill.Y+1, blocks)

	log.Println("virtual drill")
	virtualDrill.Print()
}
func TestGenerateVirtualDrill2(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	drill1 := entity.Drill{
		Name: "1", X: 0, Y: 0, Z: 0,
		Layers:       []int{0, 1, 2, 3},
		LayerHeights: []float64{0, -1, -2, -3},
	}
	drill2 := entity.Drill{
		Name: "2", X: 1, Y: 0, Z: 0,
		Layers:       []int{0, 1, 2, 2},
		LayerHeights: []float64{0, -1, -2, -3},
	}
	drill3 := entity.Drill{
		Name: "3", X: 0, Y: 1, Z: 0,
		Layers:       []int{0, 1, 2, 3},
		LayerHeights: []float64{0, -1.2, -2.3, -3},
	}
	drill4 := entity.Drill{
		Name: "4", X: 1, Y: 1, Z: 0,
		Layers:       []int{0, 1, 2, 3},
		LayerHeights: []float64{0, -1.5, -2.3, -3},
	}
	drillSet := []entity.Drill{drill1, drill2, drill3, drill4}
	blocks := MakeBlocks(drillSet, 0.02)
	fmt.Println(blocks)
	var virtualDrills []entity.Drill
	for x := 0.0; x < 1; x += 0.1 {
		virtualDrills = append(virtualDrills, generateVirtualDrill(drillSet, x, 0.5, blocks))
	}
	for _, v := range virtualDrills {
		v.Print()
	}
	utils.DrawDrills(virtualDrills)
}
func TestGenerateVirtualDrill3(t *testing.T) {
	drillNames := []string{"TZZK92", "TZJT31", "TZZK40", "TZJT28", "TZZK69", "TZZK70", "TZZK72"}
	var drills []entity.Drill
	for _, name := range drillNames {
		if drill, ok := constant.GetDrillByName(name); ok {
			drills = append(drills, drill)
		}
	}
	var virtualDrills []entity.Drill
	for idx := 1; idx < len(drills); idx++ {
		vdrills := GetVirtualDrillsBetween(drills[idx-1], drills[idx], 3)
		virtualDrills = append(virtualDrills, vdrills...)
	}
	utils.DrawDrills(virtualDrills)
	//utils.DisplayDrills(virtualDrills)
}
func TestGetVirtualDrillsBetween(t *testing.T) {
	drill1 := constant.DrillSet()[0]
	drill2 := constant.DrillSet()[1]
	//vdrills := GetVirtualDrillsBetween(drill1, drill2, 5)
	vdrills := GetVirtualDrillsBetween(drill1, drill2, 5)
	fmt.Println(len(vdrills))
	utils.DrawDrills(vdrills)
}
func TestGenerateVirtualDrill4(t *testing.T) {
	drillNames := []string{"TZZK92", "TZJT31", "TZZK40", "TZJT28", "TZZK69", "TZZK70", "TZZK72"}
	var drills []entity.Drill
	for _, name := range drillNames {
		if drill, ok := constant.GetDrillByName(name); ok {
			drills = append(drills, drill)
		}
	}
	utils.DrawDrills(drills)
}
func TestStatLayer(t *testing.T) {
	//log.SetFlags(log.Lshortfile)
	//drill1, _ := constant.GetDrillByName("TZZK01")
	//drill2, _ := constant.GetDrillByName("TZJT02")
	//drill3, _ := constant.GetDrillByName("TZJT64")
	//drill4, _ := constant.GetDrillByName("TZJT72")
	//drills := []entity.Drill{drill1, drill2, drill3, drill4}
}
func TestStats(t *testing.T) {
	var virtualDrill entity.Drill
	log.SetFlags(log.Lshortfile)
	virtualDrill = virtualDrill.MakeDrill(constant.GenVirtualDrillName(), 5580, -15020, 0)
	blocks := MakeBlocks(constant.DrillSet(), constant.BlockResZ)
	incidentDrills := obtainNearDrills(constant.DrillSet(), virtualDrill, constant.RadiusIn)
	setClassicalIdwWeights(virtualDrill, incidentDrills)
	setLengthAndZ(&virtualDrill, incidentDrills)
	virtualDrill.LayerHeights = explodedHeights(blocks, virtualDrill.Z, virtualDrill.GetBottomHeight())

	var probBlockWithWeights = make([]float64, len(blocks), len(blocks))
	//var probBlocks = make([]float64, len(blocks), len(blocks))
	//var probBlockGeneral = make([]float64, len(blocks), len(blocks))
	for idx := 1; idx < len(blocks); idx++ {
		probBlockWithWeights[idx] = utils.StatProbBlockWithWeight(incidentDrills, blocks[idx-1], blocks[idx])
		//probBlocks[idx] = statProbBlock(incidentDrills, blocks[idx-1], blocks[idx])
		//probBlockGeneral[idx] = statProbBlock(constant.DrillSet(), blocks[idx-1], blocks[idx])
	}
	//log.Println("p(block)")
	//printFloat64s(blocks)
	//printFloat64s(probBlockWithWeights)
	//printFloat64s(probBlocks)
	//printFloat64s(probBlockGeneral)

	for idx := 1; idx < len(blocks); idx++ {
		var probLayers = make([]float64, constant.StdLen, constant.StdLen)
		var probBlockLayers = make([]float64, constant.StdLen, constant.StdLen)
		var probLayerBlock2s = make([]float64, constant.StdLen, constant.StdLen)
		for lidx := int(1); lidx < constant.StdLen; lidx++ {
			probLayers[lidx] = utils.StatProbLayerWithWeight(incidentDrills, blocks[idx-1], blocks[idx], lidx)
			probBlockLayers[lidx] = utils.StatProbBlockLayer(constant.DrillSet(), blocks[idx-1], blocks[idx], lidx)
			if probBlockWithWeights[idx] >= 0.0000001 {
				probLayerBlock2s[lidx] = probLayers[lidx] * probBlockLayers[lidx] / probBlockWithWeights[idx]
			}
		}
	}
}
func TestProbLayerAndBlocksWithWeight(t *testing.T) {
	blocks := MakeBlocks(constant.DrillSet(), constant.BlockResZ)
	drills := constant.DrillSet()
	drill0 := drills[0]
	incidentDrills := obtainNearDrills(constant.DrillSet(), drill0, 50)
	setClassicalIdwWeights(drill0, incidentDrills)
	log.SetFlags(log.Lshortfile)
	for idx := 1; idx < len(blocks); idx++ {
		for layer := int(1); layer < constant.StdLen; layer++ {
			ceil, floor := blocks[idx-1], blocks[idx]
			prob1 := utils.StatProbBlockAndLayerWithWeight(incidentDrills, ceil, floor, layer)
			prob2 := utils.StatProbBlockLayer(drills, ceil, floor, layer)
			log.Printf("%f, %f\n", prob1, prob2)
		}
	}
}
func TestStatProbBlockWithWeight(t *testing.T) {
	incidentDrillSet := obtainNearDrills(constant.DrillSet(), constant.DrillSet()[0], 3)
	setClassicalIdwWeights(constant.DrillSet()[0], incidentDrillSet)
	blocks := MakeBlocks(constant.DrillSet(), constant.BlockResZ)
	for _, d := range incidentDrillSet {
		d.Print()
	}

	probBlock := utils.StatProbBlockWithWeight(incidentDrillSet, blocks[1], blocks[2])
	fmt.Println(probBlock)

}
func TestExplodeHeights(t *testing.T) {
	drill := constant.DrillSet()[0]
	blocks := MakeBlocks(constant.DrillSet(), constant.BlockResZ)
	heights := explodedHeights(blocks, drill.Z, drill.GetBottomHeight())
	drill = drill.Explode(blocks)
	if len(heights) == len(drill.LayerHeights) {
		for idx, _ := range heights {
			if heights[idx] != drill.LayerHeights[idx] {
				t.Error("error")
			}
		}
	}
}
func TestClassicalIdw(t *testing.T) {
	var drill1, drill2, drill3 entity.Drill
	drill1 = drill1.MakeDrill("1", 0, 0, 0)
	drill2 = drill2.MakeDrill("1", 2, 0, 0)
	drill3 = drill3.MakeDrill("1", 3, 0, 0)
	drills := []entity.Drill{drill2, drill3}
	weights := setClassicalIdwWeights(drill1, drills)
	fmt.Println(drills)
	fmt.Println(weights)
}
func BenchmarkGetGrid(b *testing.B) {
	for i := 0; i < b.N; i++ {
		utils.GetGrids(1, 1, -5, 5, 5, -5)
	}
}
func TestFindMaxValue(t *testing.T) {
	s := []float64{0, 3, -2}
	idx, val := utils.FindMaxFloat64s(s)
	fmt.Println(idx, val)
}
func TestIsInPolygon(t *testing.T) {
	vertx := []float64{0, 0, 1, 1.5, 1}
	verty := []float64{0, 1, 1, 0.5, 0}
	testx := []float64{0.5, 1, 0, -1, 1, 0.25, 0}
	testy := []float64{1, 1, 0, 1, 1.5, 1.25, 1.5}
	rst := []bool{true, true, true, false, false, false, false}
	for idx, _ := range rst {
		if rst[idx] != utils.IsInPolygon(vertx, verty, testx[idx], testy[idx]) {
			t.Error("error")
			t.Error(testx[idx], testy[idx], rst[idx])
		}
	}

	x, y := constant.GetBoundary()
	l, r, top, b := getDrillsRecXOY(constant.DrillSet())
	gridx, gridy := utils.GetGrids(constant.GridXY, constant.GridXY, l, r, top, b)
	var in, notin int
	for _, val1 := range gridx {
		for _, val2 := range gridy {
			if utils.IsInPolygon(x, y, val1, val2) {
				in++
			} else {
				notin++
			}
		}
	}
	log.Println("in", in, "not in", notin)

}
