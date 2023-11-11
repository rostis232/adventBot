package main

import (
	"log"

	"github.com/rostis232/adventBot/internal/pkg/app"
	"github.com/spf13/viper"
)

func main(){
	if err := initConfig(); err != nil {
		log.Fatalln("error while config loading")
	  }
	a, err := app.NewApp(viper.GetString("db.dbname"))
	if err != nil {
		log.Fatal(err)
	}
	a.Run(viper.GetString("app.port"))
}

func initConfig() error {
	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}