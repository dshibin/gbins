package gmiddleware

import (
	"github.com/dshibin/gbins/glog"
	"github.com/dshibin/gbins/greq"
	"github.com/gin-gonic/gin"
	"strings"
)

func GStartLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		postdata := ""
		if !strings.Contains(c.GetHeader("Content-Type"), "multipart/form-data") {
			postdata = greq.GetAllBodyJson(c)
		}

		path := ""
		if c.Request != nil && c.Request.URL != nil && c.Request.URL.Path != "" {
			path = c.Request.URL.Path
		}
		glog.Infof(c, "router : %s , req : %s", path, postdata)
	}
}
