package queue

import (
	"auth-api/internal/interfaces"
	"auth-api/internal/models"
	"auth-api/internal/processor/queue/producers"
	"go.uber.org/zap"
)

type Queue struct {
	log              *zap.Logger
	config           *models.Config
	consumerProvider interfaces.IQueueProducerProvider

	producersProcessor interfaces.IQueueProducersProcessor
}

func NewQueueProcessor(config *models.Config,
	consumersProvider interfaces.IQueueProducerProvider, log *zap.Logger) *Queue {
	return &Queue{
		log:              log,
		config:           config,
		consumerProvider: consumersProvider,

		producersProcessor: producers.NewProducersProcessor(config, consumersProvider, log.Named("[PRODUCERS]")),
	}
}

func (q *Queue) Producers() interfaces.IQueueProducersProcessor {
	return q.producersProcessor
}
