// Package http транспортный уровень авторизации
package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"yandex-gophkeeper-server/internal/auth"
	"yandex-gophkeeper-server/internal/auth/entity"
	"yandex-gophkeeper-server/internal/config"
	"yandex-gophkeeper-server/pkg/jwt"

	"github.com/sirupsen/logrus"
)

type handler struct {
	conf *config.Config
	uc   auth.UseCase
}

// NewHandler получить экземпляр хендлера
func NewHandler(conf *config.Config, uc auth.UseCase) *handler {
	return &handler{
		conf: conf,
		uc:   uc,
	}
}

func (h *handler) login(w http.ResponseWriter, r *http.Request) {
	var dto entity.Login
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&dto); err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	logrus.Println(dto)
	// uc
	user, err := h.uc.Login(dto.Email, dto.Password)
	if err != nil {
		if errors.Is(err, entity.ErrUserNotFound) {
			logrus.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		} else if errors.Is(err, entity.ErrWrongPassword) {
			logrus.Error(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else {
			logrus.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	token, err := jwt.GenerateJWT([]byte(h.conf.JWTKey), strconv.Itoa(user.ID))
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	b, err := json.Marshal(&entity.Response{Token: token})
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

func (h *handler) registration(w http.ResponseWriter, r *http.Request) {
	var dto entity.Registration
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&dto); err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user, err := h.uc.Registration(dto.Email, dto.Password)
	if err != nil {
		if errors.Is(err, entity.ErrUserExists) {
			logrus.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		} else {
			logrus.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	token, err := jwt.GenerateJWT([]byte(h.conf.JWTKey), strconv.Itoa(user.ID))
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	b, err := json.Marshal(&entity.Response{Token: token})
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(b)
}
