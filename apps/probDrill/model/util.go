package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"os/exec"
	"strconv"
)

func printSliceFloat64(s []float64) () {
	fmt.Print("[")
	for _, v := range s {
		fmt.Printf("%+2.2f\t", v)
	}
	fmt.Print("]\n")
}
func printSliceInt(s []int) () {
	fmt.Print("[")
	for _, v := range s {
		fmt.Printf("%4d\t", v)
	}
	fmt.Print("]\n")
}
func SetLengthAndZ(drill *Drill, nearDrills []Drill) {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	//这里，要检查 vdrill的 上顶点和下顶点要小于周围钻孔的最大值和最小值
	var maxLength, maxHeight, minHeight = -math.MaxFloat64, -math.MaxFloat64, math.MaxFloat64
	var length, z float64 = 0, 0

	//获取周围钻孔的最大高程和底部最小高程
	for _, d := range nearDrills {
		maxHeight = math.Max(d.Z, maxHeight)
		minHeight = math.Min(d.BottomHeight(), minHeight)
		maxLength = math.Max(d.GetLength(), maxLength)
	}

	//cal the z value and length with weight.
	for idx := 0; idx < len(nearDrills); idx++ {
		length += nearDrills[idx].GetLength() * nearDrills[idx].GetWeight()
		z += nearDrills[idx].Z * nearDrills[idx].GetWeight()
	}

	//确保参数在极值之内
	if length > maxLength {
		length = maxLength
	}
	if z > maxHeight {
		z = maxHeight
	}

	drill.SetLength(decimal(length))
	drill.SetZ(decimal(z))
	if drill.BottomHeight() < minHeight {
		drill.SetLength(decimal(drill.Z - minHeight))
	}
	if drill.Z <= drill.BottomHeight() {
		//May be drill.length==0??
		log.Printf("error,drill.Z <= drill.BottomHeight() ")
		log.Printf("rst:drill.Z:%f, bottom:%f, bottom height:%f, drill:%#v \n", drill.Z, minHeight, drill.BottomHeight(), drill)
		for _, d := range nearDrills {
			log.Printf("rst:drill.Z:%f, bottom:%f, incident drill: %#v\n", d.Z, d.BottomHeight(), d)
		}
		os.Exit(-1)
	}
	drill.LayerHeights[0] = drill.Z

}
func decimal(value float64) float64 {
	value = math.Trunc(value*1e2+0.5) * 1e-2
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}

func workout3dModel(drills []Drill) {
	fileJson, _ := json.Marshal(drills)
	err := ioutil.WriteFile("drill.json", fileJson, 0644)
	if err != nil {
		fmt.Printf("WriteFileJson ERROR: %+v", err)
	}
	command := ` ./drill.json`
	cmd := exec.Command("workoutModel", "-d", command)

	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Execute Shell:%s failed with error:%s", command, err.Error())
		return
	}
	fmt.Printf("Execute Shell:%s finished with output:\n%s", command, string(output))
}
