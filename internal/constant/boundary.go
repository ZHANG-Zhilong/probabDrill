package constant

import (
	"github.com/spf13/viper"
	"log"
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
	contents := readFile(viper.GetString("Boundary"))
	if strings.Index(contents, "\r\n") > 0 {
		log.Fatal("error, the file is crlf, not lf")
	}
	cs := strings.Split(contents, "\n")
	for _, p := range cs {
		temp := strings.Split(p, "  ")
		x, _ := strconv.ParseFloat(temp[0], 64)
		y, _ := strconv.ParseFloat(temp[1], 64)

		x = (x + viper.GetFloat64("OffX")) * viper.GetFloat64("ScaleXY")
		y = (y + viper.GetFloat64("OffY")) * viper.GetFloat64("ScaleXY")
		bx = append(bx, x)
		by = append(by, y)
	}
}
