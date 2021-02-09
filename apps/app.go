package apps

import (
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
	Config configs.Config
	Logger *zap.Logger
	DB     *gorm.DB
	Redis  redis.Conn
}

func InitApp() {
	App.initConfig()
	App.initLogger()
	App.initDB()
	App.initRedis()
}

func (app *Application) initConfig() {
	v := viper.New()
	v.SetConfigFile("configs/config.toml")
	err := v.ReadInConfig()
	if err != nil {
		panic("read config fail")
	}
	v.Unmarshal(&app.Config)
}

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
