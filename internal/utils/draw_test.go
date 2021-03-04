package utils

import (
	"fmt"
	"probabDrill/internal/constant"
	"testing"
)

func TestSvgDemo(t *testing.T) {
	SvgDemo()
}
func TestDisplayDrills(t *testing.T) {
	drill := constant.DrillSet()[0]
	y0 := getMappedY(drill, 600, 3.49, 6.5)
	fmt.Println(y0)
}
