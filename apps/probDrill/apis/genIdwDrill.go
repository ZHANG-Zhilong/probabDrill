package apis

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"probabDrill/apps/probDrill/model"
	"probabDrill/internal/constant"
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

// GenIdwDrill2 godoc
// @Id GenIdwDrill2
// @Summary Generation drill by idw method.
// @Description get accounts
// @Tags gen
// @Accept  multipart/form-data
// @Param x body float32 true "position"
// @Param y body float32 true "position"
// @Produce  json
// @Success 200 array []model.Drill
// @Router /genDrill [post]
func GenIdwDrill2(c *gin.Context) {
	log.Println(*c)
	x, _ := strconv.ParseFloat(c.PostForm("x"), 10)
	y, _ := strconv.ParseFloat(c.PostForm("y"), 10)
	hobbies := c.PostFormArray("hobby")
	drillSet := constant.GetHelpDrills()
	drill := service.GenVDrillIDW(drillSet, nil, nil, x, y)
	var drills []model.Drill
	drills = append(drills, drill)
	d, _ := json.Marshal(drills)
	fmt.Println("hobbies:", hobbies)
	c.String(http.StatusOK, fmt.Sprintf("%s", string(d)))
}

func GenIdwDrill(c *gin.Context) {
	c.HTML(http.StatusOK, "genDrill.html", nil)
}
