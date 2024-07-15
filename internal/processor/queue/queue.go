package queue

import (
	"auth-api/config"
	"auth-api/internal/interfaces"
	"auth-api/internal/processor/queue/producers"
	"auth-api/logger"
)

type Queue struct {
	log              *logger.Logger
	config           *config.Config
	consumerProvider interfaces.IQueueProducerProvider

	producersProcessor interfaces.IQueueProducersProcessor
}

func NewQueueProcessor(config *config.Config,
	consumersProvider interfaces.IQueueProducerProvider, log *logger.Logger) *Queue {
	return &Queue{
		log:              log,
		config:           config,
		consumerProvider: consumersProvider,

		producersProcessor: producers.NewProducersProcessor(config, consumersProvider, log.WithModule("PRODUCERS")),
	}
}

func (q *Queue) Producers() interfaces.IQueueProducersProcessor {
	return q.producersProcessor
}
