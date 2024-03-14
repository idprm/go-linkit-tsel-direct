package repository

import (
	"database/sql"

	"github.com/idprm/go-linkit-tsel/internal/domain/entity"
)

const (
	querySelectContent = "SELECT value, tid FROM contents WHERE service_id = $1 AND name = $2 LIMIT 1"
)

type ContentRepository struct {
	db *sql.DB
}

type IContentRepository interface {
	Get(int, string) (*entity.Content, error)
}

func NewContentRepository(db *sql.DB) *ContentRepository {
	return &ContentRepository{
		db: db,
	}
}

func (r *ContentRepository) Get(serviceId int, name string) (*entity.Content, error) {
	var content entity.Content
	err := r.db.QueryRow(querySelectContent, serviceId, name).Scan(&content.Value, &content.Tid)
	if err != nil {
		return &content, err
	}
	return &content, nil
}
