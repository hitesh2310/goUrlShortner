package config

import (
	"fmt"
	"main/logs"
	"main/pkg/constants"
	"main/pkg/database"
	models "main/pkg/models/configModel"

	"github.com/spf13/viper"
)

func SetUpApplication() {
	fmt.Println("SetUp config details....")
	setupConfig()
	fmt.Println("SettingUp Logs....")
	// setUpApplicationLogs()
	setUpApplicationDatabase()
	fmt.Println("DB set up done")
	setUpRedis()
	fmt.Println("Redis set up done")
}

func setupConfig() {
	viper.SetConfigName("config")    // Name of the config file (without extension)
	viper.SetConfigType("json")      // Config file type
	viper.AddConfigPath("./config/") // Path to the directory containing the config file

	// Read the config file
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("failed to read config file: %s", err))
	}

	var config models.Config
	err = viper.Unmarshal(&config)
	if err != nil {
		fmt.Println("Error to unmarshal config")
	}

	constants.ApplicationConfig = &config

}

func setUpApplicationLogs() {
	logs.SetUpApplicationLogs()
}

func setUpApplicationDatabase() {
	database.EstablishDbConnection()
}

func setUpRedis() {
	database.SetUpRedis()
}
