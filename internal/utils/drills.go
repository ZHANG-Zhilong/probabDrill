package utils

import (
	"probabDrill/internal/entity"
)

func Extend(drill entity.Drill, unifiedSeq []int, drills []entity.Drill) entity.Drill {
	//attention, drill's length question, that drill's length is not under control.
	drill = drill.UnifySeq(unifiedSeq)
	tempDrills := make([]entity.Drill, len(drills))
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

func ExtendDrills(unifiedSeq []int, nearDrills []entity.Drill) (extendedDrills []entity.Drill) {
	for _, d := range nearDrills {
		extendedDrills = append(extendedDrills, Extend(d, unifiedSeq, nearDrills))
	}
	return
}
