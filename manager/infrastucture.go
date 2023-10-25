package manager

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/terajari/bank-api/utils"
)

type InfrastuctureManager interface {
	Conn() *sqlx.DB
	Config() *utils.Config
}

type infrastuctureManager struct {
	db     *sqlx.DB
	config *utils.Config
}

func (i *infrastuctureManager) initDb() error {
	driver := i.config.DBDriver
	source := i.config.DBSource

	db, err := sqlx.Open(driver, source)
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}
	i.db = db

	return nil

}

func (i *infrastuctureManager) Conn() *sqlx.DB {
	return i.db
}

func (i *infrastuctureManager) Config() *utils.Config {
	return i.config
}

func NewInfraManager(configParam *utils.Config) (InfrastuctureManager, error) {
	infra := &infrastuctureManager{
		config: configParam,
	}

	err := infra.initDb()
	if err != nil {
		return nil, err
	}
	return infra, nil
}
