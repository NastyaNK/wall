package postgres

import (
	"wall/internal/entity"
)

func (psql *Database) AddUser(user *entity.User) (int, error) {
	rows, err := psql.db.NamedQuery("insert into users(name, age, password) VALUES (:name,:age,:password) RETURNING id", user)
	if err != nil {
		return 0, err
	}
	var id int
	rows.Next()
	err = rows.Scan(&id)
	return id, err
}

func (psql *Database) GetUser(name string) (*entity.User, error) {
	var user entity.User
	err := psql.db.Get(&user, "SELECT * FROM users WHERE name=$1", name)
	return &user, err
}
