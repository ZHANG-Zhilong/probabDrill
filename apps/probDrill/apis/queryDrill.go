package apis

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"probabDrill/internal/constant"
)

func QueryDrill(c *gin.Context) {
	c.HTML(http.StatusOK, "queryDrill.html", nil)
}
func QueryDrill2(c *gin.Context) {
	drillName := c.PostForm("name")
	if d, ok := constant.GetRealDrillByName(drillName); ok {
		djson, _ := json.Marshal(d)
		c.String(http.StatusOK, fmt.Sprintf(" drill:%s\n", string(djson)))
	} else {
		c.String(http.StatusBadRequest, fmt.Sprintf("invalid drill name: %s.\n", drillName))
	}
}
