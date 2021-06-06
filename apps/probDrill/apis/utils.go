package apis

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gonum.org/v1/gonum/blas/blas64"
	"math"
	"net/http"
	"probabDrill/apps/probDrill/model"
	"probabDrill/internal/constant"
	"probabDrill/internal/utils"
	"strconv"
)

// IsValidPoint godoc
// @Id IsValidPoint
// @Summary 钻孔位置是否有效
// @Description 钻孔位置是否在边界内
// @Tags 工具类
// @Accept mpfd
// @Produce json
// @Param x query number true "x"
// @Param y query number true "y"
// @Success 200 {string} string "{"true"}"
// @Router /utils/IsValidPoint [get]
func IsValidPoint(c *gin.Context) {
	x, _ := strconv.ParseFloat(c.Query("x"), 10)
	fmt.Println(x)
	y, _ := strconv.ParseFloat(c.Query("y"), 10)
	fmt.Println(y)
	bx, by := constant.GetBoundary()
	flag := utils.IsInPolygon(bx, by, x, y)
	req := IsValidPointReq{x, y}
	resp := IsValidPointResp{req, []bool{flag}}
	rst, _ := json.Marshal(resp)
	c.JSON(http.StatusOK, string(rst))
}

type IsValidPointReq struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
type IsValidPointResp struct {
	Param   IsValidPointReq `json:"param"`
	Results []bool          `json:"results"`
}

// GetRec godoc
// @Id GetRec
// @Summary 获取研究区域内左下和右上角点坐标
// @Description 获取研究区域内左下和右上角点坐标[x1,y1,x2,y2]
// @Tags 工具类
// @Produce json
// @Success 200 {string} string
// @Router /utils/GetRec [get]
func GetRec(ctx *gin.Context) {
	x, y := constant.GetBoundary()
	var x1, y1, x2, y2 float64
	for k := 1; k < len(x); k++ {
		x1 = math.Min(x[k-1], x[k])
		x2 = math.Max(x[k-1], x[k])
		y1 = math.Min(y[k-1], y[k])
		y2 = math.Max(y[k-1], y[k])
	}
	resp := GetRecResp{"", []float64{x1, y1, x2, y2}}
	rst, _ := json.Marshal(resp)
	fmt.Println(string(rst))
	ctx.String(http.StatusOK, string(rst))
}

type GetRecReq struct{}
type GetRecResp struct {
	Param   string    `json:"param"`
	Results []float64 `json:"results"`
}

// ProbBlocks godoc
// @Id ProbBlocks
// @Summary 研究区域内block概率矩阵
// @Description 获取研究区域内不同深度范围内blocks出现的概率
// @Tags 工具类
// @Produce json
// @Success 200 {string} string
// @Router /utils/ProbBlocks [get]
func ProbBlocks(ctx *gin.Context) {
	//for diagram  P(block) 和其他等在插值过程中用到的一些条形图
	drillSet := constant.GetHelpDrills()
	blocks := utils.MakeBlocks(drillSet, viper.GetFloat64("BlockResz"))
	pblocks, _ := utils.ProbBlocks(drillSet, blocks)
	rst := pblocks.RawMatrix()
	req := ProbBlocksReq{}
	resp := ProbBlocksResp{req, []blas64.General{rst}}
	rstJson, _ := json.Marshal(resp)
	ctx.String(http.StatusOK, string(rstJson))
}

type ProbBlocksReq struct {
}
type ProbBlocksResp struct {
	Param   ProbBlocksReq    `json:"param"`
	Results []blas64.General `json:"results"`
}

// QueryDrill godoc
// @Id QueryDrill
// @Summary 钻孔位置是否有效
// @Description 钻孔位置是否在边界内
// @Tags 工具类
// @Accept mpfd
// @Produce json
// @Param name query string true "TZZK06"
// @Success 200 {array} model.Drill "return queried drill, TZZK05"
// @Router /utils/queryDrill [get]
func QueryDrill(c *gin.Context) {
	drillName := c.Query("name")
	if d, ok := constant.GetRealDrillByName(drillName); ok {
		var drills []model.Drill
		drills = append(drills, d)
		rep := QueryDrillReq{drillName}
		resp := QueryDrillResp{rep, drills}
		respJson, _ := json.Marshal(resp)
		c.String(http.StatusOK, fmt.Sprintln(string(respJson)))
	} else {
		c.String(http.StatusBadRequest, fmt.Sprintf("invalid drill name: %s.\n", drillName))
	}
}

type QueryDrillReq struct {
	Name string `json:"name"`
}
type QueryDrillResp struct {
	Param   QueryDrillReq `json:"param"`
	Results []model.Drill `json:"results"`
}
