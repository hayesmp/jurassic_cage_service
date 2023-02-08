package service

import (
	"database/sql"
	"fmt"
	"github.com/hayesmp/jurassic-cage-service/internal"
	"github.com/hayesmp/jurassic-cage-service/postgres"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"time"
)

type JurassicCageService struct {
	env    string
	logger zerolog.Logger
	config *internal.Config
	db     *postgres.Queries
}

func Init(config *internal.Config, logger zerolog.Logger) *JurassicCageService {
	service := &JurassicCageService{
		config: config,
		env:    config.Env,
	}
	service.logger = logger
	service.initDb()

	return service
}

func (s *JurassicCageService) initDb() {
	s.logger.Info().Msg("initializing db")

	connectionString := fmt.Sprintf("host=%s port=%s dbname=%s sslmode=disable",
		s.config.PostgresHost,
		s.config.PostgresPort,
		s.config.PostgresDb)
	if len(s.config.PostgresUser) > 0 {
		connectionString = fmt.Sprintf("%s user=%s", connectionString, s.config.PostgresUser)
	}
	if len(s.config.PostgresPw) > 0 {
		connectionString = fmt.Sprintf("%s password=%s", connectionString, s.config.PostgresPw)
	}

	s.logger.Debug().Msgf("%+v", connectionString)

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
