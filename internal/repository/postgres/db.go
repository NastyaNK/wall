package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"wall/internal/entity"
	"wall/internal/repository"
)

type Database struct {
	db *sqlx.DB
}

func NewDatabase() repository.Repository {
	return &Database{}
}

func (psql *Database) Connect(config *entity.DBConfig) error {
	var err error
	psql.db, err = sqlx.Connect("postgres", fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s sslmode=%s",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Dbname,
		config.Sslmode))
	return err
}

func (psql *Database) Close() error {
	return psql.db.Close()
}
