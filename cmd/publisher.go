package cmd

import (
	"compress/gzip"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/idprm/go-linkit-tsel/internal/domain/entity"
	"github.com/idprm/go-linkit-tsel/internal/domain/repository"
	"github.com/idprm/go-linkit-tsel/internal/logger"
	"github.com/idprm/go-linkit-tsel/internal/providers/arpu"
	"github.com/idprm/go-linkit-tsel/internal/providers/rabbit"
	"github.com/idprm/go-linkit-tsel/internal/services"
	"github.com/spf13/cobra"
	"github.com/wiliehidayat87/rmqp"
)

var publisherRenewalCmd = &cobra.Command{
	Use:   "pub_renewal",
	Short: "Renewal CLI",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		/**
		 * connect pgsql
		 */
		db, err := connectPgsql()
		if err != nil {
			panic(err)
		}

		/**
		 * connect rabbitmq
		 */
		rmq := connectRabbitMq()

		/**
		 * SETUP CHANNEL
		 */
		rmq.SetUpChannel(RMQ_EXCHANGETYPE, true, RMQ_RENEWALEXCHANGE, true, RMQ_RENEWALQUEUE)

		/**
		 * Looping schedule
		 */
		timeDuration := time.Duration(1)

		for {
			timeNow := time.Now().Format("15:04")

			scheduleRepo := repository.NewScheduleRepository(db)
			scheduleService := services.NewScheduleService(scheduleRepo)

			if scheduleService.GetUnlocked(ACT_RENEWAL, timeNow) {

				scheduleService.UpdateSchedule(false, ACT_RENEWAL)

				go func() {
					populateRenewal(db, rmq)
				}()
			}

			if scheduleService.GetLocked(ACT_RENEWAL, timeNow) {
				scheduleService.UpdateSchedule(true, ACT_RENEWAL)

				/**
				** Purge queue retry if populate renewal start
				**/
				p := rabbit.NewRabbitMQ()
				p.Purge(RMQ_RETRYINSUFFQUEUE)
			}

			time.Sleep(timeDuration * time.Minute)

		}
	},
}

var publisherRetryFpCmd = &cobra.Command{
	Use:   "pub_retry_fp",
	Short: "Publisher Retry Firstpush CLI",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		/**
		 * connect pgsql
		 */
		db, err := connectPgsql()
		if err != nil {
			panic(err)
		}

		/**
		 * connect rabbitmq
		 */
		rmq := connectRabbitMq()

		/**
		 * SETUP CHANNEL
		 */
		rmq.SetUpChannel(RMQ_EXCHANGETYPE, true, RMQ_RETRYFPEXCHANGE, true, RMQ_RETRYFPQUEUE)

		/**
		 * Looping schedule
		 */
		timeDuration := time.Duration(1)

		for {

			/**
			** Populate retry if queue message is zero or 0
			**/
			p := rabbit.NewRabbitMQ()

			q, err := p.Queue(RMQ_RETRYFPQUEUE)
			if err != nil {
				log.Println(err)
			}

			var res *entity.RabbitMQResponse
			json.Unmarshal(q, &res)

			// if queue is empty
			if !res.IsRunning() {
				go func() {
					populateRetryFp(db, rmq)
				}()
			}

			time.Sleep(timeDuration * time.Minute)

		}

	},
}

var publisherRetryDpCmd = &cobra.Command{
	Use:   "pub_retry_dp",
	Short: "Publisher Retry Dailypush CLI",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		/**
		 * connect pgsql
		 */
		db, err := connectPgsql()
		if err != nil {
			panic(err)
		}

		/**
		 * connect rabbitmq
		 */
		rmq := connectRabbitMq()

		/**
		 * SETUP CHANNEL
		 */
		rmq.SetUpChannel(RMQ_EXCHANGETYPE, true, RMQ_RETRYDPEXCHANGE, true, RMQ_RETRYDPQUEUE)

		/**
		 * Looping schedule
		 */
		timeDuration := time.Duration(1)

		for {

			/**
			** Populate retry if queue message is zero or 0
			**/
			p := rabbit.NewRabbitMQ()

			q, err := p.Queue(RMQ_RETRYDPQUEUE)
			if err != nil {
				log.Println(err)
			}

			var res *entity.RabbitMQResponse
			json.Unmarshal(q, &res)

			// if queue is empty
			if !res.IsRunning() {
				go func() {
					populateRetryDp(db, rmq)
				}()
			}

			time.Sleep(timeDuration * time.Minute)

		}

	},
}

var publisherRetryInsuffCmd = &cobra.Command{
	Use:   "pub_retry_insuff",
	Short: "Publisher Retry Insuff CLI",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		/**
		 * connect pgsql
		 */
		db, err := connectPgsql()
		if err != nil {
			panic(err)
		}

		/**
		 * connect rabbitmq
		 */
		rmq := connectRabbitMq()

		/**
		 * SETUP CHANNEL
		 */
		rmq.SetUpChannel(RMQ_EXCHANGETYPE, true, RMQ_RETRYINSUFFEXCHANGE, true, RMQ_RETRYINSUFFQUEUE)

		/**
		 * Looping schedule
		 */
		timeDuration := time.Duration(1)

		for {
			timeNow := time.Now().Format("15:04")

			scheduleRepo := repository.NewScheduleRepository(db)
			scheduleService := services.NewScheduleService(scheduleRepo)

			if scheduleService.GetUnlocked(ACT_RETRY_INSUFF, timeNow) {

				scheduleService.UpdateSchedule(false, ACT_RETRY_INSUFF)

				go func() {
					populateRetryInsuff(db, rmq)
				}()
			}

			if scheduleService.GetLocked(ACT_RETRY_INSUFF, timeNow) {
				scheduleService.UpdateSchedule(true, ACT_RETRY_INSUFF)
			}

			time.Sleep(timeDuration * time.Minute)

		}

	},
}

var publisherCSVCmd = &cobra.Command{
	Use:   "pub_csv",
	Short: "CSV CLI",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		/**
		 * connect pgsql
		 */
		db, err := connectPgsql()
		if err != nil {
			panic(err)
		}

		/**
		 * Looping schedule
		 */
		timeDuration := time.Duration(1)

		for {
			timeNow := time.Now().Format("15:04")

			scheduleRepo := repository.NewScheduleRepository(db)
			scheduleService := services.NewScheduleService(scheduleRepo)

			if scheduleService.GetUnlocked(ACT_CSV, timeNow) {

				scheduleService.UpdateSchedule(false, ACT_CSV)

				go func() {
					populateCSV(db)
				}()
			}

			if scheduleService.GetLocked(ACT_CSV, timeNow) {
				scheduleService.UpdateSchedule(true, ACT_CSV)
			}

			time.Sleep(timeDuration * time.Minute)
		}

	},
}

var publisherUploadCSVCmd = &cobra.Command{
	Use:   "pub_upload_csv",
	Short: "Upload CSV CLI",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		/**
		 * connect pgsql
		 */
		db, err := connectPgsql()
		if err != nil {
			panic(err)
		}

		/**
		 * Looping schedule
		 */
		timeDuration := time.Duration(1)

		for {
			timeNow := time.Now().Format("15:04")

			scheduleRepo := repository.NewScheduleRepository(db)
			scheduleService := services.NewScheduleService(scheduleRepo)

			if scheduleService.GetUnlocked(ACT_UPLOAD_CSV, timeNow) {

				scheduleService.UpdateSchedule(false, ACT_UPLOAD_CSV)

				go func() {
					uploadCSV()
				}()
			}

			if scheduleService.GetLocked(ACT_UPLOAD_CSV, timeNow) {
				scheduleService.UpdateSchedule(true, ACT_UPLOAD_CSV)
			}

			time.Sleep(timeDuration * time.Minute)
		}

	},
}

func populateRenewal(db *sql.DB, queue rmqp.AMQP) {
	subscriptionRepo := repository.NewSubscriptionRepository(db)
	subscriptionService := services.NewSubscriptionService(subscriptionRepo)

	subs := subscriptionService.RenewalSubscription()
	for _, s := range *subs {
		var sub entity.Subscription

		sub.ID = s.ID
		sub.ServiceID = s.ServiceID
		sub.Msisdn = s.Msisdn
		sub.Channel = s.Channel
		sub.Adnet = s.Adnet
		sub.LatestKeyword = s.LatestKeyword
		sub.LatestSubject = s.LatestSubject
		sub.LatestPIN = s.LatestPIN
		sub.LatestPayload = s.LatestPayload
		sub.IpAddress = s.IpAddress
		sub.AffSub = s.AffSub
		sub.CampKeyword = s.CampKeyword
		sub.CampSubKeyword = s.CampSubKeyword
		sub.CreatedAt = s.CreatedAt

		json, _ := json.Marshal(sub)

		queue.IntegratePublish(RMQ_RENEWALEXCHANGE, RMQ_RENEWALQUEUE, RMQ_DATATYPE, "", string(json))

		time.Sleep(100 * time.Microsecond)
	}
}

func populateRetryFp(db *sql.DB, queue rmqp.AMQP) {
	subscriptionRepo := repository.NewSubscriptionRepository(db)
	subscriptionService := services.NewSubscriptionService(subscriptionRepo)

	subs := subscriptionService.RetryFpSubscription()

	for _, s := range *subs {
		var sub entity.Subscription

		sub.ID = s.ID
		sub.ServiceID = s.ServiceID
		sub.Msisdn = s.Msisdn
		sub.Channel = s.Channel
		sub.Adnet = s.Adnet
		sub.LatestKeyword = s.LatestKeyword
		sub.LatestSubject = s.LatestSubject
		sub.LatestPIN = s.LatestPIN
		sub.LatestPayload = s.LatestPayload
		sub.IpAddress = s.IpAddress
		sub.AffSub = s.AffSub
		sub.CampKeyword = s.CampKeyword
		sub.CampSubKeyword = s.CampSubKeyword
		sub.RetryAt = s.RetryAt
		sub.CreatedAt = s.CreatedAt

		json, _ := json.Marshal(sub)
		queue.IntegratePublish(RMQ_RETRYFPEXCHANGE, RMQ_RETRYFPQUEUE, RMQ_DATATYPE, "", string(json))

		time.Sleep(100 * time.Microsecond)
	}
}

func populateRetryDp(db *sql.DB, queue rmqp.AMQP) {
	subscriptionRepo := repository.NewSubscriptionRepository(db)
	subscriptionService := services.NewSubscriptionService(subscriptionRepo)

	subs := subscriptionService.RetryDpSubscription()

	for _, s := range *subs {
		var sub entity.Subscription

		sub.ID = s.ID
		sub.ServiceID = s.ServiceID
		sub.Msisdn = s.Msisdn
		sub.Channel = s.Channel
		sub.Adnet = s.Adnet
		sub.LatestKeyword = s.LatestKeyword
		sub.LatestSubject = s.LatestSubject
		sub.LatestPIN = s.LatestPIN
		sub.LatestPayload = s.LatestPayload
		sub.IpAddress = s.IpAddress
		sub.AffSub = s.AffSub
		sub.CampKeyword = s.CampKeyword
		sub.CampSubKeyword = s.CampSubKeyword
		sub.RetryAt = s.RetryAt
		sub.CreatedAt = s.CreatedAt

		json, _ := json.Marshal(sub)
		queue.IntegratePublish(RMQ_RETRYDPEXCHANGE, RMQ_RETRYDPQUEUE, RMQ_DATATYPE, "", string(json))

		time.Sleep(100 * time.Microsecond)
	}
}

func populateRetryInsuff(db *sql.DB, queue rmqp.AMQP) {
	subscriptionRepo := repository.NewSubscriptionRepository(db)
	subscriptionService := services.NewSubscriptionService(subscriptionRepo)

	subs := subscriptionService.RetryInsuffSubscription()

	for _, s := range *subs {
		var sub entity.Subscription

		sub.ID = s.ID
		sub.ServiceID = s.ServiceID
		sub.Msisdn = s.Msisdn
		sub.Channel = s.Channel
		sub.Adnet = s.Adnet
		sub.LatestKeyword = s.LatestKeyword
		sub.LatestSubject = s.LatestSubject
		sub.LatestPIN = s.LatestPIN
		sub.LatestPayload = s.LatestPayload
		sub.IpAddress = s.IpAddress
		sub.AffSub = s.AffSub
		sub.CampKeyword = s.CampKeyword
		sub.CampSubKeyword = s.CampSubKeyword
		sub.RetryAt = s.RetryAt
		sub.CreatedAt = s.CreatedAt

		json, _ := json.Marshal(sub)
		queue.IntegratePublish(RMQ_RETRYINSUFFEXCHANGE, RMQ_RETRYINSUFFQUEUE, RMQ_DATATYPE, "", string(json))

		time.Sleep(100 * time.Microsecond)
	}
}

func populateCSV(db *sql.DB) {

	fileNameSubs := "/logs/csv/subscriptions_id_telkomsel_cloudplay.csv"
	fileNameTrans := "/logs/csv/transactions_id_telkomsel_cloudplay.csv"

	subscriptionRepo := repository.NewSubscriptionRepository(db)
	subscriptionService := services.NewSubscriptionService(subscriptionRepo)
	transactionRepo := repository.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepo)

	subRecords, err := subscriptionService.SelectSubcriptionToCSV()
	if err != nil {
		log.Fatalf("error load table subscriptions: %s", err)
	}

	// delete file sub csv
	os.Remove(fileNameSubs)

	subCsv, err := os.Create(fileNameSubs)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	defer subCsv.Close()
	subW := csv.NewWriter(subCsv)
	defer subW.Flush()

	subsHeaders := []string{
		"country", "operator", "service", "source", "msisdn",
		"status", "cycle", "adnet", "revenue", "subs_date",
		"renewal_date", "freemium_end_date", "unsubs_from", "unsubs_date",
		"service_price", "currency", "profile_status", "publisher",
		"trxid", "pixel", "handset", "browser", "attempt_charging",
		"success_billing",
	}
	subW.Write(subsHeaders)

	var subsData [][]string
	for _, r := range *subRecords {
		row := []string{
			r.Country, r.Operator, r.Service, r.Source, r.Msisdn,
			r.LatestSubject, r.Cycle, r.Adnet, r.Revenue, r.SubsDate.String,
			r.RenewalDate.String, r.FreemiumEndDate, r.UnsubsFrom, r.UnsubsDate.String,
			r.ServicePrice, r.Currency, r.ProfileStatus, r.Publisher,
			r.Trxid, r.Pixel, r.Handset, r.Browser, r.AttemptCharging,
			r.SuccessBilling,
		}
		subsData = append(subsData, row)
	}

	err = subW.WriteAll(subsData) // calls Flush internally
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(20 * time.Second)

	transRecords, err := transactionService.SelectTransactionToCSV()
	if err != nil {
		log.Fatalf("error load table transactions: %s", err)
	}

	// delete file trans csv
	os.Remove(fileNameTrans)

	transCsv, err := os.Create(fileNameTrans)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	defer transCsv.Close()
	transW := csv.NewWriter(transCsv)
	defer transW.Flush()

	transHeaders := []string{
		"country", "operator", "service", "source", "msisdn",
		"event", "event_date", "cycle", "revenue", "charge_date",
		"currency", "publisher", "handset",
		"browser", "trxid", "telco_api_url", "telco_api_response",
		"sms_content", "status_sms",
	}
	transW.Write(transHeaders)

	var transData [][]string
	for _, r := range *transRecords {
		row := []string{
			r.Country, r.Operator, r.Service, r.Source, r.Msisdn,
			r.Event, r.EventDate.String, r.Cycle, r.Revenue, r.ChargeDate.String,
			r.Currency, r.Publisher, r.Handset,
			r.Browser, r.TrxId, r.TelcoApiUrl, r.TelcoApiResponse,
			r.SmsContent, r.StatusSms,
		}
		transData = append(transData, row)
	}

	err = transW.WriteAll(transData) // calls Flush internally
	if err != nil {
		log.Println(err.Error())
	}

}

func compressCSV(file1, file2 string) {

	fileNameSubsCompress := "/logs/csv/subscriptions_id_telkomsel_cloudplay.gz"
	fileNameTransCompress := "/logs/csv/transactions_id_telkomsel_cloudplay.gz"

	f1, err := os.Open(file1)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer f1.Close()

	f2, err := os.Open(file2)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer f2.Close()

	// Create a new gzip writer
	gzipSubsWriter, err := os.Create(fileNameSubsCompress)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer gzipSubsWriter.Close()

	// Create a new gzip writer
	gzipTransWriter, err := os.Create(fileNameTransCompress)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer gzipTransWriter.Close()

	fmt.Println("archive file created successfully....")

	zipSubsWriter := gzip.NewWriter(gzipSubsWriter)
	defer zipSubsWriter.Close()

	zipTransWriter := gzip.NewWriter(gzipTransWriter)
	defer zipTransWriter.Close()

	_, err = io.Copy(zipSubsWriter, f1)
	if err != nil {
		log.Fatal(err.Error())
	}
	// Close the gzip writer
	zipSubsWriter.Close()

	_, err = io.Copy(zipTransWriter, f2)
	if err != nil {
		log.Fatal(err.Error())
	}
	// Close the gzip writer
	zipTransWriter.Close()
}

func uploadCSV() {
	/**
	 * SETUP LOG
	 */
	logger := logger.NewLogger()

	fileNameSubs := "/logs/csv/subscriptions_id_telkomsel_cloudplay.csv"
	fileNameTrans := "/logs/csv/transactions_id_telkomsel_cloudplay.csv"

	fileNameSubsCompress := "/logs/csv/subscriptions_id_telkomsel_cloudplay.gz"
	fileNameTransCompress := "/logs/csv/transactions_id_telkomsel_cloudplay.gz"

	// compress in new file
	compressCSV(fileNameSubs, fileNameTrans)

	arp := arpu.NewArpu(logger)

	// upload file csv
	arp.UploadCSV(ARPU_URL_SUB, fileNameSubsCompress)

	time.Sleep(10 * time.Second)
	// upload file csv
	arp.UploadCSV(ARPU_URL_TRANS, fileNameTransCompress)
}
