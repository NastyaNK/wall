package web

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
	"wall/internal/entity"
)

func (h *Handler) Authorization(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	password := r.URL.Query().Get("password")
	user, err := h.repository.GetUser(name)
	if err != nil {
		h.Error(w, "Нет такого пользователя", http.StatusBadRequest)
		return
	}
	if user.Password != Hash(password) {
		h.Error(w, "Неверный пароль", http.StatusBadRequest)
		return
	}
	bytes, err := json.Marshal(user)
	if err != nil {
		h.logger.Println(err)
		h.Error(w, "Проблема на сервере", http.StatusInternalServerError)
		return
	}

	expiration := time.Now().Add(24 * time.Hour)
	SetCookie(w, expiration, "name", user.Name, "password", user.Password)

	w.Write(bytes)
}

func (h *Handler) Registration(w http.ResponseWriter, r *http.Request) {
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Println(err)
		h.Error(w, "Ошибка при чтении запроса", http.StatusBadRequest)
		return
	}
	var user entity.User
	err = json.Unmarshal(bytes, &user)
	if err != nil {
		h.logger.Println(err)
		h.Error(w, "Проблема в теле", http.StatusBadRequest)
		return
	}
	re := regexp.MustCompile("^[а-яА-Яa-zA-Z0-9_]+$")
	if !re.MatchString(user.Name) {
		h.Error(w, "Имя неверное", http.StatusBadRequest)
		return
	}
	pass := regexp.MustCompile("^[а-яА-Яa-zA-Z0-9_]+$")
	if !pass.MatchString(user.Password) {
		h.Error(w, "Пароль неверный", http.StatusBadRequest)
		return
	}

	if user.Age < 12 || user.Age > 98 {
		h.Error(w, "Возвраст неверный", http.StatusBadRequest)
		return
	}

	user.Password = Hash(user.Password)
	_, err = h.repository.AddUser(&user)
	if err != nil {
		h.logger.Println(err)
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint \"users_pk\"") {
			h.Error(w, "Такой пользователь уже есть, придумайте другое имя", http.StatusBadRequest)
		} else {
			h.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	expiration := time.Now().Add(24 * time.Hour)
	SetCookie(w, expiration, "name", user.Name, "password", user.Password)

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	user, err := h.Login(r)
	if err != nil {
		h.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	bytes, err := json.Marshal(user)
	if err != nil {
		h.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(bytes)
}

func (h *Handler) Login(r *http.Request) (*entity.User, error) {
	name := GetCookie(r, "name")
	pass := GetCookie(r, "password")
	if len(name) == 0 || len(pass) == 0 {
		return nil, errors.New("нужно авторизироваться")
	}
	user, err := h.repository.GetUser(name)
	if err != nil {
		return nil, errors.New("нет такого пользователя")
	}
	if user.Password != pass {
		return nil, errors.New("неверный пароль")
	}
	return user, nil
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	SetCookie(w, time.Now().Add(-7*24*time.Hour), "name", "", "password", "")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
