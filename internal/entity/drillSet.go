package entity

import (
	"io/ioutil"
	"log"
	"math"
	"os"
	"probabDrill"
	"strconv"
	"strings"
	"sync"
)

var drillSetOnce sync.Once
var drillSet []Drill
var drillsCeil, drillsFloor float64
var drillMap map[string]int = make(map[string]int)

func GetDrillsCF() (ceil, floor float64) {
	drillSetOnce.Do(initDrillsSet)
	ceil = drillsCeil
	floor = drillsFloor
	return
}
func GetDrillByName(name string) (drill Drill, ok bool) {
	drillSetOnce.Do(initDrillsSet)
	if idx, ok := drillMap[name]; ok {
		drill = drillSet[idx]
		return drill, ok
	}
	return drill, ok
}
func DrillSet() []Drill {
	drillSetOnce.Do(initDrillsSet)
	var drills []Drill = make([]Drill, len(drillSet))
	copy(drills, drillSet)
	return drills
}
func initDrillsSet() {
	log.SetFlags(log.Lshortfile)
	var drills []Drill

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
			d := Drill{
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

//common utils
func
readFile(path string) string {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	return string(content)
}

func SimpleDrillSet() (drills []Drill) {
	var drill1, drill2, drill3, drill4 Drill
	drill1 = drill1.MakeDrill("1", 0, 0, 0)
	drill2 = drill1.MakeDrill("2", 1, 0, 0)
	drill3 = drill1.MakeDrill("3", 1, 1, 0)
	drill4 = drill1.MakeDrill("4", 0, 1, 0)

	drill1.AddLayer(1, -1)
	drill1.AddLayer(1, -2)
	drill1.AddLayer(6, -3)
	drill1.AddLayer(3, -4)

	drill2.AddLayer(2, -1)
	drill2.AddLayer(5, -2)
	drill2.AddLayer(3, -3)
	drill2.AddLayer(4, -4)

	drill3.AddLayer(1, -1)
	drill3.AddLayer(5, -2)
	drill3.AddLayer(6, -3)
	drill3.AddLayer(4, -4)

	drill4.AddLayer(1, -1)
	drill4.AddLayer(2, -2)
	drill4.AddLayer(3, -3)
	drill4.AddLayer(4, -4)

	drills = []Drill{drill1, drill2, drill3, drill4}
	return
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
