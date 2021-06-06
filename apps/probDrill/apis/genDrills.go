package apis

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"probabDrill/apps/probDrill/model"
	"probabDrill/internal/constant"
	"probabDrill/internal/service"
	"probabDrill/internal/utils"
	"strconv"
	"time"
)

func CostTime() gin.HandlerFunc {
	return func(c *gin.Context) {
		nowTime := time.Now()
		c.Next() // 处理handler
		ct := time.Since(nowTime)
		url := c.Request.URL.String()
		fmt.Printf("the request URL %s cost %v\n", url, ct)
	}
}

// GenIdwDrill godoc
// @Id GenIdwDrill
// @Summary 生成反距离插值的虚拟钻孔数据
// @Description 在特定位置生成钻孔数据，采用反距离插值方法
// @Tags 服务
// @Param x query float32 true "position"
// @Param y query float32 true "position"
// @Produce  json
// @Success 200 {string} string
// @Router /gen/GenIdwDrill [get]
func GenIdwDrill(c *gin.Context) {
	x, _ := strconv.ParseFloat(c.Query("x"), 10)
	y, _ := strconv.ParseFloat(c.Query("y"), 10)
	drillSet := constant.GetHelpDrills()
	drill := service.GenVDrillIDW(drillSet, nil, nil, x, y)
	var drills []model.Drill
	drills = append(drills, drill)

	req := GenIdwDrillReq{x, y}
	resp := GenIdwDrillResp{req, drills}
	resultJson, _ := json.Marshal(resp)
	c.String(http.StatusOK, fmt.Sprintf("%s", string(resultJson)))
}

type GenIdwDrillReq struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
type GenIdwDrillResp struct {
	Param   GenIdwDrillReq `json:"param"`
	Results []model.Drill  `json:"results"`
}

// GenM1Drill godoc
// @Id GenM1Drill
// @Summary 生成概率模型虚拟钻孔数据
// @Description 在特定位置生成钻孔数据，采用基于沉积序列的三维地层概率模型构建方法
// @Tags 服务
// @Param x query float32 true "position"
// @Param y query float32 true "position"
// @Produce  json
// @Success 200 {string} string
// @Router /gen/GenM1Drill [get]
func GenM1Drill(c *gin.Context) {
	x, _ := strconv.ParseFloat(c.Query("x"), 10)
	y, _ := strconv.ParseFloat(c.Query("y"), 10)
	drillSet := constant.GetHelpDrills()
	blocks := utils.MakeBlocks(drillSet, viper.GetFloat64("BlockResZ"))
	pBlockLayerMat, _ := utils.ProbBlockLayerMatG(drillSet, blocks)
	drill := service.GenVDrillM1(drillSet, blocks, pBlockLayerMat, x, y)
	var drills []model.Drill
	drills = append(drills, drill)
	req := GenM1DrillReq{x, y}
	resp := GenM1DrillResp{req, drills}
	resultJson, _ := json.Marshal(resp)
	c.String(http.StatusOK, fmt.Sprintf("%s", string(resultJson)))
}

type GenM1DrillReq struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
type GenM1DrillResp struct {
	Param   GenM1DrillReq `json:"param"`
	Results []model.Drill `json:"results"`
}

// GenM1DrillSecond godoc
// @Id GenM1DrillSecond
// @Summary 生成次大概率模型虚拟钻孔数据
// @Description 在特定位置生成钻孔数据，采用基于沉积序列的三维地层概率模型构建方法
// @Tags 服务
// @Param x query float32 true "position"
// @Param y query float32 true "position"
// @Produce  json
// @Success 200 {string} string
// @Router /gen/GenM1DrillSecond [get]
func GenM1DrillSecond(c *gin.Context) {
	x, _ := strconv.ParseFloat(c.Query("x"), 10)
	y, _ := strconv.ParseFloat(c.Query("y"), 10)
	drillSet := constant.GetHelpDrills()
	blocks := utils.MakeBlocks(drillSet, viper.GetFloat64("BlockResZ"))
	pBlockLayerMat, _ := utils.ProbBlockLayerMatG(drillSet, blocks)
	drill := service.GenVDrillM1Second(drillSet, blocks, pBlockLayerMat, x, y)
	var drills []model.Drill
	drills = append(drills, drill)
	req := GenM1DrillReq{x, y}
	resp := GenM1DrillResp{req, drills}
	resultJson, _ := json.Marshal(resp)
	c.String(http.StatusOK, fmt.Sprintf("%s", string(resultJson)))
}
