package gconf

import (
	"encoding/json"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"sync"
	"time"
)

var grDB = &gorm.DB{}
var grMap = map[string]*Grconf{}
var mu sync.RWMutex

const (
	tb    = "gconf"
	Gconf = "gconf"
)

func init() {
	dsn := GConfByClientName(Gconf).Target
	d, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("grconf is err")
	}
	grDB = d
	initRconf()
}

func initRconf() {
	getRConf()
	go func() {
		for true {
			time.Sleep(time.Minute)
			getRConf()
		}
	}()
}

func getRConf() {
	grconfArr := []*Grconf{}
	grDB.Table(tb).Where("app = ?", GConfig().Server.App).Where("env = ?", GConfig().Global.EnvName).Find(&grconfArr)
	mu.Lock()
	for _, v := range grconfArr {
		grMap[fmt.Sprintf("%s_%s_%s", v.Name, v.App, v.Env)] = v
	}
	mu.Unlock()
}

func getKey(name string) string {
	return fmt.Sprintf("%s_%s_%s", name, GConfig().Server.App, GConfig().Global.EnvName)
}

func GetGrconf(name string) *Grconf {
	mu.RLock()
	defer mu.RUnlock()
	if gc, ok := grMap[getKey(name)]; ok {
		return gc
	}
	return nil
}

func GetRVal(name string) string {
	g := GetGrconf(name)
	if g != nil {
		return g.Val
	}
	return ""
}

func GetVFromValMap(name, k string) string {
	vmap := GetRValMap(name)
	if v, ok := vmap[k]; ok {
		return v
	}
	return ""
}

func GetRValMap(name string) (vmap map[string]string) {
	g := GetGrconf(name)
	if g != nil {
		err := json.Unmarshal([]byte(g.Val), &vmap)
		if err != nil {
			log.Fatalln("the conf val not a json")
		}
		return
	}
	return
}
