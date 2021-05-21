package constant

import (
	"fmt"
	"github.com/fogleman/poissondisc"
	"github.com/spf13/viper"
	"log"
	"math"
	"probabDrill/apps/probDrill/model"
	"strconv"
	"strings"
	"sync"
)

var drillSetOnce sync.Once
var drillSet []model.Drill
var drillsCeil, drillsFloor float64
var drillMap = make(map[string]int)

func GetRealDrillCF() (ceil, floor float64) {
	drillSetOnce.Do(initDrillsSet)
	ceil = drillsCeil
	floor = drillsFloor
	return
}
func GetRealDrillByName(name string) (drill model.Drill, ok bool) {
	drillSetOnce.Do(initDrillsSet)
	if idx, ok := drillMap[name]; ok {
		drill = drillSet[idx]
		return drill, ok
	}
	return drill, ok
}
func GetRealDrills() []model.Drill {
	drillSetOnce.Do(initDrillsSet)
	var drills []model.Drill = make([]model.Drill, len(drillSet))
	copy(drills, drillSet)
	return drills
}
func initDrillsSet() {
	log.SetFlags(log.Lshortfile)
	var drills []model.Drill

	//录入钻孔基本信息
	//contents := readFile(probabDrill.Basic)
	contents := readFile(viper.GetString("Basic"))
	cs := strings.Split(contents, "\n")
	if len(cs) < 10 {
		log.Fatal("error split")
	}
	for _, d := range cs {
		temp := strings.Split(d, ",")
		var valid = true
		for i := 0; i < len(temp); i++ {
			if len(temp[i]) == 0 {
				valid = false
			}
		}
		if valid {
			//make drill
			x, _ := strconv.ParseFloat(temp[1], 64)
			y, _ := strconv.ParseFloat(temp[2], 64)
			z, _ := strconv.ParseFloat(temp[3], 64)
			drill := model.Drill{
				Name: temp[0],
				X:    decimal((x + viper.GetFloat64("OffX")) * viper.GetFloat64("ScaleXY")),
				Y:    decimal((y + viper.GetFloat64("OffY")) * viper.GetFloat64("ScaleXY")),
				Z:    decimal((z + viper.GetFloat64("OffZ")) * viper.GetFloat64("ScaleZ")),
			}

			//add ground layers, initial value.
			drill.Layers = append(drill.Layers, 0)
			drill.LayerHeights = append(drill.LayerHeights, drill.Z)

			drills = append(drills, drill)
			drillMap[drill.Name] = len(drills) - 1
		}
	}

	//录入钻孔层位信息
	contents = readFile(viper.GetString("Layer"))
	if strings.Index(contents, "\r\n") > 0 {
		log.Fatal("error, the file is crlf, not lf")
	}
	cs = strings.Split(contents, "\n")
	for _, d := range cs {
		temp := strings.Split(d, ",")
		if len(temp) == 0 {
			continue
		}
		var layer = GetSeqByName(temp[1])
		if idx, ok := drillMap[temp[0]]; ok && drills != nil {
			depth, _ := strconv.ParseFloat(temp[2], 64)
			depth = decimal(depth)
			//剔除真实钻孔数据中的零厚度层
			drills[idx].Layers = append(drills[idx].Layers, layer)
			drills[idx].LayerHeights = append(drills[idx].LayerHeights, decimal(drills[idx].Z-depth))
		}
	}

	for _, d := range drills {
		if len(d.Layers) > 1 {
			d.UnBlock()
			drillSet = append(drillSet, d)
			drillsCeil = math.Max(d.Z, drillsCeil)
			drillsFloor = math.Min(drillsFloor, d.LayerHeights[len(d.LayerHeights)-1])
		}
	}

	//告警钻孔coordinate异常情况
	invalid := false
	for _, d := range drillSet {
		if math.Abs(d.X-0) > viper.GetFloat64("RealDrillDist") ||
			math.Abs(d.Y-0) > viper.GetFloat64("RealDrillDist") ||
			math.Abs(d.Z-0) > viper.GetFloat64("RealDrillDist") {
			fmt.Println(d)
			invalid = true
		}
	}
	if invalid {
		log.Fatal("invalid  drill data.")
	}
	for _, d := range drillSet {
		if !d.IsValid() {
			log.Fatal("invalid while initial drillSet.")
		}
	}
}

var helpDrillSetOnce sync.Once
var helpDrillsSet []model.Drill
var helpDrillCeil, helpDrillFloor float64

func GetHelpDrills() []model.Drill {
	helpDrillSetOnce.Do(initHelpDrillSet)
	return helpDrillsSet
}
func GetHelpDrillsCF() (ceil, floor float64) {
	helpDrillSetOnce.Do(initHelpDrillSet)
	return helpDrillCeil, helpDrillFloor
}
func initHelpDrillSet() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	realDrills := GetRealDrills()

	x0, y0, x1, y1 := realDrills[0].GetRec(realDrills)


	points := poissondisc.Sample(x0, y0, x1, y1, viper.GetFloat64("MinDistance"), viper.GetInt("MaxAttemptAdd"), nil)
	var helpDrills []model.Drill
	log.Println("initHelpDrillSet,泊松采样点数量为：", len(points))
	for _, p := range points {
		idwDrill := genIDWDrill(realDrills, p.X, p.Y)
		helpDrills = append(helpDrills, idwDrill)
		helpDrillCeil = math.Max(helpDrillCeil, idwDrill.Z)
		helpDrillFloor = math.Min(helpDrillFloor, idwDrill.LayerHeights[len(idwDrill.LayerHeights)-1])
	}
	helpDrillsSet = helpDrills
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.Printf("initHelpDrillSet,生成的有效的虚拟钻孔辅助数据为len(helpDrills)为：%d\n", len(helpDrillsSet))
}

var mu sync.Mutex
var drillId int

func GenVDrillName() (name string) {
	mu.Lock()
	name = "virtual" + strconv.FormatInt(int64(drillId), 10)
	drillId++
	mu.Unlock()
	return name
}
