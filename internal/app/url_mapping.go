package app

import (
	"database/sql"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html/v2"
	"github.com/idprm/go-linkit-tsel/internal/domain/repository"
	"github.com/idprm/go-linkit-tsel/internal/handler"
	"github.com/idprm/go-linkit-tsel/internal/logger"
	"github.com/idprm/go-linkit-tsel/internal/services"
	"github.com/idprm/go-linkit-tsel/internal/utils"
	"github.com/redis/go-redis/v9"
	"github.com/wiliehidayat87/rmqp"
)

var (
	PUBLIC_PATH string = utils.GetEnv("PUBLIC_PATH")
)

func mapUrls(db *sql.DB, rmpq rmqp.AMQP, rdb *redis.Client, logger *logger.Logger) *fiber.App {

	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	engine := html.New(path+"/views", ".html")
	/**
	 * Init Fiber
	 */
	router := fiber.New(fiber.Config{
		Views: engine,
	})

	/**
	 * Initialize default config
	 */
	router.Use(cors.New())

	router.Static("/static", path+"/"+PUBLIC_PATH)

	serviceRepo := repository.NewServiceRepository(db)
	serviceService := services.NewServiceService(serviceRepo)

	verifyRepo := repository.NewVerifyRepository(rdb)
	verifyService := services.NewVerifyService(verifyRepo)

	subscriptionRepo := repository.NewSubscriptionRepository(db)
	subscriptionService := services.NewSubscriptionService(subscriptionRepo)

	transactionRepo := repository.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepo)

	incomingHandler := handler.NewIncomingHandler(logger, rmpq, serviceService, verifyService, subscriptionService, transactionService)

	/**
	 * Routes Landing Page SUB & UNSUB
	 */
	router.Get("camp/:service", incomingHandler.CampaignDirect)
	router.Get("camptool", incomingHandler.CampaignTool)

	router.Post("cloudplay", incomingHandler.OptIn)
	router.Post("galays", incomingHandler.OptIn)

	router.Get("cloudplay", incomingHandler.CloudPlaySubPage)
	router.Get("cloudplay/camp", incomingHandler.CloudPlayCampaign)
	router.Get("cloudplay/campbill", incomingHandler.CloudPlayCampaignBillable)

	router.Get("cloudplay1", incomingHandler.CloudPlaySub1Page)
	router.Get("cloudplay1/camp", incomingHandler.CloudPlaySub1CampaignPage)

	router.Get("cloudplay2", incomingHandler.CloudPlaySub2Page)
	router.Get("cloudplay2/camp", incomingHandler.CloudPlaySub2CampaignPage)

	router.Get("cloudplay3", incomingHandler.CloudPlaySub3Page)
	router.Get("cloudplay3/camp", incomingHandler.CloudPlaySub3CampaignPage)

	router.Get("cloudplay4", incomingHandler.CloudPlaySub4Page)
	router.Get("cloudplay4/camp", incomingHandler.CloudPlaySub4CampaignPage)

	router.Get("cloudplay/unsub", incomingHandler.CloudPlayUnsubPage)

	router.Get("cbtsel", incomingHandler.CallbackUrl)

	router.Get("cloudplay/camptool", incomingHandler.CampaignTool)
	router.Get("cloudplay/camptooldynamic", incomingHandler.CampaignToolDynamic)

	router.Get("galays", incomingHandler.GalaysSubPage)
	router.Get("galays/camp", incomingHandler.GalaysCampaign)
	router.Get("galays/campbill", incomingHandler.GalaysCampaignBillable)

	router.Get("galays1", incomingHandler.GalaysSub1Page)
	router.Get("galays1/camp", incomingHandler.GalaysSub1CampaignPage)

	router.Get("galays/camptool", incomingHandler.CampaignTool)
	router.Get("galays/camptooldynamic", incomingHandler.CampaignToolDynamic)

	router.Get("success", incomingHandler.CallbackUrl)
	router.Get("cancel", incomingHandler.CallbackUrl)

	/**``
	 * Routes Another CP
	 */
	router.Get("cpsam", incomingHandler.CloudPlaySub1Page)
	router.Get("cpsam/camp", incomingHandler.CloudPlaySub1CampaignPage)

	/**
	 * Routes MO & DR
	 */
	router.Get("notif/mo", incomingHandler.MessageOriginated)
	router.Get("mo", incomingHandler.MessageOriginated)

	/**
	 * Routes Report
	 */
	router.Get("report/status", incomingHandler.SelectStatus)
	router.Get("report/statusdetail", incomingHandler.SelectStatusDetail)
	router.Get("report/adnet", incomingHandler.SelectAdnet)
	router.Get("report/daily", incomingHandler.ReportDaily)
	router.Post("report/arpu", incomingHandler.AveragePerUser)

	return router
}
