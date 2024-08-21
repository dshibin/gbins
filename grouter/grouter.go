package grouter

import (
	"github.com/dshibin/gbins/gconf"
	"github.com/dshibin/gbins/grouter/gmiddleware"
	"github.com/gin-gonic/gin"
	"strings"
)

var Router = gin.New()

func init() {
	if strings.ToLower(gconf.GConfig().Global.Namespace) == gconf.NamespaceProd {
		gin.SetMode(gin.ReleaseMode)
	}
	Router.Use(gin.Recovery())
	Router.Use(gmiddleware.GTotalTime())
	Router.Use(gmiddleware.GSeq())
	Router.Use(gmiddleware.GStartLog())
	Router.Use(gmiddleware.GTimeout())
}
