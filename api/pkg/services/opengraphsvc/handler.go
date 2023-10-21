package opengraphsvc

import (
	"github.com/GDGVIT/opengraph-thumbnail-backend/pkg/logger"
)

type OpenGraphSvcImpl struct {
	logger logger.Logger
}

func Handler(logger logger.Logger) *OpenGraphSvcImpl {
	return &OpenGraphSvcImpl{
		logger: logger,
	}
}
