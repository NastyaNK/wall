package tests

import (
	"encoding/json"
	"os"
	"testing"
	"wall/internal/entity"
	"wall/internal/repository"
	"wall/internal/repository/postgres"
)

func GetRepository(t *testing.T) repository.Repository {
	content, err := os.ReadFile("../configs/db.json")
	var config entity.DBConfig
	db := postgres.NewDatabase()
	err = json.Unmarshal(content, &config)
	if err != nil {
		t.Fatal("json.Unmarshal:", err)
	}
	err = db.Connect(&config)
	if err != nil {
		t.Fatal("проблемы с подключением базы", err)
	}
	return db
}

func TestAddUser(t *testing.T) {
	db := GetRepository(t)
	defer db.Close()
	user := entity.User{
		Name:     "Mal2",
		Age:      23,
		Password: "qwerty",
	}
	_, err := db.AddUser(&user)
	if err != nil {
		t.Fatal("добавление пользователя", err)
	}
}
func TestGetUser(t *testing.T) {
	db := GetRepository(t)
	defer db.Close()
	_, err := db.GetUser("Mal", "jgfd8")
	if err != nil {
		t.Fatal("вытаскивание пользователя", err)
	}
}
func TestAddPost(t *testing.T) {
	db := GetRepository(t)
	defer db.Close()
	user := entity.User{Id: 5}
	post := entity.NewPost("xyk", &user)
	_, err := db.AddPost(post)
	if err != nil {
		t.Fatal("добаваление поста", err)
	}
}
func TestGetPosts(t *testing.T) {
	db := GetRepository(t)
	defer db.Close()
	user := entity.User{Id: 5}
	_, err := db.GetPosts(&user)
	if err != nil {
		t.Fatal("вытаскивание поста 1", err)
	}
}
func TestUpdatePost(t *testing.T) {
	db := GetRepository(t)
	defer db.Close()
	post := entity.Post{Id: 5}
	post.Text = "yuuuil"
	err := db.UpdatePost(&post)
	if err != nil {
		t.Fatal("обновление поста", err)
	}
}

func TestDeletePost(t *testing.T) {
	db := GetRepository(t)
	defer db.Close()
	err := db.DeletePost(9)
	if err != nil {
		t.Fatal("delete поста", err)
	}
}
