package repository

import (
	"database/sql"

	"github.com/idprm/go-linkit-tsel/internal/domain/entity"
)

const (
	queryCountPostbackBySubKey  = "SELECT COUNT(*) as count FROM postbacks WHERE sub_keyword = $1 AND is_active = true"
	querySelectPostbackBySubKey = "SELECT id, sub_keyword, url_mo, url_dn, is_active FROM postbacks WHERE sub_keyword = $1 AND is_active = true LIMIT 1"
)

type PostbackRepository struct {
	db *sql.DB
}

type IPostbackRepository interface {
	CountBySubkey(string) (int, error)
	GetBySubKey(string) (*entity.Postback, error)
}

func NewPostbackRepository(db *sql.DB) *PostbackRepository {
	return &PostbackRepository{
		db: db,
	}
}

func (r *PostbackRepository) CountBySubkey(subkey string) (int, error) {
	var count int
	err := r.db.QueryRow(queryCountPostbackBySubKey, subkey).Scan(&count)
	if err != nil {
		return count, err
	}
	return count, nil
}

func (r *PostbackRepository) GetBySubKey(subkey string) (*entity.Postback, error) {
	var s entity.Postback
	err := r.db.QueryRow(querySelectPostbackBySubKey, subkey).Scan(&s.ID, &s.SubKeyword, &s.UrlMO, &s.UrlDN, &s.IsActive)
	if err != nil {
		return &s, err
	}
	return &s, nil
}
