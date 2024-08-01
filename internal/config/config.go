package config

import (
	"fmt"
	"os"

	viper "github.com/spf13/viper"
)

// Server is a struct that represents the server
type Server struct {
	Host string
	Port int
}

// NewServer is a function to create a new server
func NewServerConfig(v *viper.Viper) *Server {
	return &Server{
		Host: v.GetString("app.host"),
		Port: v.GetInt("app.port"),
	}
}

func LoadConfig(file string) (*viper.Viper, error) {
	appConfig := viper.New()

	configPath := "./internal/config"

	if _, err := os.Stat(configPath); err == nil {
		appConfig.AddConfigPath(configPath)

		appConfig.SetConfigName("pgsql")

		if commonErr := appConfig.ReadInConfig(); commonErr != nil {
			return nil, commonErr
		}

		if len(file) > 0 {
			appConfig.SetConfigName(file)
			fmt.Println("Loading config file:", file)
			if err := appConfig.MergeInConfig(); err != nil {
				return nil, err
			}
		}
	}
	appConfig.AutomaticEnv()

	return appConfig, nil

}
