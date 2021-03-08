package utils

import (
	"github.com/ajstarks/svgo"
	"log"
	"math"
	"os"
	"probabDrill/internal/constant"
	"probabDrill/internal/entity"
	"strconv"
)

func getMappedY(drill entity.Drill, dy, scaley float64) (y int) {
	if dy == 0 {
		y = int(constant.CanvasHeight / 10)
		return y + constant.CanvasOffsetY
	}
	if dy > 0 {
		y = constant.CanvasHeight/10 - int(scaley*float64(dy))
		return y + constant.CanvasOffsetY
	}
	if dy < 0 {
		y = int(scaley*math.Abs(dy)) + constant.CanvasHeight/10
		return y + constant.CanvasOffsetY
	}
	return y + constant.CanvasOffsetY
}
func drawDrill(canvas *svg.SVG, drill entity.Drill, x int, scaley float64) {
	if !drill.IsValid() {
		return
	}
	y0 := getMappedY(drill, drill.LayerHeights[0], scaley)
	canvas.Rect(x-constant.DrillWidth/2, y0, constant.DrillWidth, int(drill.GetLength()*scaley),
		"fill=\"none\" stroke=\"black\" stroke-width=\"1\"")
	canvas.Text(x, getMappedY(drill, drill.GetBottomHeight(), scaley)+15, drill.Name,
		"text-anchor:middle;font-size:7px;fill:black")

	for idx := 1; idx < len(drill.Layers); idx++ {
		var lasty, y int
		y = getMappedY(drill, drill.LayerHeights[idx-1], scaley)
		if idx > 1 {
			lasty = getMappedY(drill, drill.LayerHeights[idx-2], scaley)
			if y-lasty < 5 {
				y += constant.CanvasMinThickness
			}
		}

		canvas.Line(x-constant.DrillWidth/2, y, x+constant.DrillWidth/2, y,
			"stroke-width=\"1\" stroke=\"red\"")

		if lasty > 0 && y-lasty < 5 {
			canvas.Text(x, y+18, strconv.FormatInt(int64(drill.Layers[idx]), 10),
				"text-anchor:start;font-size:1px;fill:black")
		} else {
			canvas.Text(x, y+10, strconv.FormatInt(int64(drill.Layers[idx]), 10),
				"text-anchor:middle;font-size:1px;fill:black")
		}
	}

}
func nextIdx(idx int, drill entity.Drill) (id int) {
	if idx < len(drill.Layers)-1 {
		id = idx + 1
	} else {
		id = len(drill.Layers) - 1
	}
	return
}
func connect(canvas *svg.SVG, scaley float64, drill1 entity.Drill, x1 int, drill2 entity.Drill, x2 int) () {
	log.SetFlags(log.Lshortfile)
	if len(drill1.Layers) != len(drill2.Layers) {
		log.Fatal("error")
	} else {
		for idx1 := 0; idx1 < len(drill1.Layers); idx1++ {
			if drill1.Layers[idx1] != drill2.Layers[idx1] {
				log.Fatal("error")
			}
		}
	}
	var mapPos map[string]int = make(map[string]int)
	for idx1 := 0; idx1 < len(drill1.Layers); idx1++ {
		yl := getMappedY(drill1, drill1.LayerHeights[idx1], scaley)
		yr := getMappedY(drill2, drill2.LayerHeights[idx1], scaley)
		var lineId string = strconv.Itoa(yl) + strconv.Itoa(yr)
		if _, ok := mapPos[lineId]; ok {
			continue
		} else {
			mapPos[lineId] = 1
			canvas.Line(x1+constant.DrillWidth/2, yl, x2-constant.DrillWidth/2, yr,
				"stroke-width=\"1\" stroke=\"blue\"")
		}
	}
	//var flag = 1
	//for {
	//	if idx1 == len(drill1.Layers)-1 && idx2 == len(drill2.Layers)-1 {
	//		flag--
	//	}
	//	x1, x2 := x1+constant.DrillWidth/2, x2-constant.DrillWidth/2
	//	if drill1.Layers[idx1] == drill2.Layers[idx2] {
	//		y1 := getMappedY(drill1, drill1.LayerHeights[idx1], scaley)
	//		y2 := getMappedY(drill2, drill2.LayerHeights[idx2], scaley)
	//		canvas.Line(x1, y1, x2, y2, "stroke-width=\"1\" stroke=\"blue\"")
	//		idx1 = nextIdx(idx1, drill1)
	//		idx2 = nextIdx(idx2, drill2)
	//	}
	//	if drill1.Layers[idx1] < drill2.Layers[idx2] {
	//		y1 := getMappedY(drill1, drill1.LayerHeights[idx1], scaley)
	//		y2 := getMappedY(drill2, drill2.LayerHeights[idx2-1], scaley)
	//		canvas.Line(x1, y1, x2, y2, "stroke-width=\"1\" stroke=\"blue\"")
	//		idx1 = nextIdx(idx1, drill1)
	//	}
	//	if drill1.Layers[idx1] > drill2.Layers[idx2] {
	//		y1 := getMappedY(drill1, drill1.LayerHeights[idx1-1], scaley)
	//		y2 := getMappedY(drill2, drill2.LayerHeights[idx2], scaley)
	//		canvas.Line(x1, y1, x2, y2, "stroke-width=\"1\" stroke=\"blue\"")
	//		idx2 = nextIdx(idx2, drill2)
	//	}
	//	if idx1 == len(drill1.Layers)-1 || idx2 == len(drill2.Layers)-1 {
	//		idx1 = nextIdx(idx1, drill1)
	//		idx2 = nextIdx(idx2, drill2)
	//	}
	//	if flag == 0 {
	//		y1 := getMappedY(drill1, drill1.LayerHeights[idx1], scaley)
	//		y2 := getMappedY(drill2, drill2.LayerHeights[idx2], scaley)
	//		canvas.Line(x1, y1, x2, y2, "stroke-width=\"1\" stroke=\"blue\"")
	//		break
	//	}
	//}
}

func DrawDrills(drills []entity.Drill) {
	log.SetFlags(log.Lshortfile)
	drills = drills[0].RoundDrills(drills)
	if len(drills) < 2 {
		log.Fatal("钻孔过少", len(drills))
		return
	}
	if _, err := os.Stat("./out.svg"); !os.IsNotExist(err) {
		_ = os.Remove("./out.svg")
	}
	path, err := os.Create("./out.svg")
	if err != nil {
		panic(err)
	}
	defer path.Close()
	width := constant.CanvasWidth
	height := constant.CanvasHeight
	canvas := svg.New(path)
	canvas.Start(width, height)
	accumDist := []int{0}
	var lengthMax int
	for idx := 1; idx < len(drills); idx++ {
		dist := drills[idx].DistanceBetween(drills[idx-1])
		accumDist = append(accumDist, int(math.Round(dist))+accumDist[len(accumDist)-1])
		if float64(lengthMax) < drills[idx].GetLength() {
			lengthMax = int(math.Round(drills[idx].GetLength()))
		}
	}
	scalex := (float64(width) * 0.8) / float64(accumDist[len(accumDist)-1])
	scaley := (float64(height) * 0.8) / float64(lengthMax)
	for i, d := range drills {
		drawDrill(canvas, d, width/10+int(float64(accumDist[i])*scalex), scaley)
	}

	drills = UnifyDrillsStrata(&drills, CheckSeqZiChun)
	for idx := 1; idx < len(drills); idx++ {
		x1 := width/10 + int(scalex*float64(accumDist[idx-1]))
		x2 := width/10 + int(scalex*float64(accumDist[idx]))
		connect(canvas, scaley, drills[idx-1], x1, drills[idx], x2)
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
func SplitSegment(x1, y1, x2, y2 float64, n int) (vertices []float64) {
	step := (x2 - x1) / float64(n+1)
	for x := x1 + step; math.Abs(x-x1) < math.Abs(x2-x1); x += step {
		y := GetLine(x1, y1, x2, y2, x)
		vertices = append(vertices, x, y)
	}
	return
}
