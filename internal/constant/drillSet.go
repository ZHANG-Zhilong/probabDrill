package constant

import (
	"github.com/fogleman/poissondisc"
	"log"
	"math"
	"probabDrill"
	"probabDrill/internal/entity"
	"strconv"
	"strings"
	"sync"
)

var drillSetOnce sync.Once
var drillSet []entity.Drill
var drillsCeil, drillsFloor float64
var drillMap = make(map[string]int)

func GetRealDrillCF() (ceil, floor float64) {
	drillSetOnce.Do(initDrillsSet)
	ceil = drillsCeil
	floor = drillsFloor
	return
}
func GetRealDrillByName(name string) (drill entity.Drill, ok bool) {
	drillSetOnce.Do(initDrillsSet)
	if idx, ok := drillMap[name]; ok {
		drill = drillSet[idx]
		return drill, ok
	}
	return drill, ok
}
func GetRealDrills() []entity.Drill {
	drillSetOnce.Do(initDrillsSet)
	var drills []entity.Drill = make([]entity.Drill, len(drillSet))
	copy(drills, drillSet)
	return drills
}
func initDrillsSet() {
	log.SetFlags(log.Lshortfile)
	var drills []entity.Drill

	//add basic
	contents := readFile(probabDrill.Basic)
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
			d := entity.Drill{
				Name: temp[0],
				X:    decimal((x + probabDrill.OffX) * probabDrill.ScaleXY),
				Y:    decimal((y + probabDrill.OffY) * probabDrill.ScaleXY),
				Z:    decimal(z * probabDrill.ScaleZ),
			}

			//add ground layers, initial value.
			d.Layers = append(d.Layers, 0)
			d.LayerHeights = append(d.LayerHeights, d.Z)

			drills = append(drills, d)
			drillMap[d.Name] = len(drills) - 1
		}
	}

	//add layers
	contents = readFile(probabDrill.Layer)
	if strings.Index(contents, "\r\n") > 0 {
		log.Fatal("error, the file is crlf, not lf")
	}
	cs = strings.Split(contents, "\n")
	for _, d := range cs {
		temp := strings.Split(d, ",")
		if len(temp) == 0 {
			continue
		}
		var seq = GetSeqByName(temp[1])
		if idx, ok := drillMap[temp[0]]; ok {
			drills[idx].Layers = append(drills[idx].Layers, seq)
			depth, _ := strconv.ParseFloat(temp[2], 64)
			drills[idx].LayerHeights = append(drills[idx].LayerHeights, decimal(drills[idx].Z-depth))
		}
	}

	for _, d := range drills {
		if len(d.Layers) > 1 {
			drillSet = append(drillSet, d)
			drillsCeil = math.Max(d.Z, drillsCeil)
			drillsFloor = math.Min(drillsFloor, d.LayerHeights[len(d.LayerHeights)-1])
		}
	}
}

var helpDrillSetOnce sync.Once
var helpDrillsSet []entity.Drill
var helpDrillCeil, helpDrillFloor float64

func GetHelpDrills() []entity.Drill {
	helpDrillSetOnce.Do(initHelpDrillSet)
	return helpDrillsSet
}
func GetHelpDrillsCF() (ceil, floor float64) {
	helpDrillSetOnce.Do(initHelpDrillSet)
	return helpDrillCeil, helpDrillFloor
}
func initHelpDrillSet() {
	realDrills := GetRealDrills()
	x0, y0, x1, y1 := realDrills[0].GetRec(realDrills)
	points := poissondisc.Sample(x0, y0, x1, y1, probabDrill.MinDistance, probabDrill.MaxAttemptAdd, nil)
	for _, p := range points {
		idwDrill := genIDWDrill(realDrills, p.X, p.Y)
		helpDrillsSet = append(helpDrillsSet, idwDrill)
		helpDrillCeil = math.Max(helpDrillCeil, idwDrill.Z)
		helpDrillFloor = math.Min(helpDrillFloor, idwDrill.LayerHeights[len(idwDrill.LayerHeights)-1])
	}
	log.SetFlags(log.Lshortfile)
	log.Printf("len(helpDrills)=%d\n", len(helpDrillsSet))
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
