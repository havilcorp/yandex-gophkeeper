package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"ya-gophkeeper-server/internal/storage"
	"ya-gophkeeper-server/internal/storage/entity"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type handler struct {
	uc storage.UserCase
}

func NewHandler(uc storage.UserCase) *handler {
	return &handler{
		uc: uc,
	}
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	dto := entity.CreateDto{}
	dec := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := dec.Decode(&dto); err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err := h.uc.Save(1, &dto)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *handler) GetOne(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	item, err := h.uc.GetById(id)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	data, err := json.Marshal(item)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *handler) GetAll(w http.ResponseWriter, r *http.Request) {
	items, err := h.uc.GetAll()
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
	_, err = w.Write(data)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
}

func (h *handler) Remove(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = h.uc.Remove(id)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
