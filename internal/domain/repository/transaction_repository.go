package repository

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/idprm/go-linkit-tsel/internal/domain/entity"
)

const (
	queryInsertTransaction         = "INSERT INTO transactions(tx_id, service_id, msisdn, channel, camp_keyword, camp_sub_keyword, adnet, pub_id, aff_sub, keyword, pin, amount, status, status_code, status_detail, subject, ip_address, payload, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)"
	queryDeleteTransaction         = "DELETE FROM transactions WHERE service_id = $1 AND msisdn = $2 AND subject = $3 AND status = $4 AND DATE(created_at) = DATE($5)"
	querySelectGroupByStatus       = "SELECT DATE(created_at), service_id, subject, status, COUNT(msisdn) as count, SUM(ROUND(amount)) as amount FROM transactions WHERE DATE(created_at) >= DATE(NOW() + interval '-7 day') GROUP BY service_id, status, DATE(created_at), subject ORDER BY DATE(created_at) DESC, status DESC, subject ASC LIMIT 35"
	querySelectGroupByStatusDetail = "SELECT DATE(created_at), service_id, subject, status, status_detail, COUNT(msisdn) as count FROM transactions WHERE DATE(created_at) >= DATE(NOW() + interval '-7 day') GROUP BY service_id, status, status_detail, DATE(created_at), subject ORDER BY DATE(created_at) DESC LIMIT 35"
	querySelectGroupByAdnet        = "SELECT DATE(created_at), service_id, adnet, COUNT(msisdn) as count FROM transactions WHERE DATE(created_at) >= DATE(NOW() + interval '-7 day') AND subject = 'FIRSTPUSH' GROUP BY service_id, adnet, DATE(created_at) ORDER BY DATE(created_at) DESC LIMIT 35"
	querySelectTransactionCSV      = "SELECT 'ID' as country, 'telkomsel' as operator, b.name as service, a.channel as source, a.msisdn, a.subject as event, a.created_at as even_date, CONCAT(b.renewal_day, 'd') as cycle, a.amount as revenue, a.updated_at as charge_date, 'IDR' as currency, 'NA' as publisher, 'NA' as handset, 'NA' as browser, a.tx_id as trxid, b.url_telco as telco_api_url, COALESCE(a.payload, 'NA') as telco_api_response, COALESCE(c.value, 'NA') as sms_content, a.status as status_sms FROM public.transactions a LEFT JOIN public.services b ON b.id = a.service_id LEFT JOIN public.contents c ON c.service_id = b.id AND c.name = a.subject WHERE DATE(a.created_at) = DATE(NOW() - interval '1 day') ORDER BY a.id ASC"
	//querySelectTransactionCSV = "SELECT 'ID' as country, 'telkomsel' as operator, b.name as service, a.channel as source, a.msisdn, a.subject as event, a.created_at as even_date, CONCAT(b.renewal_day, 'd') as cycle, a.amount as revenue, a.updated_at as charge_date, 'IDR' as currency, 'NA' as publisher, 'NA' as handset, 'NA' as browser, a.tx_id as trxid, b.url_telco as telco_api_url, COALESCE(a.payload, 'NA') as telco_api_response, COALESCE(c.value, 'NA') as sms_content, a.status as status_sms FROM public.transactions a LEFT JOIN public.services b ON b.id = a.service_id LEFT JOIN public.contents c ON c.service_id = b.id AND c.name = a.subject WHERE DATE(a.created_at) = DATE('2024-03-02') ORDER BY a.id ASC"
)

type TransactionRepository struct {
	db *sql.DB
}

type ITransactionRepository interface {
	Save(*entity.Transaction) error
	Delete(*entity.Transaction) error
	SelectByStatus() (*[]entity.Transaction, error)
	SelectByStatusDetail() (*[]entity.Transaction, error)
	SelectByAdnet() (*[]entity.Transaction, error)
	SelectTransactionToCSV() (*[]entity.TransactionToCSV, error)
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{
		db: db,
	}
}

func (r *TransactionRepository) Save(t *entity.Transaction) error {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelfunc()
	stmt, err := r.db.PrepareContext(ctx, queryInsertTransaction)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, t.TxID, t.ServiceID, t.Msisdn, t.Channel, t.CampKeyword, t.CampSubKeyword, t.Adnet, t.PubID, t.AffSub, t.Keyword, t.PIN, t.Amount, t.Status, t.StatusCode, t.StatusDetail, t.Subject, t.IpAddress, t.Payload, time.Now(), time.Now())
	if err != nil {
		log.Printf("Error %s when inserting row into transactions table", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return err
	}
	log.Printf("%d transactions created ", rows)
	return nil
}

func (r *TransactionRepository) Delete(t *entity.Transaction) error {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelfunc()
	stmt, err := r.db.PrepareContext(ctx, queryDeleteTransaction)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, t.ServiceID, t.Msisdn, t.Subject, t.Status, time.Now())
	if err != nil {
		log.Printf("Error %s when remove row into transactions table", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return err
	}
	log.Printf("%d transactions deleted ", rows)
	return nil
}

func (r *TransactionRepository) SelectByStatus() (*[]entity.Transaction, error) {
	rows, err := r.db.Query(querySelectGroupByStatus)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var trans []entity.Transaction
	for rows.Next() {
		var t entity.Transaction
		if err := rows.Scan(&t.CreatedAt, &t.ServiceID, &t.Subject, &t.Status, &t.Msisdn, &t.Amount); err != nil {
			return nil, err
		}
		trans = append(trans, t)
	}
	if err = rows.Err(); err != nil {
		return &trans, err
	}
	return &trans, nil
}

func (r *TransactionRepository) SelectByStatusDetail() (*[]entity.Transaction, error) {
	rows, err := r.db.Query(querySelectGroupByStatusDetail)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var trans []entity.Transaction
	for rows.Next() {
		var t entity.Transaction
		if err := rows.Scan(&t.CreatedAt, &t.ServiceID, &t.Msisdn, &t.Subject, &t.Status, &t.StatusDetail); err != nil {
			return nil, err
		}
		trans = append(trans, t)
	}
	if err = rows.Err(); err != nil {
		return &trans, err
	}
	return &trans, nil
}

func (r *TransactionRepository) SelectByAdnet() (*[]entity.Transaction, error) {
	rows, err := r.db.Query(querySelectGroupByAdnet)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var trans []entity.Transaction
	for rows.Next() {
		var t entity.Transaction
		if err := rows.Scan(&t.CreatedAt, &t.ServiceID, &t.Msisdn, &t.Adnet); err != nil {
			return nil, err
		}
		trans = append(trans, t)
	}
	if err = rows.Err(); err != nil {
		return &trans, err
	}
	return &trans, nil
}

func (r *TransactionRepository) SelectTransactionToCSV() (*[]entity.TransactionToCSV, error) {
	rows, err := r.db.Query(querySelectTransactionCSV)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trans []entity.TransactionToCSV

	for rows.Next() {
		var t entity.TransactionToCSV
		if err := rows.Scan(&t.Country, &t.Operator, &t.Service, &t.Source, &t.Msisdn, &t.Event, &t.EventDate, &t.Cycle, &t.Revenue, &t.ChargeDate, &t.Currency, &t.Publisher, &t.Handset, &t.Browser, &t.TrxId, &t.TelcoApiUrl, &t.TelcoApiResponse, &t.SmsContent, &t.StatusSms); err != nil {
			return nil, err
		}
		trans = append(trans, t)
	}

	if err = rows.Err(); err != nil {
		return &trans, err
	}

	return &trans, nil
}
