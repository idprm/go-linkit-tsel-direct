package repository

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/idprm/go-linkit-tsel/internal/domain/entity"
)

const (
	queryInsertHistory = "INSERT INTO histories(service_id, msisdn, channel, adnet, keyword, subject, status, ip_address, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)"
)

type HistoryRepository struct {
	db *sql.DB
}

func NewHistoryRepository(db *sql.DB) *HistoryRepository {
	return &HistoryRepository{
		db: db,
	}
}

type IHistoryRepository interface {
	Save(*entity.History) error
}

func (r *HistoryRepository) Save(h *entity.History) error {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelfunc()
	stmt, err := r.db.PrepareContext(ctx, queryInsertHistory)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, h.ServiceID, h.Msisdn, h.Channel, h.Adnet, h.Keyword, h.Subject, h.Status, h.IpAddress, time.Now())
	if err != nil {
		log.Printf("Error %s when inserting row into histories table", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return err
	}
	log.Printf("%d histories created ", rows)
	return nil
}
