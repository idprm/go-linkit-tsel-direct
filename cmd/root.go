package cmd

import (
	"database/sql"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
	"github.com/wiliehidayat87/rmqp"
)

var (
	APP_HOST       string = getEnv("APP_HOST")
	APP_PORT       string = getEnv("APP_PORT")
	APP_TZ         string = getEnv("APP_TZ")
	APP_URL        string = getEnv("APP_URL")
	URI_POSTGRES   string = getEnv("URI_POSTGRES")
	URI_REDIS      string = getEnv("URI_REDIS")
	URI_AMQP       string = getEnv("URI_AMQP")
	RMQ_HOST       string = getEnv("RMQ_HOST")
	RMQ_USER       string = getEnv("RMQ_USER")
	RMQ_PASS       string = getEnv("RMQ_PASS")
	RMQ_PORT       string = getEnv("RMQ_PORT")
	RMQ_URL        string = getEnv("RMQ_URL")
	ARPU_URL_SUB   string = getEnv("ARPU_URL_SUB")
	ARPU_URL_TRANS string = getEnv("ARPU_URL_TRANS")
	LOG_PATH       string = getEnv("LOG_PATH")
)

const (
	RMQ_EXCHANGETYPE        string = "direct"
	RMQ_DATATYPE            string = "application/json"
	RMQ_MOEXCHANGE          string = "E_MO"
	RMQ_MOQUEUE             string = "Q_MO"
	RMQ_RENEWALEXCHANGE     string = "E_RENEWAL"
	RMQ_RENEWALQUEUE        string = "Q_RENEWAL"
	RMQ_RETRYFPEXCHANGE     string = "E_RETRY_FP"
	RMQ_RETRYFPQUEUE        string = "Q_RETRY_FP"
	RMQ_RETRYDPEXCHANGE     string = "E_RETRY_DP"
	RMQ_RETRYDPQUEUE        string = "Q_RETRY_DP"
	RMQ_RETRYINSUFFEXCHANGE string = "E_RETRY_INSUFF"
	RMQ_RETRYINSUFFQUEUE    string = "Q_RETRY_INSUFF"
	RMQ_NOTIFEXCHANGE       string = "E_NOTIF"
	RMQ_NOTIFQUEUE          string = "Q_NOTIF"
	RMQ_POSTBACKMOEXCHANGE  string = "E_POSTBACK_MO"
	RMQ_POSTBACKMOQUEUE     string = "Q_POSTBACK_MO"
	RMQ_POSTBACKMTEXCHANGE  string = "E_POSTBACK_MT"
	RMQ_POSTBACKMTQUEUE     string = "Q_POSTBACK_MT"
	MT_FIRSTPUSH            string = "MT_FIRSTPUSH"
	ACT_RENEWAL             string = "RENEWAL"
	ACT_RETRY_FP            string = "RETRY_FP"
	ACT_RETRY_DP            string = "RETRY_DP"
	ACT_RETRY_INSUFF        string = "RETRY_INSUFF"
	ACT_CSV                 string = "CSV"
)

var (
	rootCmd = &cobra.Command{
		Use:   "cobra-cli",
		Short: "A generator for Cobra based Applications",
		Long:  `Cobra is a CLI library for Go that empowers applications.`,
	}
)

func init() {
	// setup timezone
	loc, _ := time.LoadLocation(APP_TZ)
	time.Local = loc
	/**
	 * WEBSERVER SERVICE
	 */
	rootCmd.AddCommand(listenerCmd)

	/**
	 * RABBITMQ SERVICE
	 */
	rootCmd.AddCommand(consumerMOCmd)
	rootCmd.AddCommand(consumerRenewalCmd)
	rootCmd.AddCommand(consumerRetryFpCmd)
	rootCmd.AddCommand(consumerRetryDpCmd)
	rootCmd.AddCommand(consumerRetryInsuffCmd)
	rootCmd.AddCommand(consumerNotifCmd)
	rootCmd.AddCommand(consumerPostbackMOCmd)
	rootCmd.AddCommand(consumerPostbackMTCmd)

	rootCmd.AddCommand(publisherRenewalCmd)
	rootCmd.AddCommand(publisherRetryFpCmd)
	rootCmd.AddCommand(publisherRetryDpCmd)
	rootCmd.AddCommand(publisherRetryInsuffCmd)
	rootCmd.AddCommand(publisherCSVCmd)

}

func Execute() error {
	return rootCmd.Execute()
}

func getEnv(key string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		log.Panicf("Error %v", key)
	}
	return value
}

// Connect to postgresql
func connectPgsql() (*sql.DB, error) {
	db, err := sql.Open("postgres", URI_POSTGRES)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Connect to redis
func connectRedis() (*redis.Client, error) {
	opts, err := redis.ParseURL(URI_REDIS)
	if err != nil {
		return nil, err
	}
	return redis.NewClient(opts), nil
}

// Connect to rabbitmq
func connectRabbitMq() rmqp.AMQP {
	var rb rmqp.AMQP
	port, _ := strconv.Atoi(RMQ_PORT)
	rb.SetAmqpURL(RMQ_HOST, port, RMQ_USER, RMQ_PASS)
	rb.SetUpConnectionAmqp()
	return rb
}
