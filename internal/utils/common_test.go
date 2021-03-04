package utils

import (
	"probabDrill/internal/constant"
	"testing"
)

func TestSvgDemo(t *testing.T) {
	SvgDemo()
}
func TestDisplayDrills(t *testing.T) {
	drills := constant.DrillSet()[:3]
	DrawDrills(drills)
}
