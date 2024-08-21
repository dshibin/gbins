package gconf

const (
	NamespaceProd = "prod" // 正式环境namespace配置
)

type GlobalConfig struct {
	Global struct {
		Namespace string `json:"namespace"`
		EnvName   string `json:"env_name"`
	} `json:"global"`
	Server struct {
		App     string `json:"app"`
		Ip      string `json:"ip"`
		Port    string `json:"port"`
		Limit   int    `json:"limit"`
		Timeout int    `json:"timeout"`
	} `json:"server"`
	Client []Client `json:"client"`
}

type Client struct {
	Name    string            `json:"name"`
	Target  string            `json:"target"`
	Network string            `json:"network"`
	Timeout int64             `json:"timeout"`
	Config  map[string]string `json:"config"`
}

type Grconf struct {
	Id   int32  `json:"id"`
	Env  string `json:"env_name"`
	Name string `json:"name"`
	App  string `json:"app"`
	Typ  int32  `json:"typ"`
	Val  string `json:"val"`
}
