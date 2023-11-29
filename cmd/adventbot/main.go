package main

import (
	"log"

	"github.com/rostis232/adventBot/config"
	"github.com/rostis232/adventBot/internal/pkg/app"
	"github.com/spf13/viper"
)

func main(){
	if err := initConfig(); err != nil {
		log.Fatalln("error while cfg loading")
	  }
	
	  cfg := &config.Config{
		DBname: viper.GetString("db.dbname"),
		RedisAddress: viper.GetString("db.redis"),
		Port: viper.GetString("app.port"),
		TGsecretCode: viper.GetString("app.tg_secretkey"),
		AdminLogin: viper.GetString("app.admin_login"),
		AdminPass: viper.GetString("app.admin_pass"),
		AppName: viper.GetString("app.app_name"),
		AppDesc: viper.GetString("app.app_desc"),
		TGlink: viper.GetString("app.tglink"),
		InstaLink: viper.GetString("app.instalink"),
		MyTime: viper.GetString("app.time"),
	  }
	a, err := app.NewApp(cfg)
	if err != nil {
		log.Fatal(err)
	}
	a.Run()
}

func initConfig() error {
	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}