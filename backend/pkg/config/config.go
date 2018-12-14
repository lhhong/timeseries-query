package config

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
)

type ConfigModel struct {
	Database   DatabaseConfig
	Redis      RedisConfig
	HttpServer HttpConfig
}

type DatabaseConfig struct {
	Hostname string
	Port     int
	Username string
	Password string
	Database string
}

type RedisConfig struct {
	Hostname string
	Port     int
}

type HttpConfig struct {
	Port int
}

var Config *ConfigModel

func LoadConfig(cmd *cobra.Command) (ConfigModel, error) {

	var c ConfigModel

	err := viper.BindPFlags(cmd.Flags())
	if err != nil {
		return c, err
	}

	viper.SetEnvPrefix("TSQ")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Set viper path and read configuration
	if configFile, _ := cmd.Flags().GetString("config"); configFile != "" {
		viper.SetConfigFile(configFile)
		err := viper.ReadInConfig()

		// Handle errors reading the config file
		if err != nil {
			log.Fatalln("Fatal error config file", err)
			return c, err
		}
	} else {

		viper.AddConfigPath("conf")
		viper.SetConfigName("default")
		err := viper.ReadInConfig()

		// Handle errors reading the config file
		if err != nil {
			log.Fatalln("Fatal error config file", err)
			return c, err
		}

		if os.Getenv("ENV") == "prod" {
			viper.SetConfigName("production")
			err := viper.MergeInConfig()
			if err != nil {
				log.Fatalln("Fatal error config file", err)
				return c, err
			}
		}
	}

	if err := viper.Unmarshal(&c); err != nil {
		log.Fatalln("couldn't read config", err)
	}

	Config = &c

	return c, nil

}
