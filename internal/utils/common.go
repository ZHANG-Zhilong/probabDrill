package utils

import (
	"fmt"
	"math"
	"probabDrill/internal/entity"
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

func UnifyStratum(drills *[]entity.Drill) (uniLayers []int) {
	if len(*drills) < 2 {
		return
	}
	var newLayer []int = (*drills)[0].Layers
	CheckSeq(&newLayer)
	for idx := 1; idx < len(*drills); idx++ {
		layers := (*drills)[idx].Layers
		newLayer = UnifySeq(layers, newLayer)
	}
	return newLayer
}
func UnifySeq(seq1, seq2 []int) (seq []int) {
	idx1, idx2 := 1, 1
	seq = []int{0}
	CheckSeq(&seq1)
	CheckSeq(&seq2)
	for {
		if idx1 == len(seq1) && idx2 == len(seq2) {
			break
		}
		if idx2 == len(seq2) && idx1 < len(seq1) ||
			idx1 < len(seq1) && idx2 < len(seq2) && seq1[idx1] < seq2[idx2] {
			addLayerForUnify(&seq, seq1[idx1])
			idx1++
			continue
		}
		if idx1 == len(seq1) && idx2 < len(seq2) ||
			idx1 < len(seq1) && idx2 < len(seq2) && seq1[idx1] > seq2[idx2] {
			addLayerForUnify(&seq, seq2[idx2])
			idx2++
			continue
		}
		if seq1[idx1] == seq2[idx2] {
			addLayerForUnify(&seq, seq2[idx2])
			idx1++
			idx2++
			continue
		}
	}
	return seq
}
func isNormal(layers *[]int, idx int) (ok bool) {
	if idx <= len(*layers)-1 && idx-1 >= 0 {
		return prePos(layers, idx)+1 == (*layers)[idx]
	}
	if idx >= 1 && idx+1 <= len(*layers)-1 {
		return (*layers)[idx]+1 == (*layers)[idx+1]
	}
	return
}
func isLack(layers *[]int, idx int) (ok bool) {
	//has gap
	if idx <= len(*layers)-1 && idx-1 >= 1 {
		return (*layers)[idx] > prePos(layers, idx)+1
	}
	return
}
func isInverse(layers *[]int, idx int) (ok bool) {
	//case 1: drill top, ->   😄  smaller  go top
	if idx == 1 && idx+1 < len(*layers) {
		if (*layers)[idx] > (*layers)[idx+1] {
			return true
		}
	}

	//case 2: drill bottom -> bigger 😄  go bottom
	if idx == len(*layers)-1 && idx-1 >= 1 {
		if prePos(layers, idx) > (*layers)[idx] {
			return true
		}
	}

	//case 3: drill middle -> bigger 😄 bigger  go down
	//case 4: drill middle -> bigger 😄 smaller  go down
	//case 5: drill middle -> smaller 😄 smaller  go up
	if idx-1 >= 1 && idx+1 < len(*layers) {

		if prePos(layers, idx) > (*layers)[idx] {
			return true
		}

		if prePos(layers, idx)+1 < (*layers)[idx] && (*layers)[idx] > (*layers)[idx+1] {
			return true
		}
	}

	return false
}
func prePos(layers *[]int, idx int) (layer int) {
	for {
		if idx-1 == 0 {
			break
		}
		if (*layers)[idx-1] > 0 {
			return (*layers)[idx-1]
		} else {
			idx--
		}
		//hello
	}
	return 0
}
func CheckSeq(layers *[]int) {
	for _, layer := range *layers {
		if layer < 0 {
			return
		}
	}
	//mark inverse and repeat
	var layerMap map[int]int = make(map[int]int)
	for i := 1; i < len(*layers); i++ {
		if isNormal(layers, i) { //normal first
			layerMap[(*layers)[i]] = 1
			continue
		}
		if isLack(layers, i) { //lack first
			layerMap[(*layers)[i]] = 1
			continue
		}
		if isInverse(layers, i) {
			(*layers)[i] = -(*layers)[i]
			layerMap[(*layers)[i]] = 1
			continue
		}
		if _, ok := layerMap[(*layers)[i]]; ok {
			(*layers)[i] = -(*layers)[i]
		}
	}
}
func addLayerForUnify(layers *[]int, layer int) {
	*layers = append(*layers, layer)
}
