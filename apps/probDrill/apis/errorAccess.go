package apis

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"probabDrill/internal/constant"
	"probabDrill/internal/service"
)

// GetStudyAreaAvgDrillPEIdw godoc
// @Id GetStudyAreaAvgDrillPEIdw
// @Summary 获取研究区域内钻孔平均百分比误差反距离插值方法
// @Description 获取研究区域内钻孔平均百分比误差，需要自行获取研究区域的矩形边界
// @Tags 误差评估
// @Produce json
// @Success 200 {string} string
// @Router /err/GetStudyAreaAvgDrillPEIdw [get]
func GetStudyAreaAvgDrillPEIdw(ctx *gin.Context) {
	//通过多次，计算，获取研究区域内不同位置虚拟钻孔Pe
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	drills := constant.GetHelpDrills()
	if x, y, pe, err := service.GetPeByArea(drills, 10, service.GenVDrillsIDW, 10); err == nil {
		resp := GetStudyAreaAvgDrillPEResp{"", [][]float64{x, y, pe}}
		rst, _ := json.Marshal(resp)
		ctx.JSON(http.StatusOK, string(rst))
	} else {
		log.Println("TestGetPeAvgByLayer", err)
	}
}

// GetStudyAreaAvgDrillPEM1 godoc
// @Id GetStudyAreaAvgDrillPEM1
// @Summary 获取研究区域内钻孔平均百分比误差概率模型构建方法
// @Description 获取研究区域内钻孔平均百分比误差，需要自行获取研究区域的矩形边界
// @Tags 误差评估
// @Produce json
// @Success 200 {string} string
// @Router /err/GetStudyAreaAvgDrillPEM1 [get]
func GetStudyAreaAvgDrillPEM1(ctx *gin.Context) {
	//通过多次，计算，获取研究区域内不同位置虚拟钻孔Pe
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	drills := constant.GetHelpDrills()
	if x, y, pe, err := service.GetPeByArea(drills, 10, service.GenVDrillsM1, 10); err == nil {
		resp := GetStudyAreaAvgDrillPEResp{"", [][]float64{x, y, pe}}
		rst, _ := json.Marshal(resp)
		ctx.JSON(http.StatusOK, string(rst))
	} else {
		log.Println("TestGetPeAvgByLayer", err)
	}
}

type GetStudyAreaAvgDrillPEResp struct {
	Param   string      `json:"param"`
	Results [][]float64 `json:"results"`
}

// GetAvgPEByLayerIDW godoc
// @Id GetAvgPEByLayerIDW
// @Summary 钻孔分层PE
// @Description 获取研究区域内钻孔钻孔分层PE，应用反距离插值方法
// @Tags 误差评估
// @Produce json
// @Success 200 {string} string
// @Router /err/GetAvgPEByLayerIDW [get]
func GetAvgPEByLayerIDW(ctx *gin.Context) {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	//按地质分层，分别计算每一个分层的平均百分比误差
	drills := constant.GetHelpDrills()
	//drills := constant.GetRealDrills()
	if pes, err := service.GetPeAvgByLayer(drills, 10, service.GenVDrillsIDW); err == nil {
		resp := GetStudyAreaAvgDrillPEResp{"", [][]float64{pes}}
		rst, _ := json.Marshal(resp)
		ctx.JSON(http.StatusOK, string(rst))
	} else {
		log.Println("TestGetPeAvgByLayer", err)
	}
}

// GetAvgPEByLayerM1 godoc
// @Id GetAvgPEByLayerM1
// @Summary 钻孔分层PE
// @Description 获取研究区域内钻孔钻孔分层PE，应用反距离插值方法
// @Tags 误差评估
// @Produce json
// @Success 200 {string} string
// @Router /err/GetAvgPEByLayerM1 [get]
func GetAvgPEByLayerM1(ctx *gin.Context) {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	//按地质分层，分别计算每一个分层的平均百分比误差
	drills := constant.GetHelpDrills()
	//drills := constant.GetRealDrills()
	if pes, err := service.GetPeAvgByLayer(drills, 10, service.GenVDrillsM1); err == nil {
		resp := GetStudyAreaAvgDrillPEResp{"", [][]float64{pes}}
		rst, _ := json.Marshal(resp)
		ctx.JSON(http.StatusOK, string(rst))
	} else {
		log.Println("TestGetPeAvgByLayer", err)
	}
}

// DrillAroundPeCloud godoc
// @Id DrillAroundPeCloud
// @Summary 真实钻孔周围500米虚拟钻孔的百分比误差
// @Description 真实钻孔周围虚拟钻孔的百分比误差，应用反距离插值方法 [][]float64{xs, ys, pes}
// @Tags 误差评估
// @Param name query string true "TZZK06"
// @Produce json
// @Success 200 {string} string
// @Router /err/DrillAroundPeCloud [get]
func DrillAroundPeCloud(ctx *gin.Context) {
	drillName := ctx.Query("name")
	drillSet := constant.GetHelpDrills()
	realDrill, _ := constant.GetRealDrillByName(drillName)
	var xs, ys, pes []float64
	for y := float64(-500); y <= 500; y += 50 {
		for x := float64(-500); x <= 500; x += 50 {
			virtualDrill := service.GenVDrillIDW(drillSet, nil, nil, realDrill.X+x, realDrill.Y+y)
			pe := service.GetPeByDrill(realDrill, virtualDrill)
			xs = append(xs, x)
			ys = append(ys, y)
			pes = append(pes, pe)
		}
	}
	resp := GetStudyAreaAvgDrillPEResp{"", [][]float64{xs, ys, pes}}
	rst, _ := json.Marshal(resp)
	ctx.JSON(http.StatusOK, string(rst))
}

// DrillAroundPeCloudM1 godoc
// @Id DrillAroundPeCloudM1
// @Summary 真实钻孔周围500米虚拟钻孔的百分比误差
// @Description 真实钻孔周围虚拟钻孔的百分比误差，应用反距离插值方法 [][]float64{xs, ys, pes}
// @Tags 误差评估
// @Param name query string true "TZZK06"
// @Produce json
// @Success 200 {string} string
// @Router /err/DrillAroundPeCloudM1 [get]
func DrillAroundPeCloudM1(ctx *gin.Context) {
	drillName := ctx.Query("name")
	drillSet := constant.GetHelpDrills()
	realDrill, _ := constant.GetRealDrillByName(drillName)
	var xs, ys, pes []float64
	for y := float64(-500); y <= 500; y += 50 {
		for x := float64(-500); x <= 500; x += 50 {
			virtualDrill := service.GenVDrillIDW(drillSet, nil, nil, realDrill.X+x, realDrill.Y+y)
			pe := service.GetPeByDrill(realDrill, virtualDrill)
			xs = append(xs, x)
			ys = append(ys, y)
			pes = append(pes, pe)
		}
	}
	resp := GetStudyAreaAvgDrillPEResp{"", [][]float64{xs, ys, pes}}
	rst, _ := json.Marshal(resp)
	ctx.JSON(http.StatusOK, string(rst))
}
