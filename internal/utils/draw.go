package utils

import (
	svg "github.com/ajstarks/svgo"
	"github.com/spf13/viper"
	"log"
	"math"
	"os"
	"probabDrill/apps/probDrill/model"
	"probabDrill/internal/constant"
	"strconv"
)

func getMappedY(dy, scaley float64) (y int) {
	if dy == 0 {
		y = viper.GetInt("CanvasHeight") / 10
		return y + viper.GetInt("CanvasOffSetY")
	}
	if dy > 0 {
		y = viper.GetInt("CanvasHeight")/10 - int(math.Round(scaley*dy))
		return y + viper.GetInt("CanvasOffSetY")
	}
	if dy < 0 {
		y = int(math.Round(math.Abs(scaley*dy))) + viper.GetInt("CanvasHeight")/10
		return y + viper.GetInt("CanvasOffSetY")
	}
	return y + viper.GetInt("CanvasOffSetY")
}
func drawDrill(canvas *svg.SVG, drill *model.Drill, x int, scaley float64) {
	if !drill.IsValid() {
		return
	}
	y0 := getMappedY(drill.LayerHeights[0], scaley)
	yb := getMappedY(drill.LayerHeights[len(drill.LayerHeights)-1], scaley)
	canvas.Line(x, y0, x, yb, "stroke-width=\"1\" stroke=\"black\"")
	//canvas.Rect(x-probabDrill.DrillWidth/2, y0, probabDrill.DrillWidth, yb-y0,
	//	"fill=\"none\" stroke=\"black\" stroke-width=\"1\"")
	canvas.Text(x, yb+25, drill.Name,
		"text-anchor:middle;font-size:4px;fill:black")
	for idx := 1; idx < len(drill.Layers); idx++ {
		y := getMappedY(drill.LayerHeights[idx], scaley)
		lasty := getMappedY(drill.LayerHeights[idx-1], scaley)
		//if y-lasty < 5 {
		//	y += model.CanvasMinThickness
		//}
		//canvas.Line(x-probabDrill.DrillWidth/2, y, x+probabDrill.DrillWidth/2, y,
		//	"stroke-width=\"1\" stroke=\"red\"")

		name := constant.GetNameBySeq(drill.Layers[idx])
		if lasty > 0 && y-lasty < viper.GetInt("CanvasMinThickness") {
			//canvas.Text(x+2, y-1, name, "text-anchor:start;font-size:5px;fill:black")
		} else {
			canvas.Text(x-1, y-1, name, "text-anchor:end;font-size:5px;fill:black")
		}
	}
}

func connect(canvas *svg.SVG, scaley float64, drill1 *model.Drill, x1 int, drill2 *model.Drill, x2 int) () {
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
		var lineId = strconv.Itoa(yl) + strconv.Itoa(yr)
		if _, ok := mapPos[lineId]; ok {
			continue
		} else {
			mapPos[lineId] = 1
			canvas.Line(x1, yl, x2, yr, "stroke-width=\"1\" stroke=\"blue\"")
			//canvas.Bezier(x1, yl, x1+5, yl, x2-5, yr, x2, yr, "stroke-width=\"1\" stroke=\"blue\" fill=\"none\"")
		}
	}
}

func DrawDrills(drills []model.Drill, picPath string) {
	log.SetFlags(log.Lshortfile)
	if _, err := os.Stat(picPath); !os.IsNotExist(err) {
		_ = os.Remove(picPath)
	}
	path, err := os.Create(picPath)
	if err != nil {
		panic(err)
	}
	defer path.Close()
	width := viper.GetInt("CanvasWidth")
	height := viper.GetInt("CanvasHeight")
	canvas := svg.New(path)
	canvas.Start(width, height)


	mappedDrillX := []float64{0}
	var maxDrillLength float64
	for idx := 1; idx < len(drills); idx++ {
		dist := (drills)[idx].Distance((drills)[idx-1])
		mappedDrillX = append(mappedDrillX, mappedDrillX[len(mappedDrillX)-1]+dist)
		maxDrillLength = math.Max(maxDrillLength, drills[idx].GetLength())
	}
	scalex := (float64(width) * 0.8) / mappedDrillX[len(mappedDrillX)-1]
	scaley := (float64(height) * 0.7) / maxDrillLength
	for i, d := range drills {
		drawDrill(canvas, &d, width/10+int(mappedDrillX[i]*scalex), scaley)
	}
	drills = constant.UnifyDrillsSeq(drills, constant.CheckSeqZiChun)
	for idx := 1; idx < len(drills); idx++ {
		x1 := width/10 + int(scalex*mappedDrillX[idx-1])
		x2 := width/10 + int(scalex*mappedDrillX[idx])
		connect(canvas, scaley, &drills[idx-1], x1, &drills[idx], x2)
	}
	canvas.End()
}
