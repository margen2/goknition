package main

import (
	"log"

	"github.com/margen2/goknition/backend/db"
	"github.com/spf13/viper"
)

// Load loads the config.json configuration file
func loadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	viper.SetDefault("DB_USER", "go")
	viper.SetDefault("DB_PASSWORD", "golang")
	viper.SetDefault("DB_NAME", "goknition_new")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			err = viper.SafeWriteConfig()
			if err != nil {
				log.Fatal(err)
			}
			log.Fatal("resetting application")
		} else {
			log.Fatal(err)
		}
	}

	err := db.SetConnection(
		viper.GetString("DB_USER"),
		viper.GetString("DB_PASSWORD"),
		viper.GetString("DB_NAME"))
	if err != nil {
		log.Fatal(err)
	}
}
