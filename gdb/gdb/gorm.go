package gdb

import (
	"fmt"
	"github.com/dshibin/gbins/gconf"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
)

var (
	db          = &gorm.DB{}
	dsn         = ""
	serviceName = "mysql"
)

func init() {
	selectOrm()
	err := setOrm()
	if err != nil {
		panic(fmt.Sprintf("orm is wrong : ", err.Error))
	}
}

func NewOrm(c *gin.Context) *gorm.DB {
	_ = setOrm()
	return db.WithContext(c)
}

func selectOrm() {
	if gconf.GConfByClientName(serviceName).Target != "" {
		serviceName = gconf.GConfByClientName(serviceName).Target
	}
	return
}

func setOrm() error {
	if dsn != gconf.GetRVal(serviceName) {
		dsn = gconf.GetRVal(serviceName)
		d, err := gorm.Open(mysql.Open(dsn), debugOrm())
		if err == nil {
			pluginSet(d)
			defaultOrm(d)
			db = d
		}
	}
	return nil
}

func pluginSet(d *gorm.DB) {
	_ = d.Use(&TracePlugin{})
	return
}

func debugOrm() *gorm.Config {
	if strings.ToLower(gconf.GConfig().Global.Namespace) == gconf.NamespaceProd {
		return &gorm.Config{}
	}
	return &gorm.Config{}
}

func defaultOrm(do *gorm.DB) {
	d, err := do.DB()
	if err != nil {
		fmt.Print("the db is wrong")
	}
	if d == nil {
		return
	}
	d.SetMaxIdleConns(100)
	d.SetMaxOpenConns(400)
	return
}
