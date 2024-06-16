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
	uc   auth.UserCase
}

func NewHandler(conf *config.Config, uc auth.UserCase) *handler {
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
	_, err = w.Write(b)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
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
	_, err = w.Write(b)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
