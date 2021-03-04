package utils

import (
	"github.com/ajstarks/svgo"
	"log"
	"math"
	"os"
	"probabDrill/internal/entity"
	"strconv"
)

func getMappedY(drill entity.Drill, height int, dy, scaley float64) (y int) {
	if dy == 0 {
		y = int(height / 10)
		return
	}
	if dy > 0 {
		y = height/10 - int(scaley*float64(dy))
		return
	}
	if dy < 0 {
		y = int(scaley*math.Abs(dy)) + height/10
		return
	}
	return
}

func drawDrill(canvas *svg.SVG, width, height int, drill entity.Drill, x, y int, scaley float64) {
	if !drill.IsValid() {
		return
	}
	drillWidth := width / 50
	for idx := 1; idx < len(drill.Layers); idx++ {
		x0 := x - drillWidth/2
		y0 := getMappedY(drill, height, drill.LayerFloorHeights[idx-1], scaley)
		blockLength := getMappedY(drill, height, drill.LayerFloorHeights[idx], scaley) - y0

		canvas.Rect(x0, y0, drillWidth, blockLength, "fill=\"none\" stroke=\"black\" stroke-width=\"1\"")
		canvas.Text(x, y0+blockLength*2/3,
			strconv.FormatInt(drill.Layers[idx], 10), "text-anchor:middle;font-size:7px;fill:black")
	}
	canvas.Text(x, getMappedY(drill, height, drill.GetBottomHeight(), scaley)+15, drill.Name,
		"text-anchor:middle;font-size:7px;fill:black")
}
func nextIdx(idx int, drill entity.Drill) (id int) {
	if idx < len(drill.Layers)-1 {
		id = idx + 1
	} else {
		id = len(drill.Layers) - 1
	}
	return
}
func connect(canvas *svg.SVG, width, height int, scaley float64, drill1 entity.Drill, x1 int, drill2 entity.Drill, x2 int) () {
	drillWidth := width / 50
	idx1, idx2 := 0, 0
	var flag = 1
	for {
		if idx1 == len(drill1.Layers)-1 && idx2 == len(drill2.Layers)-1 {
			flag--
		}
		x1, x2 := x1+drillWidth/2, x2-drillWidth/2
		if drill1.Layers[idx1] == drill2.Layers[idx2] {
			y1 := getMappedY(drill1, height, drill1.LayerFloorHeights[idx1], scaley)
			y2 := getMappedY(drill2, height, drill2.LayerFloorHeights[idx2], scaley)
			canvas.Line(x1, y1, x2, y2, "stroke-width=\"1\" stroke=\"blue\"")
			idx1 = nextIdx(idx1, drill1)
			idx2 = nextIdx(idx2, drill2)
		}
		if drill1.Layers[idx1] < drill2.Layers[idx2] {
			y1 := getMappedY(drill1, height, drill1.LayerFloorHeights[idx1], scaley)
			y2 := getMappedY(drill2, height, drill2.LayerFloorHeights[idx2-1], scaley)
			canvas.Line(x1, y1, x2, y2, "stroke-width=\"1\" stroke=\"blue\"")
			idx1 = nextIdx(idx1, drill1)
		}
		if drill1.Layers[idx1] > drill2.Layers[idx2] {
			y1 := getMappedY(drill1, height, drill1.LayerFloorHeights[idx1-1], scaley)
			y2 := getMappedY(drill2, height, drill2.LayerFloorHeights[idx2], scaley)
			canvas.Line(x1, y1, x2, y2, "stroke-width=\"1\" stroke=\"blue\"")
			idx2 = nextIdx(idx2, drill2)
		}
		if idx1 == len(drill1.Layers)-1 || idx2 == len(drill2.Layers)-1 {
			idx1 = nextIdx(idx1, drill1)
			idx2 = nextIdx(idx2, drill2)
		}
		if flag == 0 {
			y1 := getMappedY(drill1, height, drill1.LayerFloorHeights[idx1], scaley)
			y2 := getMappedY(drill2, height, drill2.LayerFloorHeights[idx2], scaley)
			canvas.Line(x1, y1, x2, y2, "stroke-width=\"1\" stroke=\"blue\"")
			break
		}
	}

	//for _, layer := range drill1.Layers {
	//	x1, x2 := x1+drillWidth/2, x2-drillWidth/2
	//	b1s := drill1.GetBottomHeightByLayer(layer)
	//	b2s := drill2.GetBottomHeightByLayer(layer)
	//	if len(b1s) < 1 || len(b2s) < 1 {
	//		return
	//	}
	//	y1 := getMappedY(drill1, height, b1s[0], scaley)
	//	y2 := getMappedY(drill2, height, b2s[0], scaley)
	//	canvas.Line(x1, y1, x2, y2, "stroke-width=\"1\" stroke=\"blue\"")
	//}
}
func DrawDrills(drills []entity.Drill, dist int) {
	log.SetFlags(log.Lshortfile)
	if len(drills) < 2 {
		log.Fatal("钻孔过少", len(drills))
		return
	}
	width := len(drills) * dist
	height := 600
	if _, err := os.Stat("./out.svg"); !os.IsNotExist(err) {
		_ = os.Remove("./out.svg")
	}
	path, err := os.Create("./out.svg")
	if err != nil {
		panic(err)
	}
	defer path.Close()

	canvas := svg.New(path)
	canvas.Start(width, height)
	distAccum := []int{0}
	var lengthMax int
	for idx := 1; idx < len(drills); idx++ {
		dist := drills[idx].DistanceBetween(drills[idx-1])
		distAccum = append(distAccum, int(math.Ceil(dist))+distAccum[len(distAccum)-1])
		if float64(lengthMax) < drills[idx].GetLength() {
			lengthMax = int(math.Ceil(drills[idx].GetLength()))
		}
	}
	scalex := (float64(width) * 0.8) / float64(distAccum[len(distAccum)-1])
	scaley := (float64(height) * 0.8) / float64(lengthMax)
	x0, y0 := width/10, height/10
	for i, d := range drills {
		drawDrill(canvas, width, height, d, x0+int(float64(distAccum[i])*scalex), y0, scaley)
	}
	for idx := 1; idx < len(drills); idx++ {
		x1 := x0 + int(scalex*float64(distAccum[idx-1]))
		x2 := x0 + int(scalex*float64(distAccum[idx]))
		connect(canvas, width, height, scaley, drills[idx-1], x1, drills[idx], x2)
	}
	canvas.End()
}

func SvgDemo() {

	width := 500
	height := 500
	if _, err := os.Stat("./out.svg"); !os.IsNotExist(err) {
		_ = os.Remove("./out.svg")
	}
	path, err := os.Create("./out.svg")
	if err != nil {
		panic(err)
	}
	defer path.Close()
	canvas := svg.New(path)
	canvas.Start(width, height)
	canvas.Circle(width/2, height/2, 100)
	canvas.Text(width/2, height/2, "Hello, SVG", "text-anchor:middle;font-size:30px;fill:white")
	canvas.CenterRect(100, 100, 50, 50, "fill=\"none\" stroke=\"blue\" stroke-width=\"20\"")
	canvas.Line(0, 0, width, height, "stroke-width=\"5\" stroke=\"blue\"")
	canvas.Bezier(0, 250, 75, 250, 150, 350, 200, 350, "fill=\"none\" stroke-width=\"5\" stroke=\"blue\"")
	canvas.End()
}
func GetLine(x1, y1, x2, y2, x float64) (y float64) {
	log.SetFlags(log.Lshortfile)
	if x2 == x1 {
		log.Fatal("error")
	}
	y = (x-x1)*(y2-y1)/(x2-x1) + y1
	return
}
