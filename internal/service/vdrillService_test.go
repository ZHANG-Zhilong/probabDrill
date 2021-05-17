package service

import (
	"fmt"
	"github.com/spf13/viper"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"
	"log"
	"probabDrill/internal/constant"
	"probabDrill/internal/entity"
	"probabDrill/internal/utils"
	"testing"
)

func TestPrintFigure(t *testing.T) {
	//for diagram  P(block) 和其他等在插值过程中用到的一些条形图
	drillSet := constant.GetHelpDrills()
	blocks := utils.MakeBlocks(drillSet, viper.GetFloat64("BlockResz"))
	//pLayerBlock, _ := utils.ProbLayerBlockMat(drillSet, blocks)
	pblocks, _ := utils.ProbBlocks(drillSet, blocks)
	fa := mat.Formatted(pblocks, mat.Prefix(""), mat.Squeeze())
	fmt.Printf("with all values:\na = %v\n\n", fa)
	//fmt.Printf("with only non-zero values:\na = % v\n\n", fa)

}

//use P(layer_ij|block_i) first max probability to fill the block of the drill.
func TestGenVDrillsM1(t *testing.T) {
	drillNames := []string{"TZZK92", "TZJT31", "TZZK40", "TZJT28", "TZZK69", "TZZK70", "TZZK72"}
	var realDrills []entity.Drill
	for _, name := range drillNames {
		if drill, ok := constant.GetRealDrillByName(name); ok {
			realDrills = append(realDrills, drill)
		}
	}
	drillSet := constant.GetHelpDrills()
	var vdrills []entity.Drill
	if realDrills != nil {
		vdrills = append(vdrills, realDrills[0])
	}
	for idx := 1; idx < len(realDrills); idx++ {
		midDrills := GenVDrillsBetween(drillSet, realDrills[idx-1], realDrills[idx], 1, GenVDrillsM1)
		vdrills = append(vdrills, midDrills...)
		vdrills = append(vdrills, realDrills[idx])
	}
	//extend drills by die zhi principle.
	unifiedSeq := constant.GetUnifiedSeq(realDrills, constant.CheckSeqZiChun)
	vdrills = utils.ExtendDrills(unifiedSeq, vdrills)

	//trunc drills.
	//utils.TruncDrills(vdrills)

	//draw drills by WXD svg.
	utils.Drill2WXD(vdrills)

	//draw drills by svg in simple form.
	utils.DrawDrills(vdrills, "./TestGenVDrillsM1.svg")
}

//use P(layer_ij|block_i) second to fill the block of the drill.
func TestGenVDrillsM1Second(t *testing.T) {
	drillNames := []string{"TZZK92", "TZJT31", "TZZK40", "TZJT28", "TZZK69", "TZZK70", "TZZK72"}
	var realDrills []entity.Drill
	for _, name := range drillNames {
		if drill, ok := constant.GetRealDrillByName(name); ok {
			realDrills = append(realDrills, drill)
		}
	}
	drillSet := constant.GetHelpDrills()
	var vdrills []entity.Drill
	if realDrills != nil {
		vdrills = append(vdrills, realDrills[0])
	}
	for idx := 1; idx < len(realDrills); idx++ {
		midDrills := GenVDrillsBetween(drillSet, realDrills[idx-1], realDrills[idx], 1, GenVDrillsM1Second)
		vdrills = append(vdrills, midDrills...)
		vdrills = append(vdrills, realDrills[idx])
	}
	//extend drills by die zhi principle.
	unifiedSeq := constant.GetUnifiedSeq(realDrills, constant.CheckSeqZiChun)
	vdrills = utils.ExtendDrills(unifiedSeq, vdrills)

	//trunc drills.
	//utils.TruncDrills(vdrills)

	//draw drills by WXD svg.
	utils.Drill2WXD(vdrills)

	//draw drills by svg in simple form.
	utils.DrawDrills(vdrills, "./TestGenVDrillsM1.svg")
}

//generate 3d geology model by interpolation method idw.
func TestGenVDrillsIDW(t *testing.T) {
	drillNames := []string{"TZZK92", "TZJT31", "TZZK40", "TZJT28", "TZZK69", "TZZK70", "TZZK72"}
	var realDrills []entity.Drill
	for _, name := range drillNames {
		if drill, ok := constant.GetRealDrillByName(name); ok {
			realDrills = append(realDrills, drill)
		}
	}
	var vdrills []entity.Drill
	//drillSet := constant.GetHelpDrills()
	drillSet := constant.GetRealDrills()
	if realDrills != nil {
		vdrills = append(vdrills, realDrills[0])
	}
	for idx := 1; idx < len(realDrills); idx++ {
		middleDrills := GenVDrillsBetween(drillSet, realDrills[idx-1], realDrills[idx], 1, GenVDrillsIDW)
		vdrills = append(vdrills, middleDrills...)
		vdrills = append(vdrills, realDrills[idx])
	}
	//extend drills by die zhi principle.
	unifiedSeq := constant.GetUnifiedSeq(realDrills, constant.CheckSeqZiChun)
	vdrills = utils.ExtendDrills(unifiedSeq, vdrills)
	//trunc drills.
	//utils.TruncDrills(vdrills)
	fmt.Println(unifiedSeq)
	utils.Drill2WXD(vdrills)
	utils.DrawDrills(vdrills, "./TestGenVDrillsIDW.svg")
}

func TestGenVDrillsBetween(t *testing.T) {
	drill1 := constant.GetRealDrills()[0]
	drill2 := constant.GetRealDrills()[1]
	drillSet := constant.GetRealDrills()
	vdrills := GenVDrillsBetween(drillSet, drill1, drill2, 5, GenVDrillsM1)
	utils.DrawDrills(vdrills, "./between.svg")
}

func TestDrawDrills(t *testing.T) {
	//只使用真实的钻孔数据绘制剖面图
	drillNames := []string{"TZZK92", "TZJT31", "TZZK40", "TZJT28", "TZZK69", "TZZK70", "TZZK72"}
	var drills []entity.Drill
	for _, name := range drillNames {
		if drill, ok := constant.GetRealDrillByName(name); ok {
			drills = append(drills, drill)
		}
	}
	utils.DrawDrills(drills, "./realDrill.svg")
	constant.DisplayDrills(drills)
}

func TestRealDrillsIDW(t *testing.T) {
	//drillNames := []string{"TZZK92", "TZJT31", "TZZK40", "TZJT28", "TZZK69", "TZZK70", "TZZK72"}
	drillNames := []string{"BP01", "BP02"}
	var realDrills []entity.Drill
	for _, name := range drillNames {
		if drill, ok := constant.GetRealDrillByName(name); ok {
			realDrills = append(realDrills, drill)
		}
	}
	//extend drills by die zhi principle.
	//unifiedSeq := constant.GetUnifiedSeq(realDrills, constant.CheckSeqZiChun)
	//realDrills = utils.ExtendDrills(unifiedSeq, realDrills)
	utils.Drill2WXD(realDrills)
	for _, d := range realDrills {
		d.Display()
	}
	utils.DrawDrills(realDrills, "./fating-real.svg")
}

func TestGetPeAvgByLayer(t *testing.T) {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	//按地质分层，分别计算每一个分层的平均百分比误差
	drills := constant.GetHelpDrills()
	//drills := constant.GetRealDrills()
	if pes, err := GetPeAvgByLayer(drills, 10, GenVDrillsM1); err == nil {
		log.Println("M1:", utils.Average(pes), pes)
	} else {
		log.Fatal("TestGetPeAvgByLayer", err)
	}

	if pes, err := GetPeAvgByLayer(drills, 10, GenVDrillsIDW); err == nil {
		log.Println("IDW:", utils.Average(pes), pes)
	} else {
		log.Fatal("TestGetPeAvgByLayer", err)
	}
}
func TestGetPeCloudFigure(t *testing.T) {
	//通过多次，计算，获取研究区域内不同位置虚拟钻孔Pe
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	drills := constant.GetHelpDrills()
	if x, y, pe, err := GetPeByArea(drills, 10, GenVDrillsIDW, 10); err == nil {
		log.Printf("共统计个%d数据点\n", len(x))
		fmt.Print("x=")
		fmt.Print(x)
		fmt.Println(";")

		fmt.Print("y=")
		fmt.Print(y)
		fmt.Println(";")

		fmt.Print("z=")
		fmt.Print(pe)
		fmt.Println(";")

		x1, y1, x2, y2 := drills[0].GetRec(constant.GetRealDrills())
		fmt.Printf("[X,Y,Z]=griddata(x,y,z,linspace(%f,%f,100)',linspace(%f,%f,100),'v4');\n", x1, x2, y1, y2)
		log.Printf("max pes:%f, min pes: %f, avg pes: %f\n", floats.Max(pe), floats.Min(pe), utils.Average(pe))
	} else {
		log.Fatal("TestGetPeAvgByLayer", err)
	}

}

func TestDrillAroundPeCloudFigure(t *testing.T) {
	//真实钻孔周围虚拟钻孔的百分比误差
	drillSet := constant.GetHelpDrills()
	drillSet = append(drillSet, constant.GetRealDrills()...)
	blocks := utils.MakeBlocks(drillSet, viper.GetFloat64("blockResZ"))
	pBlockLayerMat, _ := utils.ProbBlockLayerMatG(drillSet, blocks)
	drillNames := []string{"TZZK92", "TZJT31", "TZZK40", "TZJT28", "TZZK69", "TZZK70", "TZZK72"}
	for _, drillName := range drillNames {
		realDrill, _ := constant.GetRealDrillByName(drillName)
		var xs, ys, pes []float64
		for y := float64(-500); y <= 500; y += 50 {
			for x := float64(-500); x <= 500; x += 50 {
				virtualDrill := GenVDrillM1(drillSet, blocks, pBlockLayerMat, realDrill.X+x, realDrill.Y+y)
				pe := GetPeByDrill(realDrill, virtualDrill)
				xs = append(xs, x)
				ys = append(ys, y)
				pes = append(pes, pe)
			}
		}
		fmt.Printf("%% %s 钻孔周围虚拟钻孔平均百分比误差云图\n", drillName)
		fmt.Print("x=")
		fmt.Print(xs)
		fmt.Println(";")

		fmt.Print("y=")
		fmt.Print(ys)
		fmt.Println(";")

		fmt.Print("z=")
		fmt.Print(pes)
		fmt.Println(";")
		fmt.Printf("[X,Y,Z]=griddata(x,y,z,linspace(-500,500,50)',linspace(-500,500,50),'v4');\n\n\n")

		log.Printf("max pes:%f, min pes: %f, avg pes: %f\n", floats.Max(pes), floats.Min(pes), utils.Average(pes))
	}

}
