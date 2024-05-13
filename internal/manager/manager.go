package manager

import (
	"auth-api/internal/interfaces"
	"auth-api/internal/models"
	"auth-api/internal/pkg/rabbitmq"
	"auth-api/internal/pkg/redis"
	"auth-api/internal/processor"
	"auth-api/internal/repository"
	"auth-api/internal/service"
	"auth-api/server"
	_ "github.com/jackc/pgx/v4/pgxpool"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Manager struct {
	db                    *sqlx.DB
	log                   *zap.Logger
	config                *models.Config
	cacheProvider         interfaces.ICacheProvider
	queueConsumerProvider interfaces.IQueueProducerProvider

	processor  interfaces.IProcessor
	repository interfaces.IRepository
	service    interfaces.IService
	server     interfaces.IServer
}

func InitManager(db *sqlx.DB, log *zap.Logger, config *models.Config,
	queueConsumerProvider *rabbitmq.ProducerClient, cacheProvider *redis.Conn) *Manager {
	m := &Manager{
		db:                    db,
		log:                   log,
		config:                config,
		cacheProvider:         cacheProvider,
		queueConsumerProvider: queueConsumerProvider,
	}

	m.processor = processor.InitProcessor(m.queueConsumerProvider, m.cacheProvider, config, log)
	m.repository = repository.InitRepository(db, config, m.Processor().Cache())
	m.service = service.InitService(m, config, log)
	m.server = server.InitServer(m, config, log)
	return m
}

func (m *Manager) Repository() interfaces.IRepository {
	return m.repository
}

func (m *Manager) Service() interfaces.IService {
	return m.service
}

func (m *Manager) Processor() interfaces.IProcessor {
	return m.processor
}

func (m *Manager) Server() interfaces.IServer {
	return m.server
}
