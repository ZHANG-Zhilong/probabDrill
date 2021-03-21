package utils

import (
	"fmt"
	"log"
	"math"
	"strconv"
)

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
func GetLine(x1, y1, x2, y2, x float64) (y float64) {
	log.SetFlags(log.Lshortfile)
	if x2 == x1 {
		log.Fatal("error")
	}
	y = (x-x1)*(y2-y1)/(x2-x1) + y1
	return
}
func MiddleKPoints(x1, y1, x2, y2 float64, n int) (vertices []float64) {
	step := (x2 - x1) / float64(n)
	for x := x1 + step; math.Abs(x-x1) < math.Abs(x2-x1); x += step {
		y := GetLine(x1, y1, x2, y2, x)
		vertices = append(vertices, x, y)
	}
	return
}
func Decimal(value float64) float64 {
	value = math.Trunc(value*1e2+0.5) * 1e-2
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}

func PercentageError(observe, estimate []float64) (pe float64, err error) {
	var rmse float64
	if len(observe) != len(estimate) {
		return -1, fmt.Errorf(":input param error")
	}
	for idx, _ := range observe {
		rmse += math.Pow(observe[idx]-estimate[idx], 2)
	}
	rmse = math.Sqrt(rmse / float64(len(observe)))

	var estimateSum float64
	for _, v := range estimate {
		estimateSum += v
	}
	pe = rmse / (estimateSum / float64(len(estimate)))
	return pe, nil
}
func Average(arr []float64) (avg float64) {
	var sum float64
	for _, v := range arr {
		sum += v
	}
	return sum / float64(len(arr))
}
