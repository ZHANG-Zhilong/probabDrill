package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"probabDrill/apps/probDrill/apis"
)

func InitRouter() {
	r := gin.Default()
	r.Use(apis.CostTime())
	r.LoadHTMLGlob("template/*")
	r.Handle("GET", "/", func(context *gin.Context) {
		context.HTML(http.StatusOK, "index.html", nil)
	})

	v1g := r.Group("/v1")
	{
		v1g.GET("/genDrill", apis.GenIdwDrill)
		v1g.POST("/genDrill", apis.GenIdwDrill2)

		v1g.GET("/queryDrill", apis.QueryDrill)
		v1g.POST("/queryDrill", apis.QueryDrill2)
	}
	err := r.Run(":8080")
	if err != nil {
		return
	}
}
