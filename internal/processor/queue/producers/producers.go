package producers

import (
	"auth-api/config"
	"auth-api/internal/interfaces"
	"auth-api/logger"
)

type Producers struct {
	log              *logger.Logger
	config           *config.Config
	consumerProvider interfaces.IQueueProducerProvider
}

func NewProducersProcessor(
	config *config.Config,
	consumersProvider interfaces.IQueueProducerProvider,
	log *logger.Logger) *Producers {
	return &Producers{
		log:              log,
		config:           config,
		consumerProvider: consumersProvider,
	}
}
