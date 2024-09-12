package cmd

import (
	"log"

	"github.com/idprm/go-linkit-tsel/internal/app"
	"github.com/idprm/go-linkit-tsel/internal/logger"
	"github.com/spf13/cobra"
)

var listenerCmd = &cobra.Command{
	Use:   "listener",
	Short: "Listener Service CLI",
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
		 * connect redis
		 */
		rds, err := connectRedis()
		if err != nil {
			panic(err)
		}

		/**
		 * SETUP LOG
		 */
		logger := logger.NewLogger()

		/**
		 * SETUP CHANNEL
		 */
		rmq.SetUpChannel(RMQ_EXCHANGE_TYPE, true, RMQ_MO_EXCHANGE, true, RMQ_MO_QUEUE)
		rmq.SetUpChannel(RMQ_EXCHANGE_TYPE, true, RMQ_TRAFFIC_EXCHANGE, true, RMQ_TRAFFIC_QUEUE)

		router := app.StartApplication(db, rmq, rds, logger)
		log.Fatal(router.Listen(":" + APP_PORT))
	},
}
