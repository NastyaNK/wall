package entity

import "time"

type Post struct {
	Id      int       `db:"id" json:"id"`
	Text    string    `db:"text"  json:"text"`
	UserId  int       `db:"user_id"  json:"userId"`
	Created time.Time `db:"created"  json:"created"`
	Updated time.Time `db:"updated" json:"updated"`
}

func NewPost(text string, user *User) *Post {
	return &Post{Text: text, UserId: user.Id, Created: time.Now(), Updated: time.Now()}
}
