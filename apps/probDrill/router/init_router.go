package router

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"probabDrill/apps/probDrill/apis"
	"strings"
)

func InitRouter() {
	r := gin.Default()
	r.RedirectTrailingSlash = true
	r.Use(apis.CostTime())
	r.LoadHTMLGlob("./*.md")
	r.Handle("GET", "/", func(context *gin.Context) {
		//context.HTML(http.StatusOK, "index.html", nil)
		context.File("./README.md")
	})
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1g := r.Group("/v1")
	geng := v1g.Group("/gen")
	errg := v1g.Group("/err")
	utilsg := v1g.Group("/utils")
	{
		geng.GET("/GenIdwDrill", apis.GenIdwDrill)
		geng.GET("/GenM1Drill", apis.GenM1Drill)
		geng.GET("/GenM1DrillSecond", apis.GenM1DrillSecond)
	}
	{
		errg.GET("/GetStudyAreaAvgDrillPEIdw", apis.GetStudyAreaAvgDrillPEIdw)
		errg.GET("/GetStudyAreaAvgDrillPEM1", apis.GetStudyAreaAvgDrillPEM1)
		utilsg.GET("/GetAvgPEByLayerIDW", apis.GetAvgPEByLayerIDW)
		utilsg.GET("/GetAvgPEByLayerM1", apis.GetAvgPEByLayerM1)
		utilsg.GET("/DrillAroundPeCloud", apis.DrillAroundPeCloud)
		utilsg.GET("/DrillAroundPeCloudM1", apis.DrillAroundPeCloudM1)
	}
	{
		utilsg.GET("/ProbBlocks", apis.ProbBlocks)
		utilsg.GET("/IsValidPoint", apis.IsValidPoint)
		utilsg.GET("/queryDrill", apis.QueryDrill)
		utilsg.GET("/GetRec", apis.GetRec)
	}

	port := strings.Join([]string{viper.GetString("listen.ip"), ":", viper.GetString("listen.port")}, "")
	err := r.Run(port)
	if err != nil {
		return
	}
}
