package processor

import (
	"auth-api/config"
	"auth-api/internal/interfaces"
	"auth-api/internal/processor/apis"
	"auth-api/internal/processor/cache"
	"auth-api/internal/processor/queue"
	"auth-api/logger"
)

type Processor struct {
	log                   *logger.Logger
	config                *config.Config
	cacheProvider         interfaces.ICacheProvider
	queueConsumerProvider interfaces.IQueueProducerProvider

	apisProcessor  interfaces.IAPIsProcessor
	cacheProcessor interfaces.ICacheProcessor
	queueProcessor interfaces.IQueueProcessor
}

func InitProcessor(
	queueConsumerProvider interfaces.IQueueProducerProvider, cacheProvider interfaces.ICacheProvider,
	config *config.Config, log *logger.Logger) *Processor {
	return &Processor{
		log:    log,
		config: config,

		cacheProvider:         cacheProvider,
		queueConsumerProvider: queueConsumerProvider,

		apisProcessor:  apis.NewAPIsProcessor(config, log.WithModule("APIs")),
		cacheProcessor: cache.NewCacheProcessor(cacheProvider, log.WithModule("CACHE")),
		queueProcessor: queue.NewQueueProcessor(config, queueConsumerProvider, log.WithModule("QUEUE")),
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
