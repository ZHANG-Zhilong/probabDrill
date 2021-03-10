package play

import (
	"fmt"
	"github.com/fogleman/gg"
	"github.com/fogleman/poissondisc"
)

func Play2() {
	const (
		W = 2400
		H = 1600
		R = 8
		K = 32
	)
	points := poissondisc.Sample(0, 0, W, H, R, K, nil)
	fmt.Println(len(points), "points")

	dc := gg.NewContext(W, H)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	for _, p := range points {
		dc.DrawPoint(p.X, p.Y, R*0.45)
	}
	dc.SetRGB(0, 0, 0)
	dc.Fill()
	dc.SavePNG("poissiondist1.png")
}

func Play1() {
	var x0, y0, x1, y1, r float64
	x0 = 0    // bbox min
	y0 = 0    // bbox min
	x1 = 1000 // bbox max
	y1 = 1000 // bbox max
	r = 10    // min distance between points
	k := 30   // max attempts to add neighboring point

	points := poissondisc.Sample(x0, y0, x1, y1, r, k, nil)

	for _, p := range points {
		fmt.Println(p.X, p.Y)
	}
}
