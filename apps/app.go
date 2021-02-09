package apps

import (
	"flag"
	"fmt"

	"github.com/garbein/lottery-golang/configs"
	"github.com/gomodule/redigo/redis"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var App Application

type Application struct {
	ConfigFile string
	Config     configs.Config
	Logger     *zap.Logger
	DB         *gorm.DB
	Redis      redis.Conn
}

func InitApp() {
	App.initArgs()
	App.initConfig()
	App.initLogger()
	App.initDB()
	App.initRedis()
}

// 读取命令行参数
func (app *Application) initArgs() {
	var configFile string
	flag.StringVar(&configFile, "c", "configs/config.toml", "please config a *.toml file")
	flag.Parse()
	app.ConfigFile = configFile
}

// 读取配置
func (app *Application) initConfig() {
	v := viper.New()
	v.SetConfigFile(app.ConfigFile)
	fmt.Println("config file:", app.ConfigFile)
	err := v.ReadInConfig()
	if err != nil {
		panic("read config fail")
	}
	v.Unmarshal(&app.Config)
}

// 初始化zap
func (app *Application) initLogger() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic("init logger fail")
	}
	app.Logger = logger
}

func (app *Application) initDB() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s", app.Config.DB.Username, app.Config.DB.Password, app.Config.DB.Host, app.Config.DB.Port, app.Config.DB.Dbname, app.Config.DB.Charset)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic("failed to connect database")
	}
	app.DB = db
}

func (app *Application) initRedis() {
	client, err := redis.Dial("tcp", fmt.Sprintf("%s:%s", app.Config.Redis.Host, app.Config.Redis.Port))
	if err != nil {
		panic("failed to connect redis")
	}
	app.Redis = client
}
