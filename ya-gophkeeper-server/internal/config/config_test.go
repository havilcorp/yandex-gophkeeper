package config

import (
	"errors"
	"os"
	"testing"
)

func Test_writeConfigByFile(t *testing.T) {
	t.Run("good", func(t *testing.T) {
		conf := Config{}
		file, err := os.Create("./test.json")
		if err != nil {
			t.Error(err)
		}
		file.WriteString("{}")
		file.Close()
		err = writeConfigByFile("./test.json", &conf)
		if err != nil {
			t.Error(err)
		}
		os.Remove("test.json")
	})
	t.Run("bad json", func(t *testing.T) {
		conf := Config{}
		file, err := os.Create("./test.json")
		if err != nil {
			t.Error(err)
		}
		file.WriteString("")
		file.Close()
		err = writeConfigByFile("./test.json", &conf)
		if err == nil {
			t.Error(errors.New("err is nil"))
		}
		os.Remove("test.json")
	})
}

func TestNew(t *testing.T) {
	os.Args = append(os.Args, "-address_http=:8080")
	os.Args = append(os.Args, "-address_grpc=:8081")
	os.Args = append(os.Args, `-db=postgress`)

	t.Setenv("ADDRESS_HTTP", ":8080")
	t.Setenv("ADDRESS_GRPC", ":8081")
	t.Setenv("CA_CRT", "./tls/ca.crt")
	t.Setenv("SERVER_CRT", "./tls/server.crt")
	t.Setenv("SERVER_KEY", "./tls/server.key")
	t.Setenv("DB_CONNECT", "postgres:/...")
	t.Setenv("JWT_KEY", "jwt")
	New()
}
