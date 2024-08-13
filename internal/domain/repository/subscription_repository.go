package repository

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/idprm/go-linkit-tsel/internal/domain/entity"
)

const (
	queryInsertSubscription           = "INSERT INTO subscriptions(category, service_id, msisdn, channel, camp_keyword, camp_sub_keyword, adnet, pub_id, aff_sub, latest_trxid, latest_keyword, latest_subject, latest_status,latest_pin, success, ip_address, is_active, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)"
	queryUpdateSubSuccess             = "UPDATE subscriptions SET latest_trxid = $1, latest_subject = $2, latest_status = $3, latest_pin = $4, amount = amount + $5, renewal_at = $6, charge_at = $7, success = success + $8, is_retry = $9, charging_count = charging_count + $10, charging_count_all = charging_count_all + $11, total_firstpush = total_firstpush + $12, total_renewal = total_renewal + $13, total_amount_firstpush = total_amount_firstpush + $14, total_amount_renewal = total_amount_renewal + $15, latest_payload = $16, updated_at = NOW() WHERE service_id = $17 AND msisdn = $18"
	queryUpdateSubFailed              = "UPDATE subscriptions SET latest_trxid = $1, latest_subject = $2, latest_status = $3, renewal_at = $4, retry_at = $5, failed = failed + $6, is_retry = $7, latest_payload = $8, updated_at = NOW() WHERE service_id = $9 AND msisdn = $10"
	queryUpdateSubLatest              = "UPDATE subscriptions SET latest_trxid = $1, latest_keyword = $2, latest_subject = $3, latest_status = $4, updated_at = NOW() WHERE service_id = $5 AND msisdn = $6"
	queryUpdateSubEnable              = "UPDATE subscriptions SET channel = $1, camp_keyword = $2, camp_sub_keyword = $3, adnet = $4, pub_id = $5, aff_sub = $6, latest_trxid = $7, latest_keyword = $8, latest_subject = $9, ip_address = $10, is_retry = $11, is_active = $12, updated_at = NOW(), created_at = NOW() WHERE service_id = $13 AND msisdn = $14"
	queryUpdateSubDisable             = "UPDATE subscriptions SET channel = $1, latest_trxid = $2, latest_keyword = $3, latest_subject = $4, latest_status = $5, unsub_at = $6, ip_address = $7, is_retry = $8, is_active = $9, charging_count = $10, updated_at = NOW() WHERE service_id = $11 AND msisdn = $12"
	queryUpdateSubConfirm             = "UPDATE subscriptions SET is_confirm = $1, updated_at = NOW() WHERE service_id = $2 AND msisdn = $3"
	queryUpdateSubPurge               = "UPDATE subscriptions SET purge_at = $1, purge_reason = $2, is_purge = true, is_active = false WHERE service_id = $3 AND msisdn = $4"
	queryUpdateSubLatestPayload       = "UPDATE subscriptions SET latest_payload = $1, updated_at = NOW() WHERE service_id = $2 AND msisdn = $3"
	queryUpdateSubPin                 = "UPDATE subscriptions SET latest_pin = $1, updated_at = NOW() WHERE service_id = $2 AND msisdn = $3"
	queryUpdateSubCampByToken         = "UPDATE subscriptions SET camp_keyword = $1, camp_sub_keyword = $2, adnet = $3, pub_id = $4, aff_sub = $5, is_trial = true, updated_at = NOW() WHERE latest_keyword = $6 AND DATE(created_at) = DATE(NOW()) AND camp_keyword = '' AND camp_sub_keyword = ''"
	queryUpdateSubSuccessRetry        = "UPDATE subscriptions SET latest_trxid = $1, latest_subject = $2, latest_status = $3, latest_pin = $4, amount = amount + $5, renewal_at = $6, charge_at = $7, success = success + $8, failed = failed - $9, is_retry = $10, charging_count = charging_count + $11, charging_count_all = charging_count_all + $12, total_firstpush = total_firstpush + $13, total_renewal = total_renewal + $14, total_amount_firstpush = total_amount_firstpush + $15, total_amount_renewal = total_amount_renewal + $16, latest_payload = $17, updated_at = NOW() WHERE service_id = $18 AND msisdn = $19"
	queryCountSubscription            = "SELECT COUNT(*) as count FROM subscriptions WHERE service_id = $1 AND msisdn = $2"
	queryCountActiveSubscription      = "SELECT COUNT(*) as count FROM subscriptions WHERE service_id = $1 AND msisdn = $2 AND is_active = true"
	queryCountPinSub                  = "SELECT COUNT(*) as count FROM subscriptions WHERE latest_pin = $1"
	queryCountSubActivePin            = "SELECT COUNT(*) as count FROM subscriptions WHERE category = $1 AND latest_pin = $2 AND is_active = true AND is_retry = false"
	querySelectSubscription           = "SELECT id, service_id, msisdn, channel, camp_keyword, camp_sub_keyword, adnet, pub_id, aff_sub, latest_trxid, latest_keyword, latest_subject, latest_status, latest_payload, amount, renewal_at, success, ip_address, total_firstpush, total_renewal, total_amount_firstpush, total_amount_renewal, is_retry, is_active FROM subscriptions WHERE service_id = $1 AND msisdn = $2"
	querySelectPopulateRenewal        = "SELECT id, service_id, msisdn, channel, adnet, latest_keyword, latest_subject, latest_pin, latest_payload, ip_address, aff_sub, camp_keyword, camp_sub_keyword, created_at FROM subscriptions WHERE renewal_at IS NOT NULL AND DATE(renewal_at) <= DATE(NOW()) AND is_active = true ORDER BY success DESC, DATE(created_at) DESC"
	querySelectPopulateRetryFirstpush = "SELECT id, service_id, msisdn, channel, adnet, latest_keyword, latest_subject, latest_pin, latest_payload, ip_address, aff_sub, camp_keyword, camp_sub_keyword, retry_at, created_at FROM subscriptions WHERE latest_payload <> '3:3:21' AND latest_subject = 'FIRSTPUSH' AND renewal_at IS NOT NULL AND DATE(renewal_at) = DATE(NOW() + interval '1 day') AND is_retry = true AND is_active = true ORDER BY success DESC, DATE(created_at) DESC"
	querySelectPopulateRetryDailypush = "SELECT id, service_id, msisdn, channel, adnet, latest_keyword, latest_subject, latest_pin, latest_payload, ip_address, aff_sub, camp_keyword, camp_sub_keyword, retry_at, created_at FROM subscriptions WHERE latest_payload <> '3:3:21' AND latest_subject = 'RENEWAL' AND renewal_at IS NOT NULL AND DATE(renewal_at) = DATE(NOW() + interval '1 day') AND is_retry = true AND is_active = true ORDER BY success DESC, DATE(created_at) DESC"
	querySelectPopulateRetryInsuff    = "SELECT id, service_id, msisdn, channel, adnet, latest_keyword, latest_subject, latest_pin, latest_payload, ip_address, aff_sub, camp_keyword, camp_sub_keyword, retry_at, created_at FROM subscriptions WHERE latest_payload = '3:3:21' AND renewal_at IS NOT NULL AND DATE(renewal_at) = DATE(NOW() + interval '1 day') AND is_retry = true AND is_active = true ORDER BY success DESC, DATE(created_at) DESC"
	querySelectPopulateReminder       = "SELECT id, service_id, msisdn, channel, latest_keyword, latest_pin, ip_address, aff_sub, camp_keyword, camp_sub_keyword, created_at FROM subscriptions WHERE renewal_at IS NOT NULL AND DATE(renewal_at) = DATE(NOW() + interval '2 day') AND is_retry = false AND is_active = true ORDER BY success DESC, DATE(created_at) DESC"
	querySelectPopulateTrial          = "SELECT id, service_id, msisdn, channel, latest_keyword, latest_pin, ip_address, camp_keyword, camp_sub_keyword, adnet, pub_id, aff_sub, created_at FROM subscriptions WHERE is_trial = true ORDER BY DATE(created_at) DESC"
	querySelectPopulateEmptyCamp      = "SELECT id, service_id, msisdn, channel, latest_keyword, latest_pin, ip_address, camp_keyword, camp_sub_keyword, adnet, pub_id, aff_sub, created_at FROM subscriptions WHERE camp_keyword = '' AND camp_sub_keyword = ''"
	querySelectArpu                   = "SELECT c.name, a.camp_sub_keyword, a.adnet, COUNT(a.id) as subs, COUNT(b.id) as subs_active, SUM(ROUND(a.amount)) as revenue FROM subscriptions a LEFT JOIN subscriptions b ON b.service_id = a.service_id AND b.msisdn = a.msisdn AND b.is_active = true LEFT JOIN services c ON c.id = a.service_id WHERE DATE(a.created_at) BETWEEN DATE($1) AND DATE($2) AND DATE(a.renewal_at) >= DATE($3) AND a.camp_sub_keyword = $4 GROUP BY a.adnet, a.camp_sub_keyword, a.service_id, c.name ORDER BY SUM(a.total_amount_renewal + a.total_amount_firstpush) DESC"
	querySelectSubCSV                 = "SELECT 'ID' as country, 'telkomsel' as operator, b.name as service, a.channel as source, a.msisdn, a.latest_subject, CONCAT(b.renewal_day, 'd') as cycle, a.adnet, a.total_amount_firstpush + a.total_amount_renewal as revenue, a.created_at as subs_date, a.renewal_at as renewal_date, '' as freemium_end_date, a.channel as unsubs_from, a.unsub_at as unsubs_date, b.price as service_price, 'IDR' as currency, a.is_active as profile_status, 'NA' as publisher, a.latest_trxid as trxid, 'NA' as pixel, 'NA' as handset, 'NA' as browser, success + failed as attempt_charging, success as success_billing, a.camp_sub_keyword FROM public.subscriptions a LEFT JOIN public.services b ON b.id = a.service_id WHERE DATE(a.updated_at) = DATE(NOW() - interval '1 day') ORDER BY a.id ASC"
	querySelectPurge                  = "SELECT id, service_id, msisdn, channel, adnet, latest_keyword, latest_subject, latest_pin, ip_address, aff_sub, camp_keyword, camp_sub_keyword, retry_at, created_at, purge_at, purge_reason FROM subscriptions WHERE is_purge = true"
	// queryUpdateSubscription      = "UPDATE subscriptions SET channel = $1, adnet = $2, pub_id = $3, aff_sub = $4, latest_trxid = $5, latest_keyword = $6, latest_subject = $7, latest_status = $8, amount = amount + $9, trial_at = $10, renewal_at = $11, unsub_at = $12, charge_at = $13, retry_at = $14, success = success + $15, ip_address = $16, is_trial = $17, is_retry = $18, is_active = $19, total_firstpush = total_firstpush + $20, total_renewal = total_renewal + $21, total_amount_firstpush = total_amount_firstpush + $22, total_amount_renewal = total_amount_renewal + $23, updated_at = NOW() WHERE service_id = $24 AND msisdn = $25"
)

type SubscriptionRepository struct {
	db *sql.DB
}

type ISubscriptionRepository interface {
	Save(*entity.Subscription) error
	UpdateSuccess(*entity.Subscription) error
	UpdateFailed(*entity.Subscription) error
	UpdateLatest(*entity.Subscription) error
	UpdateEnable(*entity.Subscription) error
	UpdateDisable(*entity.Subscription) error
	UpdateConfirm(*entity.Subscription) error
	UpdatePurge(*entity.Subscription) error
	UpdateLatestPayload(*entity.Subscription) error
	UpdatePin(*entity.Subscription) error
	UpdateCampByToken(*entity.Subscription) error
	UpdateSuccessRetry(*entity.Subscription) error
	Count(int, string) (int, error)
	CountActive(int, string) (int, error)
	CountPin(int) (int, error)
	CountPinActive(string, string) (int, error)
	Get(int, string) (*entity.Subscription, error)
	Renewal() (*[]entity.Subscription, error)
	RetryFp() (*[]entity.Subscription, error)
	RetryDp() (*[]entity.Subscription, error)
	RetryInsuff() (*[]entity.Subscription, error)
	Reminder() (*[]entity.Subscription, error)
	Trial() (*[]entity.Subscription, error)
	EmptyCamp() (*[]entity.Subscription, error)
	AveragePerUser(string, string, string, string) (*[]entity.AveragePerUser, error)
	SelectSubcriptionToCSV() (*[]entity.SubscriptionToCSV, error)
	SelectSubcriptionPurge() (*[]entity.Subscription, error)
}

func NewSubscriptionRepository(db *sql.DB) *SubscriptionRepository {
	return &SubscriptionRepository{
		db: db,
	}
}

func (r *SubscriptionRepository) Save(s *entity.Subscription) error {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelfunc()
	stmt, err := r.db.PrepareContext(ctx, queryInsertSubscription)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, s.Category, s.ServiceID, s.Msisdn, s.Channel, s.CampKeyword, s.CampSubKeyword, s.Adnet, s.PubID, s.AffSub, s.LatestTrxId, s.LatestKeyword, s.LatestSubject, s.LatestStatus, s.LatestPIN, s.Success, s.IpAddress, s.IsActive, time.Now(), time.Now())
	if err != nil {
		log.Printf("Error %s when inserting row into subscriptions table", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return err
	}
	log.Printf("%d subscriptions created ", rows)
	return nil
}

func (r *SubscriptionRepository) UpdateSuccess(s *entity.Subscription) error {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelfunc()
	stmt, err := r.db.PrepareContext(ctx, queryUpdateSubSuccess)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, s.LatestTrxId, s.LatestSubject, s.LatestStatus, s.LatestPIN, s.Amount, s.RenewalAt, s.ChargeAt, s.Success, s.IsRetry, s.ChargingCount, s.ChargingCountAll, s.TotalFirstpush, s.TotalRenewal, s.TotalAmountFirstpush, s.TotalAmountRenewal, s.LatestPayload, s.ServiceID, s.Msisdn)
	if err != nil {
		log.Printf("Error %s when update row into subscriptions table", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return err
	}
	log.Printf("%d subscriptions updated ", rows)

	return nil
}

func (r *SubscriptionRepository) UpdateFailed(s *entity.Subscription) error {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelfunc()
	stmt, err := r.db.PrepareContext(ctx, queryUpdateSubFailed)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, s.LatestTrxId, s.LatestSubject, s.LatestStatus, s.RenewalAt, s.RetryAt, s.Failed, s.IsRetry, s.LatestPayload, s.ServiceID, s.Msisdn)
	if err != nil {
		log.Printf("Error %s when update row into subscriptions table", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return err
	}
	log.Printf("%d subscriptions updated ", rows)

	return nil
}

func (r *SubscriptionRepository) UpdateLatest(s *entity.Subscription) error {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelfunc()
	stmt, err := r.db.PrepareContext(ctx, queryUpdateSubLatest)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, s.LatestTrxId, s.LatestKeyword, s.LatestSubject, s.LatestStatus, s.ServiceID, s.Msisdn)
	if err != nil {
		log.Printf("Error %s when update row into subscriptions table", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return err
	}
	log.Printf("%d subscriptions updated ", rows)

	return nil
}

func (r *SubscriptionRepository) UpdateEnable(s *entity.Subscription) error {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelfunc()
	stmt, err := r.db.PrepareContext(ctx, queryUpdateSubEnable)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, s.Channel, s.CampKeyword, s.CampSubKeyword, s.Adnet, s.PubID, s.AffSub, s.LatestTrxId, s.LatestKeyword, s.LatestSubject, s.IpAddress, s.IsRetry, s.IsActive, s.ServiceID, s.Msisdn)
	if err != nil {
		log.Printf("Error %s when update row into subscriptions table", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return err
	}
	log.Printf("%d subscriptions updated ", rows)

	return nil
}

func (r *SubscriptionRepository) UpdateDisable(s *entity.Subscription) error {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelfunc()
	stmt, err := r.db.PrepareContext(ctx, queryUpdateSubDisable)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, s.Channel, s.LatestTrxId, s.LatestKeyword, s.LatestSubject, s.LatestStatus, s.UnsubAt, s.IpAddress, s.IsRetry, s.IsActive, s.ChargingCount, s.ServiceID, s.Msisdn)
	if err != nil {
		log.Printf("Error %s when update row into subscriptions table", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return err
	}
	log.Printf("%d subscriptions updated ", rows)

	return nil
}

func (r *SubscriptionRepository) UpdateConfirm(s *entity.Subscription) error {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelfunc()
	stmt, err := r.db.PrepareContext(ctx, queryUpdateSubConfirm)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, s.IsConfirm, s.ServiceID, s.Msisdn)
	if err != nil {
		log.Printf("Error %s when update row into subscriptions table", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return err
	}
	log.Printf("%d subscriptions updated ", rows)

	return nil
}

func (r *SubscriptionRepository) UpdatePurge(s *entity.Subscription) error {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelfunc()
	stmt, err := r.db.PrepareContext(ctx, queryUpdateSubPurge)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, s.PurgeAt, s.PurgeReason, s.ServiceID, s.Msisdn)
	if err != nil {
		log.Printf("Error %s when update row into subscriptions table", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return err
	}
	log.Printf("%d subscriptions updated ", rows)

	return nil
}

func (r *SubscriptionRepository) UpdateLatestPayload(s *entity.Subscription) error {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelfunc()
	stmt, err := r.db.PrepareContext(ctx, queryUpdateSubLatestPayload)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, s.LatestPayload, s.ServiceID, s.Msisdn)
	if err != nil {
		log.Printf("Error %s when update row into subscriptions table", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return err
	}
	log.Printf("%d subscriptions updated ", rows)

	return nil
}

func (r *SubscriptionRepository) UpdatePin(s *entity.Subscription) error {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelfunc()
	stmt, err := r.db.PrepareContext(ctx, queryUpdateSubPin)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, s.LatestPIN, s.ServiceID, s.Msisdn)
	if err != nil {
		log.Printf("Error %s when update row into subscriptions table", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return err
	}
	log.Printf("%d subscriptions updated ", rows)

	return nil
}

func (r *SubscriptionRepository) UpdateCampByToken(s *entity.Subscription) error {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelfunc()
	stmt, err := r.db.PrepareContext(ctx, queryUpdateSubCampByToken)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, s.CampKeyword, s.CampSubKeyword, s.Adnet, s.PubID, s.AffSub, s.LatestKeyword)
	if err != nil {
		log.Printf("Error %s when update row into subscriptions table", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return err
	}
	log.Printf("%d subscriptions updated ", rows)

	return nil
}

func (r *SubscriptionRepository) UpdateSuccessRetry(s *entity.Subscription) error {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelfunc()
	stmt, err := r.db.PrepareContext(ctx, queryUpdateSubSuccessRetry)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, s.LatestTrxId, s.LatestSubject, s.LatestStatus, s.LatestPIN, s.Amount, s.RenewalAt, s.ChargeAt, s.Success, s.Failed, s.IsRetry, s.ChargingCount, s.ChargingCountAll, s.TotalFirstpush, s.TotalRenewal, s.TotalAmountFirstpush, s.TotalAmountRenewal, s.LatestPayload, s.ServiceID, s.Msisdn)
	if err != nil {
		log.Printf("Error %s when update row into subscriptions table", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return err
	}
	log.Printf("%d subscriptions updated ", rows)

	return nil
}

func (r *SubscriptionRepository) Count(serviceId int, msisdn string) (int, error) {
	var count int
	err := r.db.QueryRow(queryCountSubscription, serviceId, msisdn).Scan(&count)
	if err != nil {
		return count, err
	}
	return count, nil
}

func (r *SubscriptionRepository) CountActive(serviceId int, msisdn string) (int, error) {
	var count int
	err := r.db.QueryRow(queryCountActiveSubscription, serviceId, msisdn).Scan(&count)
	if err != nil {
		return count, err
	}
	return count, nil
}

func (r *SubscriptionRepository) CountPin(pin int) (int, error) {
	var count int
	err := r.db.QueryRow(queryCountPinSub, pin).Scan(&count)
	if err != nil {
		return count, err
	}
	return count, nil
}

func (r *SubscriptionRepository) CountPinActive(category, pin string) (int, error) {
	var count int
	err := r.db.QueryRow(queryCountSubActivePin, category, pin).Scan(&count)
	if err != nil {
		return count, err
	}
	return count, nil
}

func (r *SubscriptionRepository) Get(serviceId int, msisdn string) (*entity.Subscription, error) {
	var s entity.Subscription
	err := r.db.QueryRow(querySelectSubscription, serviceId, msisdn).Scan(&s.ID, &s.ServiceID, &s.Msisdn, &s.Channel, &s.CampKeyword, &s.CampSubKeyword, &s.Adnet, &s.PubID, &s.AffSub, &s.LatestTrxId, &s.LatestKeyword, &s.LatestSubject, &s.LatestStatus, &s.LatestPayload, &s.Amount, &s.RenewalAt, &s.Success, &s.IpAddress, &s.TotalFirstpush, &s.TotalRenewal, &s.TotalAmountFirstpush, &s.TotalAmountRenewal, &s.IsRetry, &s.IsActive)
	if err != nil {
		return &s, err
	}
	return &s, nil
}

func (r *SubscriptionRepository) Renewal() (*[]entity.Subscription, error) {
	rows, err := r.db.Query(querySelectPopulateRenewal)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []entity.Subscription

	for rows.Next() {

		var s entity.Subscription
		if err := rows.Scan(&s.ID, &s.ServiceID, &s.Msisdn, &s.Channel, &s.Adnet, &s.LatestKeyword, &s.LatestSubject, &s.LatestPIN, &s.LatestPayload, &s.IpAddress, &s.AffSub, &s.CampKeyword, &s.CampSubKeyword, &s.CreatedAt); err != nil {
			return nil, err
		}
		subs = append(subs, s)
	}

	if err = rows.Err(); err != nil {
		return &subs, err
	}

	return &subs, nil
}

func (r *SubscriptionRepository) RetryFp() (*[]entity.Subscription, error) {
	rows, err := r.db.Query(querySelectPopulateRetryFirstpush)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []entity.Subscription

	for rows.Next() {

		var s entity.Subscription
		if err := rows.Scan(&s.ID, &s.ServiceID, &s.Msisdn, &s.Channel, &s.Adnet, &s.LatestKeyword, &s.LatestSubject, &s.LatestPIN, &s.LatestPayload, &s.IpAddress, &s.AffSub, &s.CampKeyword, &s.CampSubKeyword, &s.RetryAt, &s.CreatedAt); err != nil {
			return nil, err
		}
		subs = append(subs, s)
	}

	if err = rows.Err(); err != nil {
		return &subs, err
	}

	return &subs, nil
}

func (r *SubscriptionRepository) RetryDp() (*[]entity.Subscription, error) {
	rows, err := r.db.Query(querySelectPopulateRetryDailypush)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []entity.Subscription

	for rows.Next() {

		var s entity.Subscription
		if err := rows.Scan(&s.ID, &s.ServiceID, &s.Msisdn, &s.Channel, &s.Adnet, &s.LatestKeyword, &s.LatestSubject, &s.LatestPIN, &s.LatestPayload, &s.IpAddress, &s.AffSub, &s.CampKeyword, &s.CampSubKeyword, &s.RetryAt, &s.CreatedAt); err != nil {
			return nil, err
		}
		subs = append(subs, s)
	}

	if err = rows.Err(); err != nil {
		return &subs, err
	}

	return &subs, nil
}

func (r *SubscriptionRepository) RetryInsuff() (*[]entity.Subscription, error) {
	rows, err := r.db.Query(querySelectPopulateRetryInsuff)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []entity.Subscription

	for rows.Next() {

		var s entity.Subscription
		if err := rows.Scan(&s.ID, &s.ServiceID, &s.Msisdn, &s.Channel, &s.Adnet, &s.LatestKeyword, &s.LatestSubject, &s.LatestPIN, &s.LatestPayload, &s.IpAddress, &s.AffSub, &s.CampKeyword, &s.CampSubKeyword, &s.RetryAt, &s.CreatedAt); err != nil {
			return nil, err
		}
		subs = append(subs, s)
	}

	if err = rows.Err(); err != nil {
		return &subs, err
	}

	return &subs, nil
}

func (r *SubscriptionRepository) Reminder() (*[]entity.Subscription, error) {
	rows, err := r.db.Query(querySelectPopulateReminder)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []entity.Subscription

	for rows.Next() {

		var s entity.Subscription
		if err := rows.Scan(&s.ID, &s.ServiceID, &s.Msisdn, &s.Channel, &s.LatestKeyword, &s.LatestPIN, &s.IpAddress, &s.AffSub, &s.CampKeyword, &s.CampSubKeyword, &s.CreatedAt); err != nil {
			return nil, err
		}
		subs = append(subs, s)
	}

	if err = rows.Err(); err != nil {
		return &subs, err
	}

	return &subs, nil
}

func (r *SubscriptionRepository) Trial() (*[]entity.Subscription, error) {
	rows, err := r.db.Query(querySelectPopulateTrial)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []entity.Subscription

	for rows.Next() {
		var s entity.Subscription
		if err := rows.Scan(&s.ID, &s.ServiceID, &s.Msisdn, &s.Channel, &s.LatestKeyword, &s.LatestPIN, &s.IpAddress, &s.CampKeyword, &s.CampSubKeyword, &s.Adnet, &s.PubID, &s.AffSub, &s.CreatedAt); err != nil {
			return nil, err
		}
		subs = append(subs, s)
	}

	if err = rows.Err(); err != nil {
		return &subs, err
	}

	return &subs, nil
}

func (r *SubscriptionRepository) EmptyCamp() (*[]entity.Subscription, error) {
	rows, err := r.db.Query(querySelectPopulateEmptyCamp)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []entity.Subscription

	for rows.Next() {
		var s entity.Subscription
		if err := rows.Scan(&s.ID, &s.ServiceID, &s.Msisdn, &s.Channel, &s.LatestKeyword, &s.LatestPIN, &s.IpAddress, &s.CampKeyword, &s.CampSubKeyword, &s.Adnet, &s.PubID, &s.AffSub, &s.CreatedAt); err != nil {
			return nil, err
		}
		subs = append(subs, s)
	}

	if err = rows.Err(); err != nil {
		return &subs, err
	}

	return &subs, nil
}

func (r *SubscriptionRepository) AveragePerUser(start, end, renewal, subkey string) (*[]entity.AveragePerUser, error) {
	rows, err := r.db.Query(querySelectArpu, start, end, renewal, subkey)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []entity.AveragePerUser

	for rows.Next() {
		var s entity.AveragePerUser
		if err := rows.Scan(&s.Name, &s.Service, &s.Adnet, &s.Subs, &s.SubsActive, &s.Revenue); err != nil {
			return nil, err
		}
		subs = append(subs, s)
	}

	if err = rows.Err(); err != nil {
		return &subs, err
	}

	return &subs, nil
}

func (r *SubscriptionRepository) SelectSubcriptionToCSV() (*[]entity.SubscriptionToCSV, error) {
	rows, err := r.db.Query(querySelectSubCSV)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []entity.SubscriptionToCSV

	for rows.Next() {
		var s entity.SubscriptionToCSV
		if err := rows.Scan(&s.Country, &s.Operator, &s.Service, &s.Source, &s.Msisdn, &s.LatestSubject, &s.Cycle, &s.Adnet, &s.Revenue, &s.SubsDate, &s.RenewalDate, &s.FreemiumEndDate, &s.UnsubsFrom, &s.UnsubsDate, &s.ServicePrice, &s.Currency, &s.ProfileStatus, &s.Publisher, &s.Trxid, &s.Pixel, &s.Handset, &s.Browser, &s.AttemptCharging, &s.SuccessBilling, &s.CampSubKeyword); err != nil {
			return nil, err
		}
		subs = append(subs, s)
	}

	if err = rows.Err(); err != nil {
		return &subs, err
	}

	return &subs, nil
}

func (r *SubscriptionRepository) SelectSubcriptionPurge() (*[]entity.Subscription, error) {
	rows, err := r.db.Query(querySelectPurge)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []entity.Subscription

	for rows.Next() {
		var s entity.Subscription
		if err := rows.Scan(&s.ID, &s.ServiceID, &s.Msisdn, &s.Channel, &s.Adnet, &s.LatestKeyword, &s.LatestSubject, &s.LatestPIN, &s.IpAddress, &s.AffSub, &s.CampKeyword, &s.CampSubKeyword, &s.RetryAt, &s.CreatedAt, &s.PurgeAt, &s.PurgeReason); err != nil {
			return nil, err
		}
		subs = append(subs, s)
	}

	if err = rows.Err(); err != nil {
		return &subs, err
	}

	return &subs, nil
}
