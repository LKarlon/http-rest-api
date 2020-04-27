package store

import (
	"github.com/LKarlon/http-rest-api.git/api/models"
)

type InnRepo struct {
	store *Store
}

func (r *InnRepo) Add(passport, inn string) error {
	_, err := r.store.db.Exec(
		"INSERT INTO inn_test.inn_shema.ready_data (passport, inn) VALUES ($1, $2)",
		passport,
		inn,
	)
	if err != nil{
		return err
	}
	return nil
}

func (r *InnRepo) FindInn(passport string) (models.INNReady, error) {
	m := models.INNReady{}
	if err := r.store.db.QueryRow(
		"SELECT passport, inn FROM inn_test.inn_shema.ready_data WHERE passport = $1",
		passport,
		).Scan(&m.Inn, &m.Passport); err != nil{
		return m, err
	}
	return m, nil
}
