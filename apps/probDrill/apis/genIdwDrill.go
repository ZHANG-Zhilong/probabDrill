package apis

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"probabDrill/internal/constant"
	"probabDrill/internal/entity"
	"probabDrill/internal/service"
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
func GenIdwDrill2(c *gin.Context) {
	x, _ := strconv.ParseFloat(c.PostForm("x"), 10)
	y, _ := strconv.ParseFloat(c.PostForm("y"), 10)
	hobbies := c.PostFormArray("hobby")
	drillSet := constant.GetHelpDrills()
	drill := service.GenVDrillIDW(drillSet, nil, nil, x, y)
	var drills []entity.Drill
	drills = append(drills, drill)
	d, _ := json.Marshal(drills)
	fmt.Println("hobbies:", hobbies)
	c.String(http.StatusOK, fmt.Sprintf("%s", string(d)))
}
func GenIdwDrill(c *gin.Context) {
	c.HTML(http.StatusOK, "genDrill.html", nil)
}
