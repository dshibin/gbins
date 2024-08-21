package gconf

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/spf13/cast"
	"io/ioutil"
	"sync/atomic"
)

const defaultPath = "conf.json"

var (
	ConfigPath = defaultPath
	gm         = atomic.Value{}
)

func init() {
	SetGConfig(LoadConfig(fmt.Sprintf("./conf/%s", GetConfigPath())))
}

func GetConfigPath() string {
	if ConfigPath == defaultPath {
		flag.StringVar(&ConfigPath, "conf", defaultPath, "load config path")
		flag.Parse()
	}
	return ConfigPath
}

func LoadConfig(path string) *GlobalConfig {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		panic("the conf file is wrong")
	}
	globalConfig := &GlobalConfig{}
	err = json.Unmarshal(buf, globalConfig)
	if err != nil {
		panic(fmt.Sprintf("decode config file failed : %s err :%s", string(buf), err.Error()))
	}
	return globalConfig
}

func GConfig() *GlobalConfig {
	return gm.Load().(*GlobalConfig)
}

func GAddr() string {
	ip := GConfig().Server.Ip
	if ip == "" {
		ip = "0.0.0.0"
	}
	port := GConfig().Server.Port
	if port == "" {
		port = "8080"
	}
	return fmt.Sprintf("%s:%s", ip, port)
}

func GConfByClientName(name string) (cli Client) {
	gclientArr := GConfig()
	for _, c := range gclientArr.Client {
		if c.Name == name {
			return c
		}
	}
	return
}

func GetConfigByClient(target, name string) (ret string) {
	cli := GConfByClientName(target)
	if len(cli.Config) == 0 {
		return
	}
	if r, ok := cli.Config[name]; ok {
		ret = cast.ToString(r)
	}
	return
}

func SetGConfig(cfg *GlobalConfig) {
	defaultInit(cfg)
	gm.Store(cfg)
}

// 默认配置
func defaultInit(cfg *GlobalConfig) {
	if cfg.Server.Limit <= 0 {
		cfg.Server.Limit = 1000
	}
}
