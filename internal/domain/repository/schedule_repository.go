package repository

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/idprm/go-linkit-tsel/internal/domain/entity"
)

const (
	queryCountUnlocked  = "SELECT COUNT(*) as count FROM schedules WHERE name = $1 AND to_char(publish_at, 'HH24:MI') = $2 AND is_unlocked = true"
	queryCountLocked    = "SELECT COUNT(*) as count FROM schedules WHERE name = $1 AND to_char(unlocked_at, 'HH24:MI') = $2 AND is_unlocked = false"
	queryUpdateSchedule = "UPDATE schedules SET is_unlocked = $1 WHERE name = $2"
)

type ScheduleRepository struct {
	db *sql.DB
}

type IScheduleRepository interface {
	CountUnlocked(string, string) (int, error)
	CountLocked(string, string) (int, error)
	Update(*entity.Schedule) error
}

func NewScheduleRepository(db *sql.DB) *ScheduleRepository {
	return &ScheduleRepository{
		db: db,
	}
}

func (r *ScheduleRepository) CountUnlocked(name, hour string) (int, error) {
	var count int
	err := r.db.QueryRow(queryCountUnlocked, name, hour).Scan(&count)
	if err != nil {
		return count, err
	}
	return count, nil
}

func (r *ScheduleRepository) CountLocked(name, hour string) (int, error) {
	var count int
	err := r.db.QueryRow(queryCountLocked, name, hour).Scan(&count)
	if err != nil {
		return count, err
	}
	return count, nil
}

func (r *ScheduleRepository) Update(s *entity.Schedule) error {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelfunc()
	stmt, err := r.db.PrepareContext(ctx, queryUpdateSchedule)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, s.IsUnlocked, s.Name)
	if err != nil {
		log.Printf("Error %s when update row into schedules table", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return err
	}
	log.Printf("%d schedules updated ", rows)

	return nil
}
