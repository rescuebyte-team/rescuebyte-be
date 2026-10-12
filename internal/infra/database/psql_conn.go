package database

import (
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/rescuebyte/rescuebyte-be/internal/infra/env"
	"github.com/rescuebyte/rescuebyte-be/pkg/log"
)

func NewPgsqlConn() *sqlx.DB {
	dataSourceName := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
		env.AppEnv.DBHost,
		env.AppEnv.DBPort,
		env.AppEnv.DBUser,
		env.AppEnv.DBPass,
		env.AppEnv.DBName,
	)

	db, err := sqlx.Connect("pgx", dataSourceName)
	if err != nil {
		log.Panic(log.Fields{
			"error": err.Error(),
		}, "[DB][NewPgsqlConn] Failed to connect to PostgreSQL database")
	}

	db.SetMaxOpenConns(env.AppEnv.DBMaxConnection)
	db.SetMaxIdleConns(env.AppEnv.DBMaxIdleConnection)
	db.SetConnMaxLifetime(env.AppEnv.DBMaxLifetimeConnection)
	db.SetConnMaxIdleTime(env.AppEnv.DBMaxIdleTimeConnection)

	return db
}
