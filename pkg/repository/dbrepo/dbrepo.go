package dbrepo

import (
	"database/sql"
	"yamifood/pkg/config"
	"yamifood/pkg/repository"
)

type postgresDbRepository struct {
	App *config.AppConfig
	Db  *sql.DB
}

func NewPostgresRepo(conn *sql.DB, ac *config.AppConfig) repository.DatabaseRepo {
	return &postgresDbRepository{
		App: ac,
		Db:  conn,
	}
}
