package service

import (
	"database/sql"
	"fmt"
	"github.com/hayesmp/jurassic-cage-service/internal"
	"github.com/hayesmp/jurassic-cage-service/postgres"
	"github.com/rs/zerolog"
	"time"
)

type JurassicCageService struct {
	env    string
	logger zerolog.Logger
	config *internal.Config
	db     *postgres.Queries
}

func (s *JurassicCageService) Init(config *internal.Config, logger zerolog.Logger, db *postgres.Queries) *JurassicCageService {
	service := &JurassicCageService{
		config: config,
		env:    config.Env,
	}
	service.logger = logger
	s.initDb()

	return service
}

func (s *JurassicCageService) initDb() {
	s.logger.Info().Msg("initializing db")

	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		s.config.PostgresHost,
		s.config.PostgresPort,
		s.config.PostgresUser,
		s.config.PostgresPw,
		s.config.PostgresDb)

	for s.db == nil {
		db, err := sql.Open(
			"postgres", connectionString,
		)
		if err != nil {
			s.logger.Err(err).Msg("could not initialize db. trying again...")
			time.Sleep(2 * time.Second)
			continue
		}
		sqlcDb := postgres.New(db)

		s.db = sqlcDb
	}

	s.logger.Info().Msg("db initialized")
}
