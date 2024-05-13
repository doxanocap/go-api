package processor

import (
	"auth-api/internal/interfaces"
	"auth-api/internal/models"
	"auth-api/internal/processor/apis"
	"auth-api/internal/processor/cache"
	"auth-api/internal/processor/queue"
	"go.uber.org/zap"
)

type Processor struct {
	log                   *zap.Logger
	config                *models.Config
	cacheProvider         interfaces.ICacheProvider
	queueConsumerProvider interfaces.IQueueProducerProvider

	apisProcessor  interfaces.IAPIsProcessor
	cacheProcessor interfaces.ICacheProcessor
	queueProcessor interfaces.IQueueProcessor
}

func InitProcessor(
	queueConsumerProvider interfaces.IQueueProducerProvider, cacheProvider interfaces.ICacheProvider,
	config *models.Config, log *zap.Logger) *Processor {
	return &Processor{
		log:    log,
		config: config,

		cacheProvider:         cacheProvider,
		queueConsumerProvider: queueConsumerProvider,

		apisProcessor:  apis.NewAPIsProcessor(config, log.Named("[APIs]")),
		cacheProcessor: cache.NewCacheProcessor(cacheProvider, log.Named("[CACHE]")),
		queueProcessor: queue.NewQueueProcessor(config, queueConsumerProvider, log.Named("[QUEUE]")),
	}
}

func (p *Processor) Cache() interfaces.ICacheProcessor {
	return p.cacheProcessor
}

func (p *Processor) Queue() interfaces.IQueueProcessor {
	return p.queueProcessor
}

func (p *Processor) APIs() interfaces.IAPIsProcessor {
	return p.apisProcessor
}
