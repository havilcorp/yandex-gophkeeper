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
}

func New() *Config {
	conf := Config{
		AddressHttp: "localhost:8080",
	}
	var AddressHttp, AddressGRPC string

	if data, err := os.ReadFile("./config.json"); err == nil {
		if err := json.Unmarshal(data, &conf); err != nil {
			logrus.Error(err)
		}
	}

	flag.StringVar(&AddressHttp, "address_http", "", "address and port to run server")
	flag.StringVar(&AddressGRPC, "address_grpc", "", "address and port to run grpc server")
	flag.Parse()

	if AddressHttp != "" {
		conf.AddressHttp = AddressHttp
	}

	if AddressGRPC != "" {
		conf.AddressGRPC = AddressGRPC
	}

	if envAddressHttp := os.Getenv("ADDRESS_HTTP"); envAddressHttp != "" {
		conf.AddressHttp = envAddressHttp
	}

	if envAddressGRPC := os.Getenv("ADDRESS_GRPC"); envAddressGRPC != "" {
		conf.AddressGRPC = envAddressGRPC
	}

	return &conf
}
