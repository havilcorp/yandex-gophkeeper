package config

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/sirupsen/logrus"
)

type Config struct {
	AddressHttp string `json:"address_http"`
	AddressGRPC string `json:"address_grpc"`
	CACrt       string `json:"ca_crt"`
	ServerCrt   string `json:"server_crt"`
	ServerKey   string `json:"server_key"`
	DBConnect   string `json:"db_connect"`
	JWTKey      string `json:"jwt_key"`
}

func writeConfigByFile(file string, conf *Config) error {
	if data, err := os.ReadFile(file); err == nil {
		if err = json.Unmarshal(data, &conf); err != nil {
			logrus.Error(err)
			return err
		}
	}
	return nil
}

func New() *Config {
	conf := Config{
		AddressHttp: "localhost:8080",
	}
	var AddressHttp, AddressGRPC, DBConnect, CACrt, ServerCrt, ServerKey, JWTKey string

	writeConfigByFile("./config.json", &conf)

	flag.StringVar(&AddressHttp, "address_http", "", "address and port to run server")
	flag.StringVar(&AddressGRPC, "address_grpc", "", "address and port to run grpc server")
	flag.StringVar(&CACrt, "ca_crt", "", "")
	flag.StringVar(&ServerCrt, "server_crt", "", "")
	flag.StringVar(&ServerKey, "server_key", "", "")
	flag.StringVar(&DBConnect, "db", "", "address database connect")
	flag.StringVar(&JWTKey, "jwt_key", "", "jwt key")
	flag.Parse()

	if AddressHttp != "" {
		conf.AddressHttp = AddressHttp
	}

	if AddressGRPC != "" {
		conf.AddressGRPC = AddressGRPC
	}

	if DBConnect != "" {
		conf.DBConnect = DBConnect
	}

	if envAddressHttp := os.Getenv("ADDRESS_HTTP"); envAddressHttp != "" {
		conf.AddressHttp = envAddressHttp
	}

	if envAddressGRPC := os.Getenv("ADDRESS_GRPC"); envAddressGRPC != "" {
		conf.AddressGRPC = envAddressGRPC
	}

	if envCaCrt := os.Getenv("CA_CRT"); envCaCrt != "" {
		conf.CACrt = envCaCrt
	}

	if envServerCrt := os.Getenv("SERVER_CRT"); envServerCrt != "" {
		conf.ServerCrt = envServerCrt
	}

	if envServerKey := os.Getenv("SERVER_KEY"); envServerKey != "" {
		conf.ServerKey = envServerKey
	}

	if envDBConnect := os.Getenv("DB_CONNECT"); envDBConnect != "" {
		conf.DBConnect = envDBConnect
	}

	if envJWTKey := os.Getenv("JWT_KEY"); envJWTKey != "" {
		conf.JWTKey = envJWTKey
	}

	return &conf
}
