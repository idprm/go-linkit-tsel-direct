package cmd

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/idprm/go-linkit-tsel/internal/logger"
	"github.com/spf13/cobra"
)

var consumerMOCmd = &cobra.Command{
	Use:   "mo",
	Short: "Consumer MO Service CLI",
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
		rmq.SetUpChannel(RMQ_EXCHANGETYPE, true, RMQ_NOTIFEXCHANGE, true, RMQ_NOTIFQUEUE)
		rmq.SetUpChannel(RMQ_EXCHANGETYPE, true, RMQ_POSTBACKMOEXCHANGE, true, RMQ_POSTBACKMOQUEUE)

		messagesData := rmq.Subscribe(1, false, RMQ_MOQUEUE, RMQ_MOEXCHANGE, RMQ_MOQUEUE)

		// Initial sync waiting group
		var wg sync.WaitGroup

		// Loop forever listening incoming data
		forever := make(chan bool)

		processor := NewProcessor(db, rmq, rds, logger)

		// Set into goroutine this listener
		go func() {

			// Loop every incoming data
			for d := range messagesData {

				wg.Add(1)
				processor.MO(&wg, d.Body)
				wg.Wait()

				// Manual consume queue
				d.Ack(false)

			}

		}()

		fmt.Println("[*] Waiting for data...")

		<-forever
	},
}

var consumerRenewalCmd = &cobra.Command{
	Use:   "renewal",
	Short: "Consumer Renewal Service CLI",
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
		rmq.SetUpChannel(RMQ_EXCHANGETYPE, true, RMQ_RENEWALEXCHANGE, true, RMQ_RENEWALQUEUE)
		rmq.SetUpChannel(RMQ_EXCHANGETYPE, true, RMQ_NOTIFEXCHANGE, true, RMQ_NOTIFQUEUE)
		rmq.SetUpChannel(RMQ_EXCHANGETYPE, true, RMQ_POSTBACKMTEXCHANGE, true, RMQ_POSTBACKMTQUEUE)

		messagesData := rmq.Subscribe(1, false, RMQ_RENEWALQUEUE, RMQ_RENEWALEXCHANGE, RMQ_RENEWALQUEUE)

		// Initial sync waiting group
		var wg sync.WaitGroup

		// Loop forever listening incoming data
		forever := make(chan bool)

		processor := NewProcessor(db, rmq, rds, logger)

		// Set into goroutine this listener
		go func() {

			// Loop every incoming data
			for d := range messagesData {

				wg.Add(1)
				processor.Renewal(&wg, d.Body)
				wg.Wait()

				// Manual consume queue
				d.Ack(false)

			}

		}()

		fmt.Println("[*] Waiting for data...")

		<-forever
	},
}

var consumerRetryFpCmd = &cobra.Command{
	Use:   "retry_fp",
	Short: "Consumer Retry Firstpush Service CLI",
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
		rmq.SetUpChannel(RMQ_EXCHANGETYPE, true, RMQ_RETRYFPEXCHANGE, true, RMQ_RETRYFPQUEUE)
		rmq.SetUpChannel(RMQ_EXCHANGETYPE, true, RMQ_NOTIFEXCHANGE, true, RMQ_NOTIFQUEUE)
		rmq.SetUpChannel(RMQ_EXCHANGETYPE, true, RMQ_POSTBACKMTEXCHANGE, true, RMQ_POSTBACKMTQUEUE)

		messagesData := rmq.Subscribe(1, false, RMQ_RETRYFPQUEUE, RMQ_RETRYFPEXCHANGE, RMQ_RETRYFPQUEUE)

		// Initial sync waiting group
		var wg sync.WaitGroup

		// Loop forever listening incoming data
		forever := make(chan bool)

		processor := NewProcessor(db, rmq, rds, logger)

		// Set into goroutine this listener
		go func() {

			// Loop every incoming data
			for d := range messagesData {

				wg.Add(1)
				processor.RetryFp(&wg, d.Body)
				wg.Wait()

				// Manual consume queue
				d.Ack(false)

			}

		}()

		fmt.Println("[*] Waiting for data...")

		<-forever
	},
}

var consumerRetryDpCmd = &cobra.Command{
	Use:   "retry_dp",
	Short: "Consumer Retry Dailypush Service CLI",
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
		rmq.SetUpChannel(RMQ_EXCHANGETYPE, true, RMQ_RETRYDPEXCHANGE, true, RMQ_RETRYDPQUEUE)
		rmq.SetUpChannel(RMQ_EXCHANGETYPE, true, RMQ_NOTIFEXCHANGE, true, RMQ_NOTIFQUEUE)
		rmq.SetUpChannel(RMQ_EXCHANGETYPE, true, RMQ_POSTBACKMTEXCHANGE, true, RMQ_POSTBACKMTQUEUE)

		messagesData := rmq.Subscribe(1, false, RMQ_RETRYDPQUEUE, RMQ_RETRYDPEXCHANGE, RMQ_RETRYDPQUEUE)

		// Initial sync waiting group
		var wg sync.WaitGroup

		// Loop forever listening incoming data
		forever := make(chan bool)

		processor := NewProcessor(db, rmq, rds, logger)

		// Set into goroutine this listener
		go func() {

			// Loop every incoming data
			for d := range messagesData {

				wg.Add(1)
				processor.RetryDp(&wg, d.Body)
				wg.Wait()

				// Manual consume queue
				d.Ack(false)

			}

		}()

		fmt.Println("[*] Waiting for data...")

		<-forever
	},
}

var consumerRetryInsuffCmd = &cobra.Command{
	Use:   "retry_insuff",
	Short: "Consumer Retry Insuff Service CLI",
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
		rmq.SetUpChannel(RMQ_EXCHANGETYPE, true, RMQ_RETRYINSUFFEXCHANGE, true, RMQ_RETRYINSUFFQUEUE)
		rmq.SetUpChannel(RMQ_EXCHANGETYPE, true, RMQ_NOTIFEXCHANGE, true, RMQ_NOTIFQUEUE)
		rmq.SetUpChannel(RMQ_EXCHANGETYPE, true, RMQ_POSTBACKMTEXCHANGE, true, RMQ_POSTBACKMTQUEUE)

		messagesData := rmq.Subscribe(1, false, RMQ_RETRYINSUFFQUEUE, RMQ_RETRYINSUFFEXCHANGE, RMQ_RETRYINSUFFQUEUE)

		// Initial sync waiting group
		var wg sync.WaitGroup

		// Loop forever listening incoming data
		forever := make(chan bool)

		processor := NewProcessor(db, rmq, rds, logger)

		// Set into goroutine this listener
		go func() {

			// Loop every incoming data
			for d := range messagesData {

				wg.Add(1)
				processor.RetryInsuff(&wg, d.Body)
				wg.Wait()

				// Manual consume queue
				d.Ack(false)

			}

		}()

		fmt.Println("[*] Waiting for data...")

		<-forever
	},
}

var consumerRetryInsuffBill0Cmd = &cobra.Command{
	Use:   "retry_insuff_bill0",
	Short: "Consumer Retry Insuff Bill 0 Service CLI",
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
		rmq.SetUpChannel(RMQ_EXCHANGETYPE, true, RMQ_RETRYINSUFFBILL0EXCHANGE, true, RMQ_RETRYINSUFFBILL0QUEUE)
		rmq.SetUpChannel(RMQ_EXCHANGETYPE, true, RMQ_NOTIFEXCHANGE, true, RMQ_NOTIFQUEUE)
		rmq.SetUpChannel(RMQ_EXCHANGETYPE, true, RMQ_POSTBACKMTEXCHANGE, true, RMQ_POSTBACKMTQUEUE)

		messagesData := rmq.Subscribe(1, false, RMQ_RETRYINSUFFBILL0QUEUE, RMQ_RETRYINSUFFBILL0EXCHANGE, RMQ_RETRYINSUFFBILL0QUEUE)

		// Initial sync waiting group
		var wg sync.WaitGroup

		// Loop forever listening incoming data
		forever := make(chan bool)

		processor := NewProcessor(db, rmq, rds, logger)

		// Set into goroutine this listener
		go func() {

			// Loop every incoming data
			for d := range messagesData {

				wg.Add(1)
				processor.RetryInsuff(&wg, d.Body)
				wg.Wait()

				// Manual consume queue
				d.Ack(false)

			}

		}()

		fmt.Println("[*] Waiting for data...")

		<-forever
	},
}

var consumerRetryInsuffBill1Cmd = &cobra.Command{
	Use:   "retry_insuff_bill1",
	Short: "Consumer Retry Insuff Bill 1 Service CLI",
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
		rmq.SetUpChannel(RMQ_EXCHANGETYPE, true, RMQ_RETRYINSUFFBILL1EXCHANGE, true, RMQ_RETRYINSUFFBILL1QUEUE)
		rmq.SetUpChannel(RMQ_EXCHANGETYPE, true, RMQ_NOTIFEXCHANGE, true, RMQ_NOTIFQUEUE)
		rmq.SetUpChannel(RMQ_EXCHANGETYPE, true, RMQ_POSTBACKMTEXCHANGE, true, RMQ_POSTBACKMTQUEUE)

		messagesData := rmq.Subscribe(1, false, RMQ_RETRYINSUFFBILL1QUEUE, RMQ_RETRYINSUFFBILL1EXCHANGE, RMQ_RETRYINSUFFBILL1QUEUE)

		// Initial sync waiting group
		var wg sync.WaitGroup

		// Loop forever listening incoming data
		forever := make(chan bool)

		processor := NewProcessor(db, rmq, rds, logger)

		// Set into goroutine this listener
		go func() {

			// Loop every incoming data
			for d := range messagesData {

				wg.Add(1)
				processor.RetryInsuff(&wg, d.Body)
				wg.Wait()

				// Manual consume queue
				d.Ack(false)

			}

		}()

		fmt.Println("[*] Waiting for data...")

		<-forever
	},
}

var consumerRetryInsuffBill2Cmd = &cobra.Command{
	Use:   "retry_insuff_bill2",
	Short: "Consumer Retry Insuff Bill 2 Service CLI",
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
		rmq.SetUpChannel(RMQ_EXCHANGETYPE, true, RMQ_RETRYINSUFFBILL2EXCHANGE, true, RMQ_RETRYINSUFFBILL2QUEUE)
		rmq.SetUpChannel(RMQ_EXCHANGETYPE, true, RMQ_NOTIFEXCHANGE, true, RMQ_NOTIFQUEUE)
		rmq.SetUpChannel(RMQ_EXCHANGETYPE, true, RMQ_POSTBACKMTEXCHANGE, true, RMQ_POSTBACKMTQUEUE)

		messagesData := rmq.Subscribe(1, false, RMQ_RETRYINSUFFBILL2QUEUE, RMQ_RETRYINSUFFBILL2EXCHANGE, RMQ_RETRYINSUFFBILL2QUEUE)

		// Initial sync waiting group
		var wg sync.WaitGroup

		// Loop forever listening incoming data
		forever := make(chan bool)

		processor := NewProcessor(db, rmq, rds, logger)

		// Set into goroutine this listener
		go func() {

			// Loop every incoming data
			for d := range messagesData {

				wg.Add(1)
				processor.RetryInsuff(&wg, d.Body)
				wg.Wait()

				// Manual consume queue
				d.Ack(false)

			}

		}()

		fmt.Println("[*] Waiting for data...")

		<-forever
	},
}

var consumerRetryInsuffBill3Cmd = &cobra.Command{
	Use:   "retry_insuff_bill3",
	Short: "Consumer Retry Insuff Bill 3 Service CLI",
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
		rmq.SetUpChannel(RMQ_EXCHANGETYPE, true, RMQ_RETRYINSUFFBILL3EXCHANGE, true, RMQ_RETRYINSUFFBILL3QUEUE)
		rmq.SetUpChannel(RMQ_EXCHANGETYPE, true, RMQ_NOTIFEXCHANGE, true, RMQ_NOTIFQUEUE)
		rmq.SetUpChannel(RMQ_EXCHANGETYPE, true, RMQ_POSTBACKMTEXCHANGE, true, RMQ_POSTBACKMTQUEUE)

		messagesData := rmq.Subscribe(1, false, RMQ_RETRYINSUFFBILL3QUEUE, RMQ_RETRYINSUFFBILL3EXCHANGE, RMQ_RETRYINSUFFBILL3QUEUE)

		// Initial sync waiting group
		var wg sync.WaitGroup

		// Loop forever listening incoming data
		forever := make(chan bool)

		processor := NewProcessor(db, rmq, rds, logger)

		// Set into goroutine this listener
		go func() {

			// Loop every incoming data
			for d := range messagesData {

				wg.Add(1)
				processor.RetryInsuff(&wg, d.Body)
				wg.Wait()

				// Manual consume queue
				d.Ack(false)

			}

		}()

		fmt.Println("[*] Waiting for data...")

		<-forever
	},
}

var consumerNotifCmd = &cobra.Command{
	Use:   "notif",
	Short: "Consumer Notif Service CLI",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

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
		rmq.SetUpChannel(RMQ_EXCHANGETYPE, true, RMQ_NOTIFEXCHANGE, true, RMQ_NOTIFQUEUE)

		messagesData := rmq.Subscribe(1, false, RMQ_NOTIFQUEUE, RMQ_NOTIFEXCHANGE, RMQ_NOTIFQUEUE)

		// Initial sync waiting group
		var wg sync.WaitGroup

		// Loop forever listening incoming data
		forever := make(chan bool)

		processor := NewProcessor(&sql.DB{}, rmq, rds, logger)

		// Set into goroutine this listener
		go func() {

			// Loop every incoming data
			for d := range messagesData {

				wg.Add(1)
				processor.Notif(&wg, d.Body)
				wg.Wait()

				// Manual consume queue
				d.Ack(false)
			}

		}()

		fmt.Println("[*] Waiting for data...")

		<-forever
	},
}

var consumerPostbackMOCmd = &cobra.Command{
	Use:   "postback_mo",
	Short: "Consumer Postback MO Service CLI",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

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
		rmq.SetUpChannel(RMQ_EXCHANGETYPE, true, RMQ_POSTBACKMOEXCHANGE, true, RMQ_POSTBACKMOQUEUE)

		messagesData := rmq.Subscribe(1, false, RMQ_POSTBACKMOQUEUE, RMQ_POSTBACKMOEXCHANGE, RMQ_POSTBACKMOQUEUE)

		// Initial sync waiting group
		var wg sync.WaitGroup

		// Loop forever listening incoming data
		forever := make(chan bool)

		processor := NewProcessor(&sql.DB{}, rmq, rds, logger)

		// Set into goroutine this listener
		go func() {

			// Loop every incoming data
			for d := range messagesData {

				wg.Add(1)
				processor.PostbackMO(&wg, d.Body)
				wg.Wait()

				// Manual consume queue
				d.Ack(false)
			}

		}()

		fmt.Println("[*] Waiting for data...")

		<-forever
	},
}

var consumerPostbackMTCmd = &cobra.Command{
	Use:   "postback_mt",
	Short: "Consumer Postback MT Service CLI",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

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
		rmq.SetUpChannel(RMQ_EXCHANGETYPE, true, RMQ_POSTBACKMTEXCHANGE, true, RMQ_POSTBACKMTQUEUE)

		messagesData := rmq.Subscribe(1, false, RMQ_POSTBACKMTQUEUE, RMQ_POSTBACKMTEXCHANGE, RMQ_POSTBACKMTQUEUE)

		// Initial sync waiting group
		var wg sync.WaitGroup

		// Loop forever listening incoming data
		forever := make(chan bool)

		processor := NewProcessor(&sql.DB{}, rmq, rds, logger)

		// Set into goroutine this listener
		go func() {

			// Loop every incoming data
			for d := range messagesData {

				wg.Add(1)
				processor.PostbackMT(&wg, d.Body)
				wg.Wait()

				// Manual consume queue
				d.Ack(false)
			}

		}()

		fmt.Println("[*] Waiting for data...")

		<-forever
	},
}
