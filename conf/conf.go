package probabDrill

type Config struct {
	GridXY              float64
	MinThicknessInDrill float64
	ScaleXY             float64
	ScaleZ              float64
	RadiusIn            int
	SearchRadius        float64
	IdwPow              float64
	MinDrillDist        float64

	DrillWidth         int
	CanvasWidth        int
	CanvasHeight       int
	CanvasMinThickness int
	CanvasOffsetY      int

	//泰州项目
	BlockResZ     float64
	OffX          float64
	OffY          float64
	OffZ          float64
	StdLen        int
	MinDistance   float64
	MaxAttemptAdd float64
	RealDrillDist float64
	Layer         string
	Basic         string
	Boundary      string
	StdLayer      string
}


const GridXY float64 = 1000 //distance of grid
const MinThicknessInDrill = 0.3
const ScaleXY float64 = 1
const ScaleZ float64 = 1
const RadiusIn int = 5 //how many real drills in search radius
const SearchRadius float64 = 3000
const IdwPow float64 = 2 //power of inverse weighting interpolation method
const MinDrillDist = 1

const DrillWidth int = 20
const CanvasWidth int = 1000
const CanvasHeight int = 700
const CanvasMinThickness int = 1
const CanvasOffsetY int = 0

//泰州项目
const BlockResZ float64 = 0.3 //length of block < 1/3 min layer thickness  >0.1
const OffX float64 = -3592727.0
const OffY float64 = -499523.0
const OffZ float64 = 0
const StdLen int = 38   //len of std layer num
const MinDistance = 700 // 泊松采样参数，min distance between points
const MaxAttemptAdd = 7 //  泊松采样参数， max attempts to add neighboring point
const RealDrillDist = 90000
const Layer string = "/Users/zhangzhilong/go/src/probabDrill/assets/layer_info.dat"
const Basic string = "/Users/zhangzhilong/go/src/probabDrill/assets/basic_info.dat"
const Boundary string = "/Users/zhangzhilong/go/src/probabDrill/assets/boundary_info.dat"
const StdLayer string = "/Users/zhangzhilong/go/src/probabDrill/assets/std_layer_info.dat"

//阀厅项目
//const BlockResZ float64 = 0.2 //length of block < 1/3 min layer thickness  >0.1
//const StdLen int = 11 //len of std layer num
//const OffX float64 = -2624079
//const OffY float64 = -515645
//const OffZ float64 = -100 //fating
//const MinDistance =  10  // 泊松采样参数，min distance between points
//const MaxAttemptAdd = 15  //  泊松采样参数， max attempts to add neighboring point
//const RealDrillDist = 1000
//const Layer string = "/Users/zhangzhilong/go/src/probabDrill/assets/fating/layer_info.dat"
//const Basic string = "/Users/zhangzhilong/go/src/probabDrill/assets/fating/basic_info.dat"
//const Boundary string = "/Users/zhangzhilong/go/src/probabDrill/assets/fating/boundary_info.dat"
//const StdLayer string = "/Users/zhangzhilong/go/src/probabDrill/assets/fating/std_layer_info.dat"

//const Basic string = "C:/Users/张志龙/go/src/probabDrill/assets/basic_info.dat"
//const Layer string = "C:/Users/张志龙/go/src/probabDrill/assets/layer_info.dat"
//const Boundary string = "C:/Users/张志龙/go/src/probabDrill/assets/boundary_info.dat"
//const StdLayer string = "C:/Users/张志龙/go/src/probabDrill/assets/std_layer_info.dat"
