package utils

import (
	"fmt"
	"github.com/ajstarks/svgo"
	"log"
	"math"
	"os"
	"probabDrill/internal/entity"
	"strconv"
)

func DisplayDrills(drills []entity.Drill) {
	for _, d := range drills {
		d.Print()
	}
	fmt.Printf("total %d drills.", len(drills))
}
func PrintFloat64s(s []float64) () {
	fmt.Print("[")
	for _, v := range s {
		if v > 0 {
			fmt.Printf("%.2f\t", v)
		} else {
			fmt.Print("\t\t")
		}
	}
	fmt.Print("]\n")
}
func Hole(vals ...float64) () {
	return
}
func IsInPolygon(x, y []float64, x0, y0 float64) (isIn bool) {

	//vert[0], vert[last]
	var i, j int = 0, len(x) - 1
	if (y[i] >= y0) != (y[j] > y0) &&
		(y0 <= y[i] && y0 <= y[j] ||
			x0 <= (y0-y[i])*(x[j]-x[i])/(y[j]-y[i])+x[i]) {
		isIn = !isIn
	}

	//y0 is among y1 and y2, ray x0
	//if k=inf -> y1==y2  y0<=y1&&y0<y2 cross
	//if k< inf	x0<x1+k(y0-y1) cross
	for i := 1; i < len(x); i++ {
		if (y[i] >= y0) != (y[j] > y0) &&
			(y0 <= y[i] && y0 <= y[j] ||
				x0 <= (y0-y[i])*(x[j]-x[i])/(y[j]-y[i])+x[i]) {
			isIn = !isIn
		}
	}

	return isIn
}
func GetGrids(px, py, l, r, t, b float64) (gridx, gridy []float64) {
	gridx = append(gridx, l)
	gridy = append(gridy, b)
	for (l + px) < r {
		gridx = append(gridx, l+px)
		l = l + px
	}
	for (b + py) < t {
		gridy = append(gridy, b+py)
		b = b + py
	}
	gridx = append(gridx, r)
	gridy = append(gridy, t)
	return
}
func FindMaxFloat64s(float64s []float64) (idx int, val float64) {
	if len(float64s) < 1 {
		return 0, 0
	}
	idx, val = -math.MaxInt64, -math.MaxFloat64
	for id, va := range float64s {
		if va > val {
			idx, val = id, va
		}
	}
	return idx, val
}

func drawDrill(
	canvas *svg.SVG,
	drill entity.Drill,
	x, y int,
	scaley float64,
) {
	if !drill.IsValid() {
		return
	}
	if drill.Z > 0 {
		y -= int(drill.Z * scaley)
	} else {
		y += int(drill.Z * scaley)
	}
	var lastx, lasty int = x - 5, y
	for idx := 1; idx < len(drill.Layers); idx++ {
		blockLength := int((drill.LayerFloorHeights[idx-1] - drill.LayerFloorHeights[idx]) * scaley)
		if blockLength < 10 {
			blockLength *= 2
		}
		canvas.Rect(lastx, lasty, 10, blockLength, "fill=\"none\" stroke=\"black\" stroke-width=\"1\"")
		canvas.Text(lastx+12, lasty+blockLength*2/3,
			strconv.FormatInt(drill.Layers[idx], 10), "text-anchor:start;font-size:7px;fill:black")
		lasty += blockLength
	}
	canvas.Text(lastx+5, lasty+15, drill.Name, "text-anchor:middle;font-size:7px;fill:black")
}
func DrawDrills(drills []entity.Drill) {
	log.SetFlags(log.Lshortfile)
	if len(drills) < 2 {
		log.Fatal("钻孔过少",len(drills))
		return
	}
	width := 500
	height := 600
	if _, err := os.Stat("./out.svg"); !os.IsNotExist(err) {
		os.Remove("./out.svg")
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
		drawDrill(canvas, d, x0+int(float64(distAccum[i])*scalex), y0, scaley)
	}
	canvas.End()
}

func SvgDemo() {

	width := 500
	height := 500
	if _, err := os.Stat("./out.svg"); !os.IsNotExist(err) {
		os.Remove("./out.svg")
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
