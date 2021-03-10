package utils

import (
	"fmt"
	"log"
	"probabDrill"
	"probabDrill/internal/constant"
	"probabDrill/internal/entity"
	"testing"
)

func TestDrill_UnifyStratum(t *testing.T) {
	drills := constant.SimpleDrillSet()
	DisplayDrills(drills)
	uniLayers := UnifyDrillsStrata(drills, CheckSeqZiChun)
	DisplayDrills(drills)
	fmt.Println(uniLayers)
}
func TestUnifySeq(t *testing.T) {
	//seq1 := []int{0, 1, 2, 3, 4}
	//seq2 := []int{0, 1, 2, 3, 2, 4}
	//seq3 := []int{0, 1, 3, 2, 4}

	seq1 := []int{0, 1, 3, 6}
	seq2 := []int{0, 2, 5, 3}
	seq3 := []int{0, 1, 5, 6}
	seqs := [][]int{seq1, seq2, seq3}
	newLayer := seq1
	for idx := 1; idx < len(seqs); idx++ {
		newLayer = getUnifiedSeq(seqs[idx], newLayer, CheckSeqZiChun)
	}
	fmt.Println(newLayer)
}
func TestCheckSeq(t *testing.T) {
	seq1 := []int{0, 1, 2, 3, 4}
	seq2 := []int{0, 1, 2, 3, 2, 4}
	seq3 := []int{0, 1, 3, 2, 4}
	seq1 = CheckSeqMinNeg(seq1)
	fmt.Println(seq1)
	seq2 = CheckSeqMinNeg(seq2)
	fmt.Println(seq2)
	seq3 = CheckSeqMinNeg(seq3)
	fmt.Println(seq3)
}

func TestUnifyStratum(t *testing.T) {
	drill1 := entity.Drill{
		Layers:       []int{0, 1, 2, 3, 4},
		LayerHeights: []float64{0, -1, -2, -3, -4},
	}
	drill2 := entity.Drill{
		Layers:       []int{0, 1, 3, 2, 4},
		LayerHeights: []float64{0, -1, -2, -3, -4},
	}
	drill3 := entity.Drill{
		Layers:       []int{0, 1, 2, 3, 4},
		LayerHeights: []float64{0, -1, -2, -3, -4},
	}
	drills1 := []entity.Drill{drill1, drill2, drill3}
	drills1 = UnifyDrillsStrata(drills1, CheckSeqZiChun)
	fmt.Println("=======")
	DisplayDrills(drills1)
	drills1 = UnifyDrillsStrata(drills1, CheckSeqMinNeg)
	fmt.Println("=======")
	DisplayDrills(drills1)

	drill4 := entity.Drill{
		Layers:       []int{0, 1, 2, 3, 4},
		LayerHeights: []float64{0, -1, -2, -3, -4},
	}
	drill5 := entity.Drill{
		Layers:       []int{0, 1, 2, 3, 2, 4},
		LayerHeights: []float64{0, -1, -2, -3, -4, -5},
	}
	drill6 := entity.Drill{
		Layers:       []int{0, 1, 3, 2, 4},
		LayerHeights: []float64{0, -1, -2, -3, -4},
	}
	drills2 := []entity.Drill{drill4, drill5, drill6}
	drills2 = UnifyDrillsStrata(drills2, CheckSeqZiChun)
	DisplayDrills(drills2)
	fmt.Println("=======")
	drills2 = UnifyDrillsStrata(drills2, CheckSeqMinNeg)
	DisplayDrills(drills2)
	fmt.Println("=======")
	DrawDrills([]entity.Drill{constant.DrillSet()[1], constant.DrillSet()[2]}, "./a.svg")
}
func TestDecimal(t *testing.T) {
	sample := []float64{1.345, 1.00000000001}
	for _, s := range sample {
		fmt.Println(Decimal(s))
	}
}

func TestDisplayDrills(t *testing.T) {
	drillSet := constant.DrillSet()
	DisplayDrills(drillSet)
}
func TestStatProbBlockLayer(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	drills := constant.DrillSet()
	blocks := MakeBlocks(drills, probabDrill.BlockResZ)
	for idx := int(1); idx < len(blocks); idx++ {
		p := StatProbBlockLayer(drills, blocks[idx-1], blocks[idx], 2)
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
		x, y := GetGrids(test.px, test.py, test.l, test.r, test.t, test.b)
		t.Errorf("px %f, py %f, l %f, r %f, t %f, b %f", test.px, test.py, test.l, test.r, test.t, test.b)
		t.Errorf("gridx, %.2f~%.2f-%.2f,len %d", test.l, test.r, test.px, len(x))
		t.Errorf("%.2f\n", x)
		t.Errorf("gridy, %.2f~%.2f-%.2f,len %d", test.t, test.b, test.py, len(y))
		t.Errorf("%.2f\n", y)
	}
}
func TestGetLayerSeq(t *testing.T) {
	drills := constant.DrillSet()
	blocks := MakeBlocks(drills, probabDrill.BlockResZ)
	drill0 := drills[0]
	heights := ExplodedHeights(blocks, drill0.Z, drill0.GetBottomHeight())

	layers := []int{0}
	for idx := 1; idx < len(heights); idx++ {
		if seq, ok := drill0.GetLayerSeq(heights[idx-1], heights[idx]); ok {
			layers = append(layers, seq)
		}
	}
	fmt.Println(layers)
	PrintFloat64s(heights)
	drill0.Print()
}
func TestExplodeDrill(t *testing.T) {
	drill := constant.DrillSet()[0]
	blocks := MakeBlocks(constant.DrillSet(), probabDrill.BlockResZ)
	drill.Print()
	drill = drill.Explode(blocks)
	drill.Print()
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
	virtualDrill = virtualDrill.MakeDrill(constant.GenVDrillName(), 5580, -15020, 0)
	blocks := MakeBlocks(constant.DrillSet(), probabDrill.BlockResZ)
	drillSet := constant.DrillSet()
	nearDrills := virtualDrill.NearDrills(drillSet, probabDrill.RadiusIn)
	SetClassicalIdwWeights(virtualDrill, nearDrills)
	entity.SetLengthAndZ(&virtualDrill, nearDrills)
	virtualDrill.LayerHeights = ExplodedHeights(blocks, virtualDrill.Z, virtualDrill.GetBottomHeight())

	var probBlockWithWeights = make([]float64, len(blocks), len(blocks))
	//var probBlocks = make([]float64, len(blocks), len(blocks))
	//var probBlockGeneral = make([]float64, len(blocks), len(blocks))
	for idx := 1; idx < len(blocks); idx++ {
		probBlockWithWeights[idx] = StatProbBlockWithWeight(nearDrills, blocks[idx-1], blocks[idx])
		//probBlocks[idx] = statProbBlock(NearDrills, blocks[idx-1], blocks[idx])
		//probBlockGeneral[idx] = statProbBlock(constant.DrillSet(), blocks[idx-1], blocks[idx])
	}
	//log.Println("p(block)")
	//printFloat64s(blocks)
	//printFloat64s(probBlockWithWeights)
	//printFloat64s(probBlocks)
	//printFloat64s(probBlockGeneral)

	for idx := 1; idx < len(blocks); idx++ {
		var probLayers = make([]float64, probabDrill.StdLen, probabDrill.StdLen)
		var probBlockLayers = make([]float64, probabDrill.StdLen, probabDrill.StdLen)
		var probLayerBlock2s = make([]float64, probabDrill.StdLen, probabDrill.StdLen)
		for lidx := int(1); lidx < probabDrill.StdLen; lidx++ {
			probLayers[lidx] = StatProbLayerWithWeight(nearDrills, blocks[idx-1], blocks[idx], lidx)
			probBlockLayers[lidx] = StatProbBlockLayer(constant.DrillSet(), blocks[idx-1], blocks[idx], lidx)
			if probBlockWithWeights[idx] >= 0.0000001 {
				probLayerBlock2s[lidx] = probLayers[lidx] * probBlockLayers[lidx] / probBlockWithWeights[idx]
			}
		}
	}
}
func TestProbLayerAndBlocksWithWeight(t *testing.T) {
	blocks := MakeBlocks(constant.DrillSet(), probabDrill.BlockResZ)
	drills := constant.DrillSet()
	drill0 := drills[0]
	drillSet := constant.DrillSet()
	nearDrills := drill0.NearDrills(drillSet, 50)
	SetClassicalIdwWeights(drill0, nearDrills)
	log.SetFlags(log.Lshortfile)
	for idx := 1; idx < len(blocks); idx++ {
		for layer := int(1); layer < probabDrill.StdLen; layer++ {
			ceil, floor := blocks[idx-1], blocks[idx]
			prob1 := StatProbBlockAndLayerWithWeight(nearDrills, ceil, floor, layer)
			prob2 := StatProbBlockLayer(drills, ceil, floor, layer)
			log.Printf("%f, %f\n", prob1, prob2)
		}
	}
}
func TestStatProbBlockWithWeight(t *testing.T) {
	drillSet := constant.DrillSet()
	nearDrills := constant.DrillSet()[0].NearDrills(drillSet, 3)
	SetClassicalIdwWeights(constant.DrillSet()[0], nearDrills)
	blocks := MakeBlocks(constant.DrillSet(), probabDrill.BlockResZ)
	for _, d := range nearDrills {
		d.Print()
	}

	probBlock := StatProbBlockWithWeight(nearDrills, blocks[1], blocks[2])
	fmt.Println(probBlock)

}
func TestExplodeHeights(t *testing.T) {
	drill := constant.DrillSet()[0]
	blocks := MakeBlocks(constant.DrillSet(), probabDrill.BlockResZ)
	heights := ExplodedHeights(blocks, drill.Z, drill.GetBottomHeight())
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
	weights := SetClassicalIdwWeights(drill1, drills)
	fmt.Println(drills)
	fmt.Println(weights)
}
func BenchmarkGetGrid(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetGrids(1, 1, -5, 5, 5, -5)
	}
}
func TestFindMaxValue(t *testing.T) {
	s := []float64{0, 3, -2}
	idx, val := FindMaxFloat64s(s)
	fmt.Println(idx, val)
}
func TestIsInPolygon(t *testing.T) {
	vertx := []float64{0, 0, 1, 1.5, 1}
	verty := []float64{0, 1, 1, 0.5, 0}
	testx := []float64{0.5, 1, 0, -1, 1, 0.25, 0}
	testy := []float64{1, 1, 0, 1, 1.5, 1.25, 1.5}
	rst := []bool{true, true, true, false, false, false, false}
	for idx, _ := range rst {
		if rst[idx] != IsInPolygon(vertx, verty, testx[idx], testy[idx]) {
			t.Error("error")
			t.Error(testx[idx], testy[idx], rst[idx])
		}
	}

	x, y := constant.GetBoundary()
	l, r, top, b := constant.DrillSet()[0].GetRec(constant.DrillSet())
	gridx, gridy := GetGrids(probabDrill.GridXY, probabDrill.GridXY, l, r, top, b)
	var in, notin int
	for _, val1 := range gridx {
		for _, val2 := range gridy {
			if IsInPolygon(x, y, val1, val2) {
				in++
			} else {
				notin++
			}
		}
	}
	log.Println("in", in, "not in", notin)

}
