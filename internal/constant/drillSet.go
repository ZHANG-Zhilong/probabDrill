package constant

import (
	"io/ioutil"
	"log"
	"os"
	"probabDrill/internal/entity"
	"strconv"
	"strings"
	"sync"
)

var once2 sync.Once
var drillSet []entity.Drill
var drillMap map[string]int = make(map[string]int)

var (
	mu      sync.Mutex
	drillId int
)

func GenVirtualDrillName() (name string) {
	mu.Lock()
	name = "virtual" + strconv.FormatInt(int64(drillId), 10)
	drillId++
	mu.Unlock()
	return name
}

//GetSeqByName return
func GetDrillByName(name string) (drill entity.Drill, ok bool) {
	init2()
	if idx, ok := drillMap[name]; ok {
		drill = drillSet[idx]
		return drill, ok
	}
	return drill, ok
}
func DrillSet() (ds []entity.Drill) {
	init2()
	ds = drillSet
	return
}

//init1 init
func init2() {
	once2.Do(func() {
		log.SetFlags(log.Lshortfile)
		var drills []entity.Drill

		//add basic
		contents := readFile(Basic)
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
					X:    (x + OffX) * ScaleXY,
					Y:    (y + OffY) * ScaleXY,
					Z:    z * ScaleZ,
				}

				//add ground layers, initial value.
				d.Layers = append(d.Layers, 0)
				d.LayerHeights = append(d.LayerHeights, d.Z)

				drills = append(drills, d)
				drillMap[d.Name] = len(drills) - 1
			}
		}

		//add layers
		contents = readFile(Layer)
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
				drills[idx].LayerHeights = append(drills[idx].LayerHeights, drills[idx].Z-depth)
			}
		}

		for _, d := range drills {
			if len(d.Layers) > 1 {
				drillSet = append(drillSet, d)
			}
		}
	})
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

func GetBoundary() (x, y []float64) {
	initBoundary()
	return bx, by
}

var onceInitBoundary sync.Once
var bx, by []float64

func initBoundary() {
	log.SetFlags(log.Lshortfile)
	onceInitBoundary.Do(func() {
		contents := readFile(Boundary)
		if strings.Index(contents, "\r\n") > 0 {
			log.Fatal("error, the file is crlf, not lf")
		}
		cs := strings.Split(contents, "\n")
		for _, p := range cs {
			temp := strings.Split(p, "  ")
			x, _ := strconv.ParseFloat(temp[0], 64)
			y, _ := strconv.ParseFloat(temp[1], 64)
			x = (x + OffX) * ScaleXY
			y = (y + OffY) * ScaleXY
			bx = append(bx, x)
			by = append(by, y)
		}
	})
}
func SimpleDrillSet() (drills []entity.Drill) {
	var drill1, drill2, drill3, drill4 entity.Drill
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

	drills = []entity.Drill{drill1, drill2, drill3, drill4}
	return
}
