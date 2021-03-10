package utils

import (
	"github.com/ajstarks/svgo"
	"log"
	"math"
	"os"
	"probabDrill"
	"probabDrill/internal/entity"
	"strconv"
)

func getMappedY(dy, scaley float64) (y int) {
	if dy == 0 {
		y = int(probabDrill.CanvasHeight / 10)
		return y + probabDrill.CanvasOffsetY
	}
	if dy > 0 {
		y = probabDrill.CanvasHeight/10 - int(math.Round(scaley*dy))
		return y + probabDrill.CanvasOffsetY
	}
	if dy < 0 {
		y = int(math.Round(math.Abs(scaley*dy))) + probabDrill.CanvasHeight/10
		return y + probabDrill.CanvasOffsetY
	}
	return y + probabDrill.CanvasOffsetY
}
func drawDrill(canvas *svg.SVG, drill *entity.Drill, x int, scaley float64) {
	if !drill.IsValid() {
		return
	}
	y0 := getMappedY(drill.LayerHeights[0], scaley)
	yb := getMappedY(drill.LayerHeights[len(drill.LayerHeights)-1], scaley)
	canvas.Rect(x-probabDrill.DrillWidth/2, y0, probabDrill.DrillWidth, yb-y0,
		"fill=\"none\" stroke=\"black\" stroke-width=\"1\"")
	canvas.Text(x, yb+15, drill.Name,
		"text-anchor:middle;font-size:7px;fill:black")
	for idx := 1; idx < len(drill.Layers); idx++ {
		y := getMappedY(drill.LayerHeights[idx], scaley)
		lasty := getMappedY(drill.LayerHeights[idx-1], scaley)
		//if y-lasty < 5 {
		//	y += constant.CanvasMinThickness
		//}
		canvas.Line(x-probabDrill.DrillWidth/2, y, x+probabDrill.DrillWidth/2, y,
			"stroke-width=\"1\" stroke=\"red\"")

		if lasty > 0 && y-lasty < probabDrill.CanvasMinThickness {
			canvas.Text(x+probabDrill.DrillWidth/2, y-1, strconv.FormatInt(int64(drill.Layers[idx]), 10),
				"text-anchor:start;font-size:2px;fill:red")
		} else {
			canvas.Text(x, y, strconv.FormatInt(int64(drill.Layers[idx]), 10),
				"text-anchor:middle;font-size:2px;fill:black")
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
func connect(canvas *svg.SVG, scaley float64, drill1 *entity.Drill, x1 int, drill2 *entity.Drill, x2 int) () {
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
		yl := getMappedY(drill1.LayerHeights[idx1], scaley)
		yr := getMappedY(drill2.LayerHeights[idx1], scaley)
		var lineId string = strconv.Itoa(yl) + strconv.Itoa(yr)
		if _, ok := mapPos[lineId]; ok {
			continue
		} else {
			mapPos[lineId] = 1
			//canvas.Line(x1+probabDrill.DrillWidth/2, yl, x2-probabDrill.DrillWidth/2, yr,
			//	"stroke-width=\"1\" stroke=\"blue\"")
			canvas.Bezier(x1, yl, x1+5, yl, x2-5, yr, x2, yr, "stroke-width=\"1\" stroke=\"blue\" fill=\"none\"")
		}
	}
}

func DrawDrills(drills []entity.Drill, picPath string) {
	log.SetFlags(log.Lshortfile)
	//drills = (*drills)[0].unusedRoundDrills(*drills)
	if len(drills) < 2 {
		log.Fatal("钻孔过少", len(drills))
		return
	}
	if _, err := os.Stat(picPath); !os.IsNotExist(err) {
		_ = os.Remove(picPath)
	}
	path, err := os.Create(picPath)
	if err != nil {
		panic(err)
	}
	defer path.Close()
	width := probabDrill.CanvasWidth
	height := probabDrill.CanvasHeight
	canvas := svg.New(path)
	canvas.Start(width, height)
	accumDist := []int{0}
	var lengthMax int
	for idx := 1; idx < len(drills); idx++ {
		dist := (drills)[idx].Distance((drills)[idx-1])
		accumDist = append(accumDist, int(math.Round(dist))+accumDist[len(accumDist)-1])
		if float64(lengthMax) < (drills)[idx].GetLength() {
			lengthMax = int(math.Round((drills)[idx].GetLength()))
		}
	}
	scalex := (float64(width) * 0.8) / float64(accumDist[len(accumDist)-1])
	scaley := (float64(height) * 0.8) / float64(lengthMax)
	for i, d := range drills {
		drawDrill(canvas, &d, width/10+int(float64(accumDist[i])*scalex), scaley)
	}

	drills = UnifyDrillsStrata(drills, CheckSeqZiChun)
	for idx := 1; idx < len(drills); idx++ {
		x1 := width/10 + int(scalex*float64(accumDist[idx-1]))
		x2 := width/10 + int(scalex*float64(accumDist[idx]))
		connect(canvas, scaley, &drills[idx-1], x1, &drills[idx], x2)
	}
	canvas.End()
}
