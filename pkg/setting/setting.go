package setting

import (
	"log"
	"time"

	"github.com/go-ini/ini" //用于解析 INI 配置文件
)

//编写与配置项保持一致的结构体（App、Server、Database）
type App struct {
	JwtSecret string
	PageSize int
	RuntimeRootPath string

	ImagePrefixUrl string
	ImageSavePath string
	ImageMaxSize int
	ImageAllowExts []string

	LogSavePath string
	LogSaveName string
	LogFileExt string
	TimeFormat string
}

var AppSetting = &App{}

type Server struct {
	RunMode string
	HttpPort int
	ReadTimeout time.Duration
	WriteTimeout time.Duration
}

var ServerSetting = &Server{}


type Database struct {
	Type string
	User string
	Password string
	Host string
	Name string
	TablePrefix string
}

var DatabaseSetting = &Database{}


type Redis struct {
	Host        string
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}

var RedisSetting = &Redis{}

var cfg *ini.File




func Setup() {
	var err error
	cfg, err = ini.Load("conf/app.ini")

	if err != nil {
		log.Fatalf("Fail to parse 'conf/app.ini': %v", err)
	}

	// //使用 MapTo 将配置项映射到结构体上
	//err = Cfg.Section("app").MapTo(AppSetting)
	// if err != nil {
	// 	log.Fatalf("Cfg.MapTo AppSetting err: %v", err)
	// }

	mapTo("app", AppSetting)
	mapTo("server", ServerSetting)
	mapTo("database", DatabaseSetting)
	mapTo("redis", RedisSetting)

	//将单位转换为MB
	AppSetting.ImageMaxSize = AppSetting.ImageMaxSize * 1024 * 1024

	// err = Cfg.Section("server").MapTo(ServerSetting)
	// if err != nil {
	// 	log.Fatalf("Cfg.MapTo ServerSetting err: %v", err)
	// }

	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.WriteTimeout * time.Second
	RedisSetting.IdleTimeout = RedisSetting.IdleTimeout * time.Second

	// err = Cfg.Section("database").MapTo(DatabaseSetting)
	// if err != nil {
	// 	log.Fatalf("Cfg.MapTo DatabaseSetting err: %v", err)
	// }

}

func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %sseting err: %v", section, err)
	}
}


