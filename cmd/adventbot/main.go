package main

import (
	"fmt"
	"log"

	"github.com/rostis232/adventBot/config"
	"github.com/rostis232/adventBot/internal/pkg/app"
	"github.com/spf13/viper"
)

func main(){
	if err := initConfig(); err != nil {
		log.Fatalln("error while config loading")
	  }
	
	  config := &config.Config{
		DBname: viper.GetString("db.dbname"),
		Port: viper.GetString("app.port"),
		TGsecretCode: viper.GetString("app.tg_secretkey"),
		AdminLogin: viper.GetString("app.admin_login"),
		AdminPass: viper.GetString("app.admin_pass"),
		AppName: viper.GetString("app.app_name"),
		AppDesc: viper.GetString("app.app_desc"),
		TGlink: viper.GetString("app.tglink"),
		InstaLink: viper.GetString("app.instalink"),
	  }
	fmt.Println("!", config.Port)
	a, err := app.NewApp(config)
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