package repository

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/idprm/go-linkit-tsel/internal/domain/entity"
)

const (
	queryInsertDailypush = "INSERT INTO dailypushes(tx_id, subscription_id, service_id, msisdn, channel, camp_keyword, camp_sub_keyword, adnet, pub_id, aff_sub, subject, status_code, status_detail, is_charge, ip_address, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)"
	queryDeleteDailypush = "DELETE FROM dailypushes WHERE service_id = $1 AND msisdn = $2 AND subject = $3 AND is_charge = $4 AND DATE(created_at) = DATE($5)"
)

type DailypushRepository struct {
	db *sql.DB
}

type IDailypushRepository interface {
	Save(*entity.Dailypush) error
	Delete(*entity.Dailypush) error
}

func NewDailypushRepository(db *sql.DB) *DailypushRepository {
	return &DailypushRepository{
		db: db,
	}
}

func (r *DailypushRepository) Save(t *entity.Dailypush) error {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelfunc()
	stmt, err := r.db.PrepareContext(ctx, queryInsertDailypush)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, t.TxId, t.SubscriptionID, t.ServiceID, t.Msisdn, t.Channel, t.CampKeyword, t.CampSubKeyword, t.Adnet, t.PubID, t.AffSub, t.Subject, t.StatusCode, t.StatusDetail, t.IsCharge, t.IpAddress, time.Now(), time.Now())
	if err != nil {
		log.Printf("Error %s when inserting row into dailypushes table", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return err
	}
	log.Printf("%d dailypushes created ", rows)
	return nil
}

func (r *DailypushRepository) Delete(t *entity.Dailypush) error {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelfunc()
	stmt, err := r.db.PrepareContext(ctx, queryDeleteDailypush)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, t.ServiceID, t.Msisdn, t.Subject, t.IsCharge, time.Now())
	if err != nil {
		log.Printf("Error %s when remove row into dailypushes table", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return err
	}
	log.Printf("%d dailypushes deleted ", rows)
	return nil
}
