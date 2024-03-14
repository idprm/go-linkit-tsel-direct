package repository

import (
	"database/sql"
)

const (
	queryCountBlacklist = "SELECT COUNT(*) as count FROM blacklists WHERE msisdn = $1"
)

type BlacklistRepository struct {
	db *sql.DB
}

type IBlacklistRepository interface {
	Count(string) (int, error)
}

func NewBlacklistRepository(db *sql.DB) *BlacklistRepository {
	return &BlacklistRepository{
		db: db,
	}
}

func (r *BlacklistRepository) Count(msisdn string) (int, error) {
	var count int
	err := r.db.QueryRow(queryCountBlacklist, msisdn).Scan(&count)
	if err != nil {
		return count, err
	}
	return count, nil
}
