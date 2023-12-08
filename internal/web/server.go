package web

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
	"wall/internal/entity"
	"wall/internal/repository"
)

type Server struct {
	router *chi.Mux
}

type Handler struct {
	repository repository.Repository
	logger     *log.Logger
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type ResponsePosts struct {
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Data    []entity.Post `json:"data"`
	Pages   int           `json:"pages"`
}

func NewServer(repo repository.Repository) *Server {
	h := Handler{repo, log.Default()}

	router := chi.NewRouter()

	router.Use(middleware.Logger)

	root := http.Dir("./front")
	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		routerCtx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(routerCtx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})

	router.Route("/api", func(router chi.Router) {

		router.Route("/posts", func(r chi.Router) {
			r.Get("/", h.GetPosts)
			//r.Get("/{id}", h.GetPost)
			r.Post("/", h.AddPost)
			r.Put("/", h.UpdatePosts)
			r.Delete("/{id}", h.DeletePost)
		})

		router.Route("/user", func(r chi.Router) {
			r.Get("/", h.GetUser)
			r.Post("/reg", h.Registration)
			r.Get("/auth", h.Authorization)
			r.Get("/logout", h.Logout)
		})

	})

	return &Server{router: router}
}

func (s *Server) Run() error {
	return http.ListenAndServe(":9898", s.router)
}

func (h *Handler) Error(w http.ResponseWriter, text string, status int) {
	w.WriteHeader(status)
	response := Response{false, text}
	bytes, _ := json.Marshal(response)
	w.Write(bytes)
}

func SetCookie(w http.ResponseWriter, time time.Time, cookie ...string) {
	for i := 0; i < len(cookie)-1; i += 2 {
		http.SetCookie(w, &http.Cookie{Path: "/", Name: cookie[i], Value: url.QueryEscape(cookie[i+1]), Expires: time})
	}
}

func GetCookie(r *http.Request, name string) string {
	c, err := r.Cookie(name)
	if err != nil {
		return ""
	}
	value, err := url.QueryUnescape(c.Value)
	if err != nil {
		return ""
	}
	return value
}

func Hash(text string) string {
	sum := sha256.Sum256([]byte(text))
	return fmt.Sprintf("%x", sum)
}
