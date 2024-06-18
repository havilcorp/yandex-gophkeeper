// Package http транспортный уровень хранилища
package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"yandex-gophkeeper-server/internal/config"
	"yandex-gophkeeper-server/internal/storage"
	"yandex-gophkeeper-server/internal/storage/entity"

	"github.com/sirupsen/logrus"
)

type handler struct {
	conf *config.Config
	uc   storage.UseCase
}

// NewHandler получить экземпляр хендлера
func NewHandler(conf *config.Config, uc storage.UseCase) *handler {
	return &handler{
		conf: conf,
		uc:   uc,
	}
}

func (h *handler) save(w http.ResponseWriter, r *http.Request) {
	userIDString := r.Header.Get("X-User-ID")
	if userIDString == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userID, err := strconv.Atoi(userIDString)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	dto := entity.CreateDto{}
	dec := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := dec.Decode(&dto); err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = h.uc.Save(userID, &dto)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *handler) getAll(w http.ResponseWriter, r *http.Request) {
	userIDString := r.Header.Get("X-User-ID")
	if userIDString == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userID, err := strconv.Atoi(userIDString)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	items, err := h.uc.GetAll(userID)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(items)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
