package router

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"net/http"
	"probabDrill/apps/probDrill/apis"
	"strings"
)

func InitRouter() {
	r := gin.Default()
	r.RedirectTrailingSlash = true
	r.Use(apis.CostTime())
	r.LoadHTMLGlob("template/*")
	r.Handle("GET", "/", func(context *gin.Context) {
		context.HTML(http.StatusOK, "index.html", nil)
	})
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1g := r.Group("/v1")
	{
		v1g.GET("/genDrill", apis.GenIdwDrill)
		v1g.POST("/genDrill", apis.GenIdwDrill2)

		//v1g.POST(/isValidPoint)

		v1g.GET("/queryDrill", apis.QueryDrill)
		v1g.POST("/queryDrill", apis.QueryDrill2)
	}

	port := strings.Join([]string{":", viper.GetString("listen.port")}, "")
	err := r.Run(port)
	if err != nil {
		return
	}
}
