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

type EsbServer struct {
	Host string
	Port int
}

func NewEsbServerConfig(v *viper.Viper) *EsbServer {
	return &EsbServer{
		Host: v.GetString("esb.host"),
		Port: v.GetInt("esb.port"),
	}
}

type ClientServer struct {
	Host string
	Port int
}

func NewClientServerConfig(v *viper.Viper) *ClientServer {
	return &ClientServer{
		Host: v.GetString("client.host"),
		Port: v.GetInt("client.port"),
	}
}

type AccServer struct {
	Host string
	Port int
}

func NewAccServerConfig(v *viper.Viper) *AccServer {
	return &AccServer{
		Host: v.GetString("acc.host"),
		Port: v.GetInt("acc.port"),
	}
}

type TransactionServer struct {
	Host string
	Port int
}

func NewTransactionServerConfig(v *viper.Viper) *TransactionServer {
	return &TransactionServer{
		Host: v.GetString("transaction.host"),
		Port: v.GetInt("transaction.port"),
	}
}

type PaymentServer struct {
	Host string
	Port int
}

func NewPaymentServerConfig(v *viper.Viper) *PaymentServer {
	return &PaymentServer{
		Host: v.GetString("payment.host"),
		Port: v.GetInt("payment.port"),
	}
}

type StatementServer struct {
	Host string
	Port int
}

func NewStatementServerConfig(v *viper.Viper) *StatementServer {
	return &StatementServer{
		Host: v.GetString("statement.host"),
		Port: v.GetInt("statement.port"),
	}
}

type Hash struct {
	Salt int
}

func NewHashConfig(v *viper.Viper) *Hash {
	return &Hash{
		Salt: v.GetInt("hash.salt"),
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
