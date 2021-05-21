package utils

import "probabDrill/apps/probDrill/model"

func Extend(drill model.Drill, unifiedSeq []int, drills []model.Drill) model.Drill {
	//attention, drill's length question, that drill's length is not under control.
	drill = drill.UnifySeq(unifiedSeq)
	tempDrills := make([]model.Drill, len(drills))
	copy(tempDrills, drills)
	nearDrills := drill.NearKDrills(tempDrills, 10)
	SetClassicalIdwWeights(drill, nearDrills)

	//find first layer's ceil heights equals drill's bottom
	var tag int
	for idx := 1; idx < len(drill.LayerHeights); idx++ {
		if drill.LayerHeights[idx-1] == drill.BottomHeight() {
			tag = idx
			break
		}
	}
	if tag == 0 {
		return drill
	}
	for idx := tag; idx < len(drill.Layers); idx++ {
		var thickness float64
		for _, d := range nearDrills {
			if thick, ok := d.LayerThickness(drill.Layers[idx]); ok {
				thickness += thick * d.GetWeight()
			}
		}
		drill.LayerHeights[idx] = drill.LayerHeights[idx-1] - thickness
	}
	drill.Update()
	drill.UnBlock()
	return drill
}

func Extend2(drill model.Drill, unifiedSeq []int, drills []model.Drill) model.Drill {
	//attention, drill's length question, that drill's length is not under control.
	drill = drill.UnifySeq(unifiedSeq)
	nearDrills := make([]model.Drill, len(drills))
	copy(nearDrills, drills)
	SetClassicalIdwWeights(drill, nearDrills)

	//find first layer's ceil heights equals drill's bottom
	var tag int
	for idx := 1; idx < len(drill.LayerHeights); idx++ {
		if drill.LayerHeights[idx-1] == drill.BottomHeight() {
			tag = idx
			break
		}
	}
	if tag == 0 {
		return drill
	}
	for idx := tag; idx < len(drill.Layers); idx++ {
		var thickness float64
		for _, d := range nearDrills {
			if thick, ok := d.LayerThickness(drill.Layers[idx]); ok {
				thickness += thick * d.GetWeight()
			}
		}
		drill.LayerHeights[idx] = drill.LayerHeights[idx-1] - thickness
	}
	drill.Update()
	drill.UnBlock()
	return drill
}
func ExtendDrills(unifiedSeq []int, drills []model.Drill) (extendedDrills []model.Drill) {
	for idx, d := range drills {
		if idx-1 >= 0 && idx+1 < len(drills) {
			extendedDrills = append(extendedDrills, Extend2(d, unifiedSeq, []model.Drill{drills[idx-1], drills[idx+1]}))
			continue
		}
		if idx == 0 && idx+1 < len(drills) {
			extendedDrills = append(extendedDrills, Extend2(d, unifiedSeq, []model.Drill{drills[idx+1]}))
			continue
		}
		if idx == len(drills)-1 && idx-1 >= 0 {
			extendedDrills = append(extendedDrills, Extend2(d, unifiedSeq, []model.Drill{drills[idx-1]}))
			continue
		}
	}
	return
}
