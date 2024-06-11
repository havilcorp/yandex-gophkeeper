package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"ya-gophkeeper-client/internal/config"
	"ya-gophkeeper-client/internal/entity"

	"github.com/sirupsen/logrus"
)

type handler struct {
	conf   *config.Config
	client *http.Client
	token  string
}

func New(conf *config.Config, client *http.Client) *handler {
	return &handler{
		conf:   conf,
		client: client,
	}
}

func (h *handler) SetToken(token string) {
	h.token = token
}

func (h *handler) Save(dto *entity.ItemDto) error {
	data, err := json.Marshal(dto)
	if err != nil {
		logrus.Error(err)
		return err
	}
	body := bytes.NewReader(data)
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/storage", h.conf.AddressHttp), body)
	if err != nil {
		logrus.Error(err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", h.token))
	resp, err := h.client.Do(req)
	if err != nil {
		logrus.Error(err)
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New("status not 200")
	}
	return nil
}

func (h *handler) GetByID(id string) (*entity.ItemDto, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/storage/%s", h.conf.AddressHttp, id), nil)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", h.token))
	resp, err := h.client.Do(req)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("status not 200")
	}
	dec := json.NewDecoder(resp.Body)
	defer resp.Body.Close()
	date := entity.ItemDto{}
	if err := dec.Decode(&date); err != nil {
		logrus.Error(err)
		return nil, err
	}
	return &date, nil
}

func (h *handler) GetList() (*[]entity.ItemDto, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/storage", h.conf.AddressHttp), nil)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", h.token))
	resp, err := h.client.Do(req)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("status not 200")
	}
	dec := json.NewDecoder(resp.Body)
	defer resp.Body.Close()
	date := make([]entity.ItemDto, 0)
	if err := dec.Decode(&date); err != nil {
		logrus.Error(err)
		return nil, err
	}
	return &date, nil
}
