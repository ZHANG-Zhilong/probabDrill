package utils

import (
	"fmt"
	"log"
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
func UnifyDrillsStrata(drills *[]entity.Drill, check func([]int) []int) *[]entity.Drill {
	if len(*drills) < 2 {
		log.SetFlags(log.Lshortfile)
		log.Fatal("error")
		return nil
	}
	var repeat, nonRepeat []entity.Drill
	var ps [][]int
	for _, drill := range *drills {
		if p := repeatPattern(drill.Layers); p != nil {
			repeat = append(repeat, drill)
			ps = append(ps, p...)
		} else {
			nonRepeat = append(nonRepeat, drill)
		}
	}
	var stdLayers []int
	if len(repeat) > 0 {
		stdLayers = repeat[0].Layers
		for idx := 0; idx < len(repeat); idx++ {
			layers := repeat[idx].Layers
			stdLayers = getUnifiedSeq(layers, stdLayers, check)
		}
	}
	if len(nonRepeat) > 0 {
		if len(stdLayers) == 0 { // no repeat layer in drill.
			stdLayers = nonRepeat[0].Layers
		}
		for idx := 0; idx < len(nonRepeat); idx++ {
			layers := nonRepeat[idx].Layers
			markRepeat(layers, ps)
			stdLayers = getUnifiedSeq(layers, stdLayers, check)
		}
	}
	var stdDrills []entity.Drill
	for _, d := range *drills {
		stdDrills = append(stdDrills, d.StdSeq(stdLayers))
	}
	return &stdDrills
}
func getUnifiedSeq(seq1, seq2 []int, check func([]int) []int) (seq []int) {
	idx1, idx2 := 1, 1
	seq = []int{0}
	check(seq1)
	check(seq2)
	for {
		if idx1 == len(seq1) && idx2 == len(seq2) {
			break
		}
		if idx2 == len(seq2) && idx1 < len(seq1) ||
			idx1 < len(seq1) && idx2 < len(seq2) && seq1[idx1] < seq2[idx2] {
			seq = append(seq, seq1[idx1])
			idx1++
			continue
		}
		if idx1 == len(seq1) && idx2 < len(seq2) ||
			idx1 < len(seq1) && idx2 < len(seq2) && seq1[idx1] > seq2[idx2] {
			seq = append(seq, seq2[idx2])
			idx2++
			continue
		}
		if seq1[idx1] == seq2[idx2] {
			seq = append(seq, seq1[idx1])
			idx1++
			idx2++
			continue
		}
	}
	return seq
}
func CheckSeqZiChun(layers []int) (seq []int) {
	for _, layer := range layers {
		if layer < 0 {
			return
		}
	}
	//mark inverse and repeat
	var layerMap map[int]int = make(map[int]int)
	for i := 1; i < len(layers); i++ {
		if isNormal(&layers, i) { //normal first
			layerMap[(layers)[i]] = 1
			continue
		}
		if isLack(&layers, i) { //lack first
			layerMap[(layers)[i]] = 1
			continue
		}
		if isInverse(&layers, i) {
			layers[i] = -layers[i]
			layerMap[layers[i]] = 1
			continue
		}
		if _, ok := layerMap[layers[i]]; ok {
			layers[i] = -layers[i]
		}
	}
	seq = make([]int, len(layers), len(layers))
	copy(seq, layers)
	return
}
func CheckSeqMinNeg(layers []int) (checkedSeq []int) {
	layerMap := make(map[int]int)
	return checkSeqMinNeg(layers, 1, layerMap, true)
}
func checkSeqMinNeg(layers []int, start int, layerMap map[int]int, lackFirst bool) (checkedSeq []int) {
	if start == len(layers) {
		cseq := make([]int, len(layers))
		copy(cseq, layers)
		delete(layerMap, layers[start-1])
		return cseq
	}
	//check repeat first.
	if _, ok := layerMap[layers[start]]; ok {
		layers[start] = -layers[start]
		seq := checkSeqMinNeg(layers, start+1, layerMap, lackFirst)
		layers[start] = -layers[start]
		delete(layerMap, layers[start])
		return seq
	}

	if isNormal(&layers, start) { //normal first
		layerMap[layers[start]] = 1
		return checkSeqMinNeg(layers, start+1, layerMap, lackFirst)
	}

	if isLack(&layers, start) && isInverse(&layers, start) {
		layerMap[layers[start]] = 1
		//regard as lack
		seq1 := checkSeqMinNeg(layers, start+1, layerMap, lackFirst)
		case1 := countNeg(&seq1)

		//regard as inverse
		layers[start] = -layers[start]
		seq2 := checkSeqMinNeg(layers, start+1, layerMap, lackFirst)
		case2 := countNeg(&seq2)

		layers[start] = -layers[start]
		delete(layerMap, layers[start])

		if lackFirst && (case1 <= case2) || lackFirst && (case1 < case2) { //lack first <=
			return seq1
		} else {
			return seq2
		}
	}

	if isLack(&layers, start) { //lack first
		layerMap[layers[start]] = 1
		return checkSeqMinNeg(layers, start+1, layerMap, lackFirst)
	}
	if isInverse(&layers, start) {
		layerMap[layers[start]] = 1
		layers[start] = -layers[start]
		seq := checkSeqMinNeg(layers, start+1, layerMap, lackFirst)
		layers[start] = -layers[start]
		delete(layerMap, layers[start])
		return seq
	}
	return
}
func repeatPattern(seq []int) (p [][]int) {
	var repeatIdx []int
	layerMap := make(map[int]int)
	for idx, l := range seq {
		if _, ok := layerMap[l]; ok {
			repeatIdx = append(repeatIdx, idx)
		} else {
			layerMap[l] = 1
		}
	}
	if len(repeatIdx) == 0 {
		return nil
	} else {
		for _, val := range repeatIdx {
			if val == len(seq)-1 {
				p = append(p, []int{seq[val-1], seq[val]})
			}
			if val+1 <= len(seq)-1 {
				p = append(p, []int{seq[val-1], seq[val], seq[val+1]})
			}
		}
		return p
	}
}
func countNeg(seq *[]int) (count int) {
	for _, val := range *seq {
		if val < 0 {
			count++
		}
	}
	return
}
func markRepeat(seq []int, pattern [][]int) {
	for idx := 1; idx < len(seq); idx++ {
		for _, p := range pattern {
			if seq[idx] != p[1] {
				continue
			} else {
				if len(p) == 2 && idx == len(seq)-1 && seq[idx-1] == p[0] {
					seq[idx] = -seq[idx]
				}
				if len(p) == 3 && idx+1 <= len(seq)-1 && seq[idx-1] == p[0] && seq[idx+1] == p[2] {
					seq[idx] = -seq[idx]
				}
			}
		}
	}
}
func isNormal(layers *[]int, idx int) (ok bool) {
	if idx <= len(*layers)-1 && idx-1 >= 0 {
		return lastPositive(layers, idx)+1 == (*layers)[idx]
	}
	if idx >= 1 && idx+1 <= len(*layers)-1 {
		return (*layers)[idx]+1 == (*layers)[idx+1]
	}
	return
}
func isLack(layers *[]int, idx int) (ok bool) {
	if idx == len(*layers)-1 {
		return (*layers)[idx] > lastPositive(layers, idx)+1
	}
	if idx+1 <= len(*layers)-1 && idx-1 >= 1 {
		return (*layers)[idx] > lastPositive(layers, idx)+1 && (*layers)[idx+1] > (*layers)[idx]
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
		if lastPositive(layers, idx) > (*layers)[idx] {
			return true
		}
	}

	//case 3: drill middle -> bigger 😄 bigger  go down
	//case 4: drill middle -> bigger 😄 smaller  go down
	//case 5: drill middle -> smaller 😄 smaller  go up
	if idx-1 >= 1 && idx+1 < len(*layers) {

		if lastPositive(layers, idx) > (*layers)[idx] {
			return true
		}

		if lastPositive(layers, idx)+1 < (*layers)[idx] &&
			(*layers)[idx] > (*layers)[idx+1] {
			return true
		}
	}

	return false
}
func lastPositive(layers *[]int, idx int) (layer int) {
	for {
		if idx-1 == 0 {
			break
		}
		if (*layers)[idx-1] > 0 {
			return (*layers)[idx-1]
		} else {
			idx--
		}
	}
	return 0
}
func GetLine(x1, y1, x2, y2, x float64) (y float64) {
	log.SetFlags(log.Lshortfile)
	if x2 == x1 {
		log.Fatal("error")
	}
	y = (x-x1)*(y2-y1)/(x2-x1) + y1
	return
}
func SplitSegment(x1, y1, x2, y2 float64, n int) (vertices []float64) {
	step := (x2 - x1) / float64(n+1)
	for x := x1 + step; math.Abs(x-x1) < math.Abs(x2-x1); x += step {
		y := GetLine(x1, y1, x2, y2, x)
		vertices = append(vertices, x, y)
	}
	return
}
