package utils

import (
	"fmt"
	"log"
	"math"
	"probabDrill/apps/probDrill/model"
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

//find max value's index and value in []float64.
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

func GetMean(float64s []float64) (avg float64, err error) {
	if len(float64s) < 1 || float64s == nil {
		return 0, nil
	}
	sum, _ := GetSum(float64s)
	avg = sum / float64(len(float64s))
	if math.IsNaN(avg) || math.IsInf(avg, 10) {
		return 0, nil
	}
	return avg, nil
}
func GetSum(float64s []float64) (sum float64, err error) {
	if len(float64s) < 1 || float64s == nil {
		return 0, nil
	}
	for _, v := range float64s {
		sum += v
	}
	if math.IsNaN(sum) || math.IsInf(sum, 10) {
		return 0, nil
	}
	return sum, nil
}

//find second max value's index and value in []float64.
func FindSecondMaxFloat64s(float64s []float64) (idx int, val float64) {
	if len(float64s) < 1 {
		return 0, 0
	}
	_, maxVal := FindMaxFloat64s(float64s)
	idx, val = -math.MaxInt64, -math.MaxFloat64
	for id, v := range float64s {
		if math.Abs(v-maxVal) < 10e-7 {
			continue
		}
		if v > val {
			idx, val = id, v
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
	step := (x2 - x1) / float64(n+1)
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
	if len(observe) != len(estimate) {
		return -1, fmt.Errorf(":PercentageError,input param error,%#v,%#v", observe, estimate)
	}
	if len(observe) == len(estimate) && len(observe) == 0 {
		return -1, fmt.Errorf(":PercentageError,input param error,%#v,%#v", observe, estimate)
	}
	var rmse float64
	for idx, _ := range observe {
		rmse += math.Pow(observe[idx]-estimate[idx], 2)
	}
	rmse = math.Sqrt(rmse / float64(len(observe)))

	var estimateSum float64
	for _, v := range estimate {
		estimateSum += v
	}
	pe = rmse / (estimateSum / float64(len(estimate)))
	return Decimal(math.Abs(pe)), nil
}
func Average(arr []float64) (avg float64) {
	var sum float64
	for _, v := range arr {
		sum += v
	}
	return sum / float64(len(arr))
}
func Drill2WXD(drills []model.Drill) (rst string) {
	fmt.Print("hole_mtx_1 = [")
	for didx, d := range drills {
		fmt.Print("[[")
		for idx := 1; idx < len(d.Layers); idx++ {
			fmt.Print(d.Layers[idx])
			if idx != len(d.Layers)-1 {
				fmt.Print(",")
			}
		}
		fmt.Print("],[")
		fmt.Print(Decimal(math.Abs(d.LayerHeights[0] - d.Z)))
		for idx := 1; idx < len(d.LayerHeights); idx++ {
			fmt.Print(",", Decimal(d.Z-d.LayerHeights[idx]))
		}
		fmt.Print("]]\n")
		if didx != len(drills)-1 {
			fmt.Print(",")
		}
	}
	fmt.Println("]")

	fmt.Print("delt_h = [")
	for idx, d := range drills {
		fmt.Print(Decimal(d.Z))
		if idx != len(drills)-1 {
			fmt.Print(",")
		}
	}
	fmt.Print("]\n")

	fmt.Print("holes = [")
	for _, d := range drills {
		fmt.Print("[")
		fmt.Print("\"", d.Name, "\"", ",")
		fmt.Print("\"", d.Name, "\"", ",")
		fmt.Print(d.X, ",", d.Y, ",", d.Z, ",", Decimal(d.GetLength()))
		fmt.Print("],")
	}
	fmt.Print("]\n\n")
	return
}
func TruncDrills(vdrills []model.Drill) {
	var avgBot float64
	for _, d := range vdrills {
		avgBot += d.BottomHeight()
	}
	avgBot = avgBot / float64(len(vdrills))
	for idx, _ := range vdrills {
		vdrills[idx] = vdrills[idx].Trunc(avgBot)
	}
}
