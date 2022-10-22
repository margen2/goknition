package main

import (
	"fmt"
	"log"

	"github.com/margen2/goknition/backend/db"
	"github.com/spf13/viper"
)

var (
	StringConnectionBD = ""
)

// Load loads the config.json configuration file
func load() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	viper.SetDefault("DB_USER", "root")
	viper.SetDefault("DB_PASSWORD", "root")
	viper.SetDefault("DB_NAME", "goknition")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			err = viper.SafeWriteConfig()
			if err != nil {
				fmt.Println("safewrite")
				log.Fatal(err)
			}
		} else {
			fmt.Println("else")
			log.Fatal(err)
		}
	}

	StringConnectionBD = fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local",
		viper.GetString("DB_USER"),
		viper.GetString("DB_PASSWORD"),
		viper.GetString("DB_NAME"),
	)

	db.SetConnection(StringConnectionBD)
}
