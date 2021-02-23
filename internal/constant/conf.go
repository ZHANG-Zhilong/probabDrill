package constant

import (
	"encoding/json"
	"fmt"
	"log"
)

type params struct {
	Basic       string
	Layer       string
	Boundary    string `json:"bound"`
	StdLayer    string
	Offsetx     float64
	Offsety     float64
	Scalex      float64
	Scaly       float64
	Scalez      float64
	Resolutionz float64
}

func play() {
	params2 := params{"a", "b", "c", "d", 1.0, 1.2, 1.3, 1.4, 1.5, 1.6}
	data, err := json.Marshal(params2)
	if err != nil {
		log.Fatalf("Json marshaling failed: %s", err)
	}
	str := `{"Basic":"a","Layer":"b","bound":"c","StdLayer":"d","Offsetx":1,"Offsety":1.2,"Scalex":1.3,"Scaly":1.4,"Scalez":1.5,"Resolutionz":1.6}`
	var stu params
	if err := json.Unmarshal([]byte(str), &stu); err != nil {
		log.Fatalf("Json Unmarshaling failed: %s", err)
	}
	fmt.Println(params2)
	fmt.Printf("%s\n", data)
	fmt.Println(stu)
}

//ResX resolution x
const ResXY float64 = 100

//ResZ resolution z
//z方向的精度要求，应当小于1/3的最小层厚
const ResZ float64 = 0.5

//OffX offsetX
const OffX float64 = -3592727.0

//OffY offsetY
const OffY float64 = -499523.0

//ScaleXY scales
const ScaleXY float64 = 1

//ScaleZ scales
const ScaleZ float64 = 1

//Sr search radius
const Sr float64 = 100

//const Basic string = "/Users/zhangzhilong/go/src/probabDrill-main/assets/basic_info.dat"
const Basic string = "C:/Users/张志龙/go/src/probabDrill-main/assets/basic_info.dat"
//const Layer string = "/Users/zhangzhilong/go/src/probabDrill-main/assets/layer_info.dat"
const Layer string = "C:/Users/张志龙/go/src/probabDrill-main/assets/layer_info.dat"
//const Boundary string = "/Users/zhangzhilong/go/src/probabDrill-main/assets/boundary_info.dat"
const Boundary string = "C:/Users/张志龙/go/src/probabDrill-main/assets/boundary_info.dat"
//const StdLayer string = "/Users/zhangzhilong/go/src/probabDrill-main/assets/std_layer_info.dat"
const StdLayer string = "C:/Users/张志龙/go/src/probabDrill-main/assets/std_layer_info.dat"

const StdLen int64 = 38

const RadiusIn int = 4

const IdwPow float64 = 2

