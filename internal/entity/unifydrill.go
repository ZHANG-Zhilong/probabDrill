package entity

import "log"

func UnifyDrillsStrata(drills []Drill, check func([]int) []int) []Drill {
	if len(drills) < 2 {
		log.SetFlags(log.Lshortfile)
		log.Fatal("error")
		return nil
	}
	var repeat, nonRepeat []Drill
	var ps [][]int
	for _, drill := range drills {
		if p := RepeatPattern(drill.Layers); p != nil {
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
			stdLayers = GetUnifiedSeq(layers, stdLayers, check)
		}
	}
	if len(nonRepeat) > 0 {
		if len(stdLayers) == 0 { // no repeat layer in drill.
			stdLayers = nonRepeat[0].Layers
		}
		for idx := 0; idx < len(nonRepeat); idx++ {
			layers := nonRepeat[idx].Layers
			MarkRepeat(layers, ps)
			stdLayers = GetUnifiedSeq(layers, stdLayers, check)
		}
	}
	var stdDrills []Drill
	for _, d := range drills {
		stdDrills = append(stdDrills, d.StdSeq(stdLayers))
	}
	return stdDrills
}
func GetUnifiedSeq(seq1, seq2 []int, check func([]int) []int) (seq []int) {
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
		if IsNormal(&layers, i) { //normal first
			layerMap[(layers)[i]] = 1
			continue
		}
		if IsLack(&layers, i) { //lack first
			layerMap[(layers)[i]] = 1
			continue
		}
		if IsInverse(&layers, i) {
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
	return CheckSeqMinNegR(layers, 1, layerMap, true)
}
func CheckSeqMinNegR(layers []int, start int, layerMap map[int]int, lackFirst bool) (checkedSeq []int) {
	if start == len(layers) {
		cseq := make([]int, len(layers))
		copy(cseq, layers)
		delete(layerMap, layers[start-1])
		return cseq
	}
	//check repeat first.
	if _, ok := layerMap[layers[start]]; ok {
		layers[start] = -layers[start]
		seq := CheckSeqMinNegR(layers, start+1, layerMap, lackFirst)
		layers[start] = -layers[start]
		delete(layerMap, layers[start])
		return seq
	}

	if IsNormal(&layers, start) { //normal first
		layerMap[layers[start]] = 1
		return CheckSeqMinNegR(layers, start+1, layerMap, lackFirst)
	}

	if IsLack(&layers, start) && IsInverse(&layers, start) {
		layerMap[layers[start]] = 1
		//regard as lack
		seq1 := CheckSeqMinNegR(layers, start+1, layerMap, lackFirst)
		case1 := CountNeg(&seq1)

		//regard as inverse
		layers[start] = -layers[start]
		seq2 := CheckSeqMinNegR(layers, start+1, layerMap, lackFirst)
		case2 := CountNeg(&seq2)

		layers[start] = -layers[start]
		delete(layerMap, layers[start])

		if lackFirst && (case1 <= case2) || lackFirst && (case1 < case2) { //lack first <=
			return seq1
		} else {
			return seq2
		}
	}

	if IsLack(&layers, start) { //lack first
		layerMap[layers[start]] = 1
		return CheckSeqMinNegR(layers, start+1, layerMap, lackFirst)
	}
	if IsInverse(&layers, start) {
		layerMap[layers[start]] = 1
		layers[start] = -layers[start]
		seq := CheckSeqMinNegR(layers, start+1, layerMap, lackFirst)
		layers[start] = -layers[start]
		delete(layerMap, layers[start])
		return seq
	}
	return
}
func RepeatPattern(seq []int) (p [][]int) {
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
func CountNeg(seq *[]int) (count int) {
	for _, val := range *seq {
		if val < 0 {
			count++
		}
	}
	return
}
func MarkRepeat(seq []int, pattern [][]int) {
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
func IsNormal(layers *[]int, idx int) (ok bool) {
	if idx <= len(*layers)-1 && idx-1 >= 0 {
		return LastPositive(layers, idx)+1 == (*layers)[idx]
	}
	if idx >= 1 && idx+1 <= len(*layers)-1 {
		return (*layers)[idx]+1 == (*layers)[idx+1]
	}
	return
}
func IsLack(layers *[]int, idx int) (ok bool) {
	if idx == len(*layers)-1 {
		return (*layers)[idx] > LastPositive(layers, idx)+1
	}
	if idx+1 <= len(*layers)-1 && idx-1 >= 1 {
		return (*layers)[idx] > LastPositive(layers, idx)+1 && (*layers)[idx+1] > (*layers)[idx]
	}
	return
}
func IsInverse(layers *[]int, idx int) (ok bool) {
	//case 1: drill top, ->   😄  smaller  go top
	if idx == 1 && idx+1 < len(*layers) {
		if (*layers)[idx] > (*layers)[idx+1] {
			return true
		}
	}
	//case 2: drill bottom -> bigger 😄  go bottom
	if idx == len(*layers)-1 && idx-1 >= 1 {
		if LastPositive(layers, idx) > (*layers)[idx] {
			return true
		}
	}

	//case 3: drill middle -> bigger 😄 bigger  go down
	//case 4: drill middle -> bigger 😄 smaller  go down
	//case 5: drill middle -> smaller 😄 smaller  go up
	if idx-1 >= 1 && idx+1 < len(*layers) {

		if LastPositive(layers, idx) > (*layers)[idx] {
			return true
		}

		if LastPositive(layers, idx)+1 < (*layers)[idx] &&
			(*layers)[idx] > (*layers)[idx+1] {
			return true
		}
	}

	return false
}
func LastPositive(layers *[]int, idx int) (layer int) {
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
