package cmd

import (
	"log"

	"github.com/idprm/go-linkit-tsel/internal/app"
	"github.com/idprm/go-linkit-tsel/internal/logger"
	"github.com/spf13/cobra"
)

var listenerCmd = &cobra.Command{
	Use:   "listener",
	Short: "Webserver CLI",
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
		rmq.SetUpChannel(RMQ_EXCHANGETYPE, true, RMQ_MOEXCHANGE, true, RMQ_MOQUEUE)

		router := app.StartApplication(db, rmq, rds, logger)
		log.Fatal(router.Listen(":" + APP_PORT))
	},
}
