package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"ya-gophkeeper-client/internal/config"
	"ya-gophkeeper-client/internal/store/entity"

	"github.com/sirupsen/logrus"
)

type handler struct {
	conf *config.Config
}

func New(conf *config.Config) *handler {
	return &handler{
		conf: conf,
	}
}

func (h *handler) Save(dto *entity.SaveDto) error {
	data, err := json.Marshal(dto)
	if err != nil {
		logrus.Error(err)
		return err
	}
	body := bytes.NewReader(data)
	contentType := "application/json"
	resp, err := http.Post(fmt.Sprintf("%s/storage", h.conf.AddressHttp), contentType, body)
	if err != nil {
		logrus.Error(err)
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New("status not 200")
	}
	return nil
}

func (h *handler) GetByID(id string) (*entity.SaveDto, error) {
	resp, err := http.Get(fmt.Sprintf("%s/storage/%s", h.conf.AddressHttp, id))
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("status not 200")
	}
	dec := json.NewDecoder(resp.Body)
	defer resp.Body.Close()
	date := entity.SaveDto{}
	if err := dec.Decode(&date); err != nil {
		logrus.Error(err)
		return nil, err
	}
	return &date, nil
}
