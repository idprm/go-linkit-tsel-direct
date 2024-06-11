package repository

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/idprm/go-linkit-tsel/internal/domain/entity"
)

const (
	queryInsertTrafficCampaign = "INSERT INTO traffics_campaign(service_id, camp_keyword, camp_sub_keyword, adnet, pub_id, aff_sub, browser, os, device, ip_address, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)"
	queryInsertTrafficMO       = "INSERT INTO traffics_mo(service_id, msisdn, channel, camp_keyword, camp_sub_keyword, adnet, pub_id, aff_sub, ip_address, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)"
)

type TrafficRepository struct {
	db *sql.DB
}

type ITrafficRepository interface {
	SaveCampaign(t *entity.TrafficCampaign) error
	SaveMO(t *entity.TrafficMO) error
}

func NewTrafficRepository(db *sql.DB) *TrafficRepository {
	return &TrafficRepository{
		db: db,
	}
}

func (r *TrafficRepository) SaveCampaign(t *entity.TrafficCampaign) error {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelfunc()
	stmt, err := r.db.PrepareContext(ctx, queryInsertTrafficCampaign)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, t.ServiceID, t.CampKeyword, t.CampSubKeyword, t.Adnet, t.PubID, t.AffSub, t.Browser, t.OS, t.Device, t.IpAddress, time.Now())
	if err != nil {
		log.Printf("Error %s when inserting row into traffics_campaign table", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return err
	}
	log.Printf("%d traffics_campaign created ", rows)
	return nil
}

func (r *TrafficRepository) SaveMO(t *entity.TrafficMO) error {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelfunc()
	stmt, err := r.db.PrepareContext(ctx, queryInsertTrafficMO)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, t.ServiceID, t.Msisdn, t.Channel, t.CampKeyword, t.CampSubKeyword, t.Adnet, t.PubID, t.AffSub, t.IpAddress, time.Now())
	if err != nil {
		log.Printf("Error %s when inserting row into traffics_campaign table", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return err
	}
	log.Printf("%d traffics_campaign created ", rows)
	return nil
}
