package probabDrill

const GridXY float64 = 1000   //distance of grid
const BlockResZ float64 = 0.7 //length of block < 1/3 min layer thickness  >0.1
const OffX float64 = -3592727.0
const OffY float64 = -499523.0
const ScaleXY float64 = 1
const ScaleZ float64 = 1
const StdLen int = 38  //len of std layer num
const RadiusIn int = 5 //how many real drills in search radius
const SearchRadius float64 = 3000
const IdwPow float64 = 2 //power of inverse weighting interpolation method
const MinThicknessInDrill =1
const MinDrillDist = 1

const Layer string = "/Users/zhangzhilong/go/src/probabDrill/assets/layer_info.dat"
const Basic string = "/Users/zhangzhilong/go/src/probabDrill/assets/basic_info.dat"
const Boundary string = "/Users/zhangzhilong/go/src/probabDrill/assets/boundary_info.dat"
const StdLayer string = "/Users/zhangzhilong/go/src/probabDrill/assets/std_layer_info.dat"

const DrillWidth int = 20
const CanvasWidth int = 1000
const CanvasHeight int = 700
const CanvasMinThickness int = 1
const CanvasOffsetY int = 0

//poissiondisc sample param for help drill set
const MinDistance = 500  // min distance between points
const MaxAttemptAdd = 10 // max attempts to add neighboring point

//const Basic string = "C:/Users/张志龙/go/src/probabDrill/assets/basic_info.dat"
//const Layer string = "C:/Users/张志龙/go/src/probabDrill/assets/layer_info.dat"
//const Boundary string = "C:/Users/张志龙/go/src/probabDrill/assets/boundary_info.dat"
//const StdLayer string = "C:/Users/张志龙/go/src/probabDrill/assets/std_layer_info.dat"
