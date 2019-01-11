package config

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Config Main Configuration Model
type Config struct {
	Database   DatabaseConfig
	Redis      RedisConfig
	HTTPServer HTTPConfig
}

// DatabaseConfig Database Config Model
type DatabaseConfig struct {
	Hostname string
	Port     int
	Username string
	Password string
	Database string
}

// RedisConfig Redis Config Model
type RedisConfig struct {
	Env      string
	Hostname string
	Port     int
}

// HTTPConfig Http Config Model
type HTTPConfig struct {
	Port int
}

// GetConfig Returns Config given cobra.Command which can contain config file
// panic on failure
func GetConfig(cmd *cobra.Command) *Config {

	var c Config

	//err := viper.BindPFlags(cmd.Flags())
	//if err != nil {
	//	panic(err)
	//}

	viper.SetEnvPrefix("TSQ")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Set viper path and read configuration
	if configFile, _ := cmd.Flags().GetString("config"); configFile != "" {
		readConfigFromFile(configFile)
	} else {
		readConfigFromDefault()
	}

	if err := viper.Unmarshal(&c); err != nil {
		log.Fatalln("couldn't read config", err)
		panic(err)
	}

	return &c

}

func readConfigFromDefault() {

	viper.AddConfigPath("conf")
	viper.SetConfigName("default")
	err := viper.ReadInConfig()

	// Handle errors reading the config file
	if err != nil {
		log.Fatalln("Fatal error config file", err)
		panic(err)
	}

	if os.Getenv("ENV") == "prod" {
		viper.SetConfigName("production")
		err := viper.MergeInConfig()
		if err != nil {
			log.Fatalln("Fatal error config file", err)
			panic(err)
		}
	}
}

func readConfigFromFile(fileName string) {
	viper.SetConfigFile(fileName)
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatalln("Fatal error reading from provided file", err)
		panic(err)
	}
}
