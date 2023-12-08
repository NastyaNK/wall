package repository

import (
	"encoding/json"
	"os"
	"wall/internal/entity"
)

type Repository interface {
	Connect(config *entity.DBConfig) error
	Close() error

	AddPost(post *entity.Post) (int, error)
	UpdatePost(post *entity.Post) (int64, error)
	DeletePost(user *entity.User, id int) (int64, error)
	GetPost(user *entity.User, id int) (*entity.Post, error)
	GetPosts(user *entity.User) ([]entity.Post, error)
	GetPagePosts(user *entity.User, page int) ([]entity.Post, error)
	GetPagesCount(user *entity.User) (int, error)

	AddUser(user *entity.User) (int, error)
	GetUser(name string) (*entity.User, error)
}

func LoadConfig(filepath string) (*entity.DBConfig, error) {
	content, err := os.ReadFile(filepath)
	var config entity.DBConfig
	err = json.Unmarshal(content, &config)
	return &config, err
}
