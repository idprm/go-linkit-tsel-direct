package repository

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/idprm/go-linkit-tsel/internal/domain/entity"
)

const (
	queryInsertTrafficCampaign = "INSERT INTO traffics_campaign(tx_id, service_id, camp_keyword, camp_sub_keyword, adnet, pub_id, aff_sub, browser, os, device, referer, ip_address, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)"
	queryInsertTrafficMO       = "INSERT INTO traffics_mo(tx_id, service_id, msisdn, channel, subject, camp_keyword, camp_sub_keyword, adnet, pub_id, aff_sub, is_charge, ip_address, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)"
	queryUpdateTrafficMOCharge = "UPDATE traffics_mo SET is_charge = $1 WHERE service_id = $2 AND msisdn = $3 AND DATE(created_at) = DATE($4)"
)

type TrafficRepository struct {
	db *sql.DB
}

type ITrafficRepository interface {
	SaveCampaign(*entity.TrafficCampaign) error
	SaveMO(*entity.TrafficMO) error
	UpdateMOCharge(*entity.TrafficMO) error
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
	res, err := stmt.ExecContext(ctx, t.TxId, t.ServiceID, t.CampKeyword, t.CampSubKeyword, t.Adnet, t.PubID, t.AffSub, t.Browser, t.OS, t.Device, t.Referer, t.IpAddress, time.Now())
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
	res, err := stmt.ExecContext(ctx, t.TxId, t.ServiceID, t.Msisdn, t.Channel, t.Subject, t.CampKeyword, t.CampSubKeyword, t.Adnet, t.PubID, t.AffSub, t.IsCharge, t.IpAddress, time.Now())
	if err != nil {
		log.Printf("Error %s when inserting row into traffics_mo table", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return err
	}
	log.Printf("%d traffics_mo created ", rows)
	return nil
}

func (r *TrafficRepository) UpdateMOCharge(t *entity.TrafficMO) error {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelfunc()
	stmt, err := r.db.PrepareContext(ctx, queryUpdateTrafficMOCharge)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return err
	}

	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, t.IsCharge, t.ServiceID, t.Msisdn, time.Now())
	if err != nil {
		log.Printf("Error %s when update row into traffics_mo table", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return err
	}
	log.Printf("%d traffics_mo updated ", rows)

	return nil
}
