package cmd

import (
	"log"

	"github.com/idprm/go-linkit-tsel/src/app"
	"github.com/idprm/go-linkit-tsel/src/config"
	"github.com/idprm/go-linkit-tsel/src/datasource/pgsql/db"
	"github.com/idprm/go-linkit-tsel/src/datasource/rabbitmq"
	"github.com/idprm/go-linkit-tsel/src/datasource/redis/rdb"
	"github.com/idprm/go-linkit-tsel/src/logger"
	"github.com/spf13/cobra"
)

var listenerCmd = &cobra.Command{
	Use:   "listener",
	Short: "Webserver CLI",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		/**
		 * LOAD CONFIG
		 */
		cfg, err := config.LoadSecret("secret.yaml")
		if err != nil {
			panic(err)
		}

		/**
		 * SETUP PGSQL
		 */
		db := db.InitDB(cfg)

		/**
		 * SETUP REDIS
		 */
		rdb := rdb.InitRedis(cfg)

		/**
		 * SETUP LOG
		 */
		logger := logger.NewLogger(cfg)

		/**
		 * SETUP RMQ
		 */
		queue := rabbitmq.InitQueue(cfg)

		/**
		 * SETUP CHANNEL
		 */
		queue.SetUpChannel(RMQ_EXCHANGETYPE, true, RMQ_MOEXCHANGE, true, RMQ_MOQUEUE)

		router := app.StartApplication(cfg, db, rdb, logger, queue)
		log.Fatal(router.Listen(":" + cfg.App.Port))
	},
}
