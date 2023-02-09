package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/hayesmp/jurassic-cage-service/internal"
	"github.com/hayesmp/jurassic-cage-service/internal/models"
	"github.com/hayesmp/jurassic-cage-service/postgres"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"math/rand"
	"path/filepath"
)

/* Some test helpers */

func loadEnv(logger zerolog.Logger) *internal.Config {
	envFile, err := filepath.Abs("../.env")
	if err != nil {
		logger.Fatal().Err(err)
	}

	err = godotenv.Load(envFile)
	if err != nil {
		logger.Fatal().Err(err)
	}
	config, err := internal.ParseConfigFromEnvironment("")
	if err != nil {
		logger.Fatal().Err(err)
	}

	return config
}

// test db interface
func initTestDb(log zerolog.Logger, config internal.Config) *postgres.Queries {
	connectionString := fmt.Sprintf("host=%s port=%s dbname=%s sslmode=disable",
		config.PostgresHost,
		config.PostgresPort,
		config.PostgresDb)
	if len(config.PostgresUser) > 0 {
		connectionString = fmt.Sprintf("%s user=%s", connectionString, config.PostgresUser)
	}
	if len(config.PostgresPw) > 0 {
		connectionString = fmt.Sprintf("%s password=%s", connectionString, config.PostgresPw)
	}

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Err(err).Msg("could not initialize db. trying again...")
		return nil
	}
	sqlcDb := postgres.New(db)

	return sqlcDb

}

// set up a test service
func integrationTestSetup(config *internal.Config) *JurassicCageService {
	logger := zerolog.Logger{}

	service := &JurassicCageService{
		config: config,
		env:    config.Env,
	}
	service.logger = logger
	service.db = initTestDb(logger, *config)

	return service
}

func createTestCage(ctx context.Context, service *JurassicCageService) models.CageResponse {
	name := fmt.Sprintf("Cage %v", uuid.New())
	cage, err := service.DbCreateCage(ctx, models.Cage{
		Name: name,
	})
	if err != nil {
		service.logger.Error().Err(err).Msg("Failed to save test cage")
	}

	return models.CageResponse{
		ID:                     cage.ID,
		Name:                   cage.Name,
		PredominateEatingHabit: cage.PredominateEatingHabit.String(),
	}
}

func deleteTestCage(ctx context.Context, service *JurassicCageService, cageId uuid.UUID) {
	err := service.DbDeleteCage(ctx, cageId)
	if err != nil {
		service.logger.Error().Err(err).Msg("failed to delete cage")
	}
}

func createTestHerbivore(ctx context.Context, service *JurassicCageService) models.Dinosaur {
	name := fmt.Sprintf("Dino %v", uuid.New())
	species := randomHerbivoreSpecies()
	dino, err := service.DbCreateDinosaur(ctx, models.DinosaurRequest{
		Name:        name,
		Species:     species,
		EatingHabit: models.ParseEatingHabit(species),
	})
	if err != nil {
		service.logger.Error().Err(err).Msg("Failed to save test cage")
	}
	return dino
}
func createTestCarnivore(ctx context.Context, service *JurassicCageService) models.Dinosaur {
	name := fmt.Sprintf("Dino %v", uuid.New())
	species := randomCarnivoreSpecies()
	dino, err := service.DbCreateDinosaur(ctx, models.DinosaurRequest{
		Name:        name,
		Species:     species,
		EatingHabit: models.ParseEatingHabit(species),
	})
	if err != nil {
		service.logger.Error().Err(err).Msg("Failed to save test cage")
	}
	return dino
}
func createTestDinosaur(ctx context.Context, service *JurassicCageService) models.Dinosaur {
	name := fmt.Sprintf("Dino %v", uuid.New())
	species := randomSpecies()
	dino, err := service.DbCreateDinosaur(ctx, models.DinosaurRequest{
		Name:        name,
		Species:     species,
		EatingHabit: models.ParseEatingHabit(species),
	})
	if err != nil {
		service.logger.Error().Err(err).Msg("Failed to save test cage")
	}
	return dino
}

func randomSpecies() string {
	species := []string{"brachiosaurus", "stegosaurus", "ankylosaurus", "triceratops", "tyrannosaurus", "velociraptor", "spinosaurus", "megalosaurus"}
	min := 0
	max := 7
	index := rand.Intn(max-min+1) + min
	return species[index]
}

func randomCarnivoreSpecies() string {
	species := []string{"tyrannosaurus", "velociraptor", "spinosaurus", "megalosaurus"}
	min := 0
	max := 3
	index := rand.Intn(max-min+1) + min
	return species[index]
}

func randomHerbivoreSpecies() string {
	species := []string{"brachiosaurus", "stegosaurus", "ankylosaurus", "triceratops"}
	min := 0
	max := 3
	index := rand.Intn(max-min+1) + min
	return species[index]
}

func deleteTestDinosaur(ctx context.Context, service *JurassicCageService, dinoId uuid.UUID) {
	err := service.DbDeleteDinosaur(ctx, dinoId)
	if err != nil {
		service.logger.Error().Err(err).Msg("failed to delete test dinosaur")
	}
}
