package app

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/idprm/go-linkit-tsel/internal/logger"
	"github.com/redis/go-redis/v9"
	"github.com/wiliehidayat87/rmqp"
)

func StartApplication(db *sql.DB, rmq rmqp.AMQP, rds *redis.Client, logger *logger.Logger) *fiber.App {
	return mapUrls(db, rmq, rds, logger)
}
