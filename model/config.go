package model

import (
	"database/sql"

	"go.uber.org/zap"
)

type ServerConfig struct {
	Addr string
	Db   *sql.DB
	Log  *zap.Logger
}
