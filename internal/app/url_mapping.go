package app

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
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
	LOG_PATH    string = utils.GetEnv("LOG_PATH")
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
	r := fiber.New(fiber.Config{
		Views: engine,
	})

	/**
	 * Initialize default config
	 */
	r.Use(cors.New())

	/**
	 * Access log on browser
	 */
	r.Use(LOG_PATH, filesystem.New(
		filesystem.Config{
			Root:         http.Dir(LOG_PATH),
			Browse:       true,
			Index:        "index.html",
			NotFoundFile: "404.html",
		},
	))

	r.Static("/static", path+"/"+PUBLIC_PATH)

	serviceRepo := repository.NewServiceRepository(db)
	serviceService := services.NewServiceService(serviceRepo)

	verifyRepo := repository.NewVerifyRepository(rdb)
	verifyService := services.NewVerifyService(verifyRepo)

	subscriptionRepo := repository.NewSubscriptionRepository(db)
	subscriptionService := services.NewSubscriptionService(subscriptionRepo)

	transactionRepo := repository.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepo)

	h := handler.NewIncomingHandler(logger, rmpq, serviceService, verifyService, subscriptionService, transactionService)

	/**
	 * Routes Landing Page SUB & UNSUB
	 */
	r.Get("camp/:service", h.CampaignDirect)
	r.Get("camptool", h.CampaignTool)
	r.Get("p/:service", h.SubPage)
	r.Get("p/:service/faq", h.FaqPage)

	r.Post("cloudplay", h.OptIn)
	r.Post("galays", h.OptIn)

	r.Get("cloudplay/term", h.CloudPlayTermPage)

	r.Get("cloudplay", h.CloudPlaySubPage)
	r.Get("cloudplay/camp", h.CloudPlayCampaign)
	r.Get("cloudplay/campbill", h.CloudPlayCampaignBillable)

	r.Get("cloudplay1", h.CloudPlaySub1Page)
	r.Get("cloudplay1/camp", h.CloudPlaySub1CampaignPage)

	r.Get("cloudplay2", h.CloudPlaySub2Page)
	r.Get("cloudplay2/camp", h.CloudPlaySub2CampaignPage)

	r.Get("cloudplay3", h.CloudPlaySub3Page)
	r.Get("cloudplay3/camp", h.CloudPlaySub3CampaignPage)

	r.Get("cloudplay4", h.CloudPlaySub4Page)
	r.Get("cloudplay4/camp", h.CloudPlaySub4CampaignPage)

	r.Get("cloudplay/unsub", h.CloudPlayUnsubPage)

	r.Get("cbtsel", h.CallbackUrl)

	r.Get("cloudplay/camptool", h.CampaignTool)
	r.Get("cloudplay/camptooldynamic", h.CampaignToolDynamic)

	r.Get("galays", h.GalaysSubPage)
	r.Get("galays/camp", h.GalaysCampaign)
	r.Get("galays/campbill", h.GalaysCampaignBillable)

	r.Get("galays1", h.GalaysSub1Page)
	r.Get("galays1/camp", h.GalaysSub1CampaignPage)

	r.Get("galays/camptool", h.CampaignTool)
	r.Get("galays/camptooldynamic", h.CampaignToolDynamic)

	r.Get("success", h.CallbackUrl)
	r.Get("cancel", h.CallbackUrl)

	r.Post("auth/:category", h.Auth)

	/**``
	 * Routes Another CP
	 */
	r.Get("cpsam", h.CloudPlaySub1Page)
	r.Get("cpsam/camp", h.CloudPlaySub1CampaignPage)

	/**
	 * Routes MO & DR
	 */
	r.Get("notif/mo", h.MessageOriginated)
	r.Get("mo", h.MessageOriginated)

	/**
	 * Routes Report
	 */
	report := r.Group("report")
	report.Get("status", h.SelectStatus)
	report.Get("statusdetail", h.SelectStatusDetail)
	report.Get("adnet", h.SelectAdnet)
	report.Get("daily", h.ReportDaily)
	report.Post("arpu", h.AveragePerUser)

	return r
}
