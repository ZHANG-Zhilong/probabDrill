package constant

import (
	"log"
	"probabDrill"
	"strconv"
	"strings"
	"sync"
)

var boundaryOnce sync.Once
var bx, by []float64

func GetBoundary() (x, y []float64) {
	boundaryOnce.Do(initBoundary)
	return bx, by
}
func initBoundary() {
	log.SetFlags(log.Lshortfile)
	contents := readFile(probabDrill.Boundary)
	if strings.Index(contents, "\r\n") > 0 {
		log.Fatal("error, the file is crlf, not lf")
	}
	cs := strings.Split(contents, "\n")
	for _, p := range cs {
		temp := strings.Split(p, "  ")
		x, _ := strconv.ParseFloat(temp[0], 64)
		y, _ := strconv.ParseFloat(temp[1], 64)
		x = (x + probabDrill.OffX) * probabDrill.ScaleXY
		y = (y + probabDrill.OffY) * probabDrill.ScaleXY
		bx = append(bx, x)
		by = append(by, y)
	}
}
