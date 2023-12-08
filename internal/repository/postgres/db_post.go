package postgres

import (
	"math"
	"time"
	"wall/internal/entity"
)

var count = 7

func (psql *Database) GetPosts(user *entity.User) ([]entity.Post, error) {
	var posts []entity.Post
	err := psql.db.Select(&posts, "SELECT * FROM posts Where user_id=$1", user.Id)
	return posts, err
}

func (psql *Database) GetPost(user *entity.User, id int) (*entity.Post, error) {
	var post entity.Post
	err := psql.db.Get(&post, "SELECT * FROM posts Where id=$1 AND user_id=$2", id, user.Id)
	return &post, err
}

func (psql *Database) AddPost(post *entity.Post) (int, error) {
	rows, err := psql.db.NamedQuery("insert into posts(text, user_id, created, updated) VALUES (:text,:user_id,:created,:updated) RETURNING id", post)
	if err != nil {
		return 0, err
	}
	rows.Next()
	err = rows.Scan(&post.Id)
	return post.Id, err
}

func (psql *Database) UpdatePost(post *entity.Post) (int64, error) {
	post.Updated = time.Now()
	result, err := psql.db.NamedExec("UPDATE posts set text=:text, updated=:updated where id=:id and user_id=:user_id", post)
	if err != nil {
		return 0, err
	}
	affected, err := result.RowsAffected()
	return affected, err
}

func (psql *Database) DeletePost(user *entity.User, id int) (int64, error) {
	result, err := psql.db.Exec("DELETE from posts where id=$1 and user_id=$2", id, user.Id)
	if err != nil {
		return 0, err
	}
	affected, err := result.RowsAffected()
	return affected, err
}

func (psql *Database) GetPagePosts(user *entity.User, page int) ([]entity.Post, error) {
	var post []entity.Post
	err := psql.db.Select(&post, "SELECT * FROM posts  Where user_id=$1 ORDER BY created DESC OFFSET $2 LIMIT $3", user.Id, (page-1)*count, count)
	return post, err
}

func (psql *Database) GetPagesCount(user *entity.User) (int, error) {
	var countPosts int
	err := psql.db.Get(&countPosts, "SELECT count(*) FROM posts Where user_id=$1", user.Id)
	return int(math.Ceil(float64(countPosts) / float64(count))), err
}
