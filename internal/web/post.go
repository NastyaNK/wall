package web

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"strconv"
	"wall/internal/entity"
)

func (h *Handler) GetPosts(w http.ResponseWriter, r *http.Request) {
	user, err := h.Login(r)
	if err != nil {
		h.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	posts, err := h.repository.GetPagePosts(user, page)

	if err != nil {
		h.Error(w, "А надо"+err.Error(), http.StatusInternalServerError)
		return
	}
	t, err := h.repository.GetPagesCount(user)
	if err != nil {
		h.Error(w, "А надо"+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	response := ResponsePosts{true, "", posts, t}
	bytes, _ := json.Marshal(response)
	w.Write(bytes)
}

func (h *Handler) AddPost(w http.ResponseWriter, r *http.Request) {
	user, err := h.Login(r)
	if err != nil {
		h.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var post entity.Post
	bytes, err := io.ReadAll(r.Body)
	err = json.Unmarshal(bytes, &post)
	if err != nil {
		h.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	post = *entity.NewPost(post.Text, user)
	id, err := h.repository.AddPost(&post)
	if err != nil {
		h.Error(w, "А надо"+err.Error(), http.StatusInternalServerError)
		return
	}
	post.Id = id

	w.WriteHeader(http.StatusOK)
	response := ResponsePosts{true, "", []entity.Post{post}, 1}
	bytes, _ = json.Marshal(response)
	w.Write(bytes)
}

func (h *Handler) UpdatePosts(w http.ResponseWriter, r *http.Request) {
	user, err := h.Login(r)
	if err != nil {
		h.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var post entity.Post
	bytes, err := io.ReadAll(r.Body)
	err = json.Unmarshal(bytes, &post)
	if err != nil {
		h.Error(w, "Неверные параметры запроса", http.StatusBadRequest)
		return
	}
	post.UserId = user.Id

	affected, err := h.repository.UpdatePost(&post)
	if affected == 0 || err != nil {
		h.Error(w, "Укажите id поста к которому вы имеете доступ", http.StatusBadRequest)
		return
	}

	postUpdated, err := h.repository.GetPost(user, post.Id)
	if err != nil {
		h.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	response := ResponsePosts{true, "", []entity.Post{*postUpdated}, 1}
	bytes, _ = json.Marshal(response)
	w.Write(bytes)
}

func (h *Handler) DeletePost(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "id")
	id, err := strconv.Atoi(param)
	if err != nil {
		h.Error(w, "id должен быть числом", http.StatusBadRequest)
		return
	}

	user, err := h.Login(r)
	if err != nil {
		h.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	affected, err := h.repository.DeletePost(user, id)
	if affected == 0 || err != nil {
		h.Error(w, "Укажите id поста к которому вы имеете доступ", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
}
