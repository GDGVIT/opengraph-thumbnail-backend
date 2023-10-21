package usersvc

import (
	"context"

	"github.com/GDGVIT/opengraph-thumbnail-backend/pkg/logger"
	"gorm.io/gorm"
)

type UserSvcImpl struct {
	gormDB        *gorm.DB
	logger        logger.Logger
	messageBroker MessageBroker
}

type MessageBroker interface {
	Publish(ctx context.Context, exchange, routingKey string, body []byte) error
}

func Handler(gormDB *gorm.DB, logger logger.Logger, messageBroker MessageBroker) *UserSvcImpl {
	return &UserSvcImpl{
		gormDB:        gormDB,
		logger:        logger,
		messageBroker: messageBroker,
	}
}
