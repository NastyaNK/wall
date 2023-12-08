package entity

type User struct {
	Id       int    `db:"id" json:"id"`
	Name     string `db:"name" json:"name"`
	Age      int    `db:"age" json:"age"`
	Password string `db:"password" json:"password"`
}
