package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-playground/assert/v2"
	"github.com/google/uuid"
	"github.com/hayesmp/jurassic-cage-service/internal/models"
	"github.com/rs/zerolog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var logger = zerolog.New(os.Stderr)

// GET "/cage/:id"
func TestJurassicCageService_GetCage(t *testing.T) {
	config := loadEnv(logger)
	service := integrationTestSetup(config)
	router := service.SetupRouter()
	ctx := context.Background()
	cage := createTestCage(ctx, service)
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", fmt.Sprintf("/cage/%s", cage.ID.String()), nil)
	router.ServeHTTP(w, req)

	logger.Info().Msgf("%+v", w.Body)

	var newCage models.CageResponse
	err := json.Unmarshal(w.Body.Bytes(), &newCage)
	if err != nil {
		logger.Error().Err(err).Msg("failed to unmarshal resp body")
	}

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "ACTIVE", newCage.Status)
	assert.Equal(t, "Unknown", newCage.PredominateEatingHabit)

	deleteTestCage(ctx, service, cage.ID)
}

func TestJurassicCageService_GetCage_Notfound(t *testing.T) {
	config := loadEnv(logger)
	service := integrationTestSetup(config)
	router := service.SetupRouter()
	w := httptest.NewRecorder()

	cage := models.Cage{ID: uuid.MustParse("c1216ec0-314c-480a-9770-a2e68e1a531c")}
	req, _ := http.NewRequest("GET", fmt.Sprintf("/cage/%s", cage.ID.String()), nil)
	router.ServeHTTP(w, req)

	var errorResp models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResp)
	if err != nil {
		logger.Error().Err(err).Msg("failed to unmarshal resp body")
	}
	assert.Equal(t, 404, w.Code)
	assert.Equal(t, "failed to retrieve cage from local db", errorResp.Error)
}
func TestJurassicCageService_CreateCage(t *testing.T) {
	config := loadEnv(logger)
	service := integrationTestSetup(config)
	ctx := context.Background()
	router := service.SetupRouter()
	w := httptest.NewRecorder()
	var jsonData = []byte(`{"name": "morpheus"}`)

	//cage := models.Cage{ID: uuid.MustParse("c1216ec0-314c-480a-9770-a2e68e1a531c")}
	req, _ := http.NewRequest("POST", "/cage", bytes.NewBuffer(jsonData))
	router.ServeHTTP(w, req)

	logger.Info().Msgf("%+v", w)

	var newCage models.CageResponse
	err := json.Unmarshal(w.Body.Bytes(), &newCage)
	if err != nil {
		logger.Error().Err(err).Msg("failed to unmarshal resp body")
	}

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "ACTIVE", newCage.Status)
	assert.Equal(t, "Unknown", newCage.PredominateEatingHabit)

	deleteTestCage(ctx, service, newCage.ID)
}

func TestJurassicCageService_GetAllCages(t *testing.T) {
	config := loadEnv(logger)
	service := integrationTestSetup(config)
	router := service.SetupRouter()
	ctx := context.Background()
	cage1 := createTestCage(ctx, service)
	cage2 := createTestCage(ctx, service)
	cage3 := createTestCage(ctx, service)
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/cage", nil)
	router.ServeHTTP(w, req)

	logger.Info().Msgf("%+v", w.Body)

	var newCages []models.CageResponse
	err := json.Unmarshal(w.Body.Bytes(), &newCages)
	if err != nil {
		logger.Error().Err(err).Msg("failed to unmarshal resp body")
	}

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, true, len(newCages) >= 3)

	// cleanup test data
	deleteTestCage(ctx, service, cage1.ID)
	deleteTestCage(ctx, service, cage2.ID)
	deleteTestCage(ctx, service, cage3.ID)
}

func TestJurassicCageService_GetAllCages_FilterByDown(t *testing.T) {
	config := loadEnv(logger)
	service := integrationTestSetup(config)
	router := service.SetupRouter()
	ctx := context.Background()
	cage1 := createTestCage(ctx, service)
	cage2, _ := service.DbCreateCage(ctx, models.Cage{
		Name: "Up Cage",
	})
	cage3, _ := service.DbCreateCage(ctx, models.Cage{
		Name: "Down Cage",
	})
	cage3R, _ := service.DbUpdateCage(ctx, models.Cage{
		Name:   cage3.Name,
		Status: models.DOWN,
	})

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/cage?status=down", nil)
	router.ServeHTTP(w, req)

	logger.Info().Msgf("%+v", w.Body)

	var newCages []models.CageResponse
	err := json.Unmarshal(w.Body.Bytes(), &newCages)
	if err != nil {
		logger.Error().Err(err).Msg("failed to unmarshal resp body")
	}

	var downCage bool
	var upCage bool
	for _, newCage := range newCages {
		if newCage.Name == cage3R.Name {
			downCage = true
		}
		if newCage.Name == cage2.Name {
			upCage = true
		}
	}

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, true, downCage)
	assert.Equal(t, false, upCage)

	// cleanup test data
	deleteTestCage(ctx, service, cage1.ID)
	deleteTestCage(ctx, service, cage2.ID)
	deleteTestCage(ctx, service, cage3.ID)
}

func TestJurassicCageService_UpdateCage(t *testing.T) {
	config := loadEnv(logger)
	service := integrationTestSetup(config)
	router := service.SetupRouter()
	ctx := context.Background()
	cage := createTestCage(ctx, service)
	w := httptest.NewRecorder()

	var jsonData = []byte(`{"status": "down"}`)

	req, _ := http.NewRequest("PUT", fmt.Sprintf("/cage/%s", cage.ID.String()), bytes.NewBuffer(jsonData))
	router.ServeHTTP(w, req)

	//logger.Info().Msgf("%+v", w.Body)

	var newCage models.CageResponse
	err := json.Unmarshal(w.Body.Bytes(), &newCage)
	if err != nil {
		logger.Error().Err(err).Msg("failed to unmarshal resp body")
	}

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, newCage.Status, "DOWN")

	// cleanup test data
	deleteTestCage(ctx, service, cage.ID)
}

func TestJurassicCageService_UpdateCage_Failure(t *testing.T) {
	config := loadEnv(logger)
	service := integrationTestSetup(config)
	router := service.SetupRouter()
	ctx := context.Background()
	cage := createTestCage(ctx, service)
	dino := createTestDinosaur(ctx, service)

	err := service.DbAddDinosaurToCage(ctx, cage.ID, dino)
	if err != nil {
		logger.Error().Err(err).Msg("failed to add dino to cage")
	}

	w := httptest.NewRecorder()

	var jsonData = []byte(`{"status": "down"}`)

	req, _ := http.NewRequest("PUT", fmt.Sprintf("/cage/%s", cage.ID.String()), bytes.NewBuffer(jsonData))
	router.ServeHTTP(w, req)

	var errorResp models.ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &errorResp)
	if err != nil {
		logger.Error().Err(err).Msg("failed to unmarshal resp body")
	}
	assert.Equal(t, 400, w.Code)
	assert.Equal(t, "cannot set cage to DOWN with 1 dinosaurs in cage", errorResp.Error)

	// cleanup test data
	deleteTestCage(ctx, service, cage.ID)
	deleteTestDinosaur(ctx, service, dino.ID)
}

func TestJurassicCageService_GetDinosaur(t *testing.T) {
	config := loadEnv(logger)
	service := integrationTestSetup(config)
	router := service.SetupRouter()
	ctx := context.Background()
	dino := createTestDinosaur(ctx, service)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", fmt.Sprintf("/dinosaur/%s", dino.ID.String()), nil)
	router.ServeHTTP(w, req)

	logger.Info().Msgf("%+v", w.Body)

	var newDino models.DinosaurResponse
	err := json.Unmarshal(w.Body.Bytes(), &newDino)
	if err != nil {
		logger.Error().Err(err).Msg("failed to unmarshal resp body")
	}
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, dino.EatingHabit.String(), newDino.EatingHabit)
	assert.Equal(t, dino.Species.String(), newDino.Species)

	// cleanup test data
	deleteTestDinosaur(ctx, service, dino.ID)
}
func TestJurassicCageService_GetDinosaur_NotFound(t *testing.T) {
	config := loadEnv(logger)
	service := integrationTestSetup(config)
	router := service.SetupRouter()

	w := httptest.NewRecorder()

	dinoId := uuid.New()

	req, _ := http.NewRequest("GET", fmt.Sprintf("/dinosaur/%s", dinoId.String()), nil)
	router.ServeHTTP(w, req)

	var errorResp models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResp)
	if err != nil {
		logger.Error().Err(err).Msg("failed to unmarshal resp body")
	}
	assert.Equal(t, 404, w.Code)
	assert.Equal(t, "error retrieving dinosaur from local db", errorResp.Error)
}

func TestJurassicCageService_CreateDinosaur(t *testing.T) {
	config := loadEnv(logger)
	service := integrationTestSetup(config)
	ctx := context.Background()
	router := service.SetupRouter()
	w := httptest.NewRecorder()
	species := randomSpecies()
	eatingHabit := models.ParseEatingHabit(species)
	var jsonData = []byte(fmt.Sprintf(`{"name": "Dino %s", "species":"%s", "eating_habit":"%s"}`,
		uuid.New().String(), species, eatingHabit))

	//cage := models.Cage{ID: uuid.MustParse("c1216ec0-314c-480a-9770-a2e68e1a531c")}
	req, _ := http.NewRequest("POST", "/dinosaur", bytes.NewBuffer(jsonData))
	router.ServeHTTP(w, req)

	logger.Info().Msgf("%+v", w)

	var newDino models.DinosaurResponse
	err := json.Unmarshal(w.Body.Bytes(), &newDino)
	if err != nil {
		logger.Error().Err(err).Msg("failed to unmarshal resp body")
	}

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, species, strings.ToLower(newDino.Species))
	assert.Equal(t, eatingHabit, strings.ToLower(newDino.EatingHabit))

	// cleanup data
	deleteTestCage(ctx, service, newDino.ID)
}

func TestJurassicCageService_GetAllDinosaurs(t *testing.T) {
	config := loadEnv(logger)
	service := integrationTestSetup(config)
	router := service.SetupRouter()
	ctx := context.Background()
	dino1 := createTestDinosaur(ctx, service)
	dino2 := createTestDinosaur(ctx, service)
	dino3 := createTestDinosaur(ctx, service)
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/dinosaur", nil)
	router.ServeHTTP(w, req)

	logger.Info().Msgf("%+v", w.Body)

	var newDinos []models.DinosaurResponse
	err := json.Unmarshal(w.Body.Bytes(), &newDinos)
	if err != nil {
		logger.Error().Err(err).Msg("failed to unmarshal resp body")
	}

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, true, len(newDinos) >= 3)

	// cleanup test data
	deleteTestDinosaur(ctx, service, dino1.ID)
	deleteTestDinosaur(ctx, service, dino2.ID)
	deleteTestDinosaur(ctx, service, dino3.ID)
}

func TestJurassicCageService_GetAllDinosaurs_FilterBySpecies(t *testing.T) {
	config := loadEnv(logger)
	service := integrationTestSetup(config)
	router := service.SetupRouter()
	ctx := context.Background()
	dino1 := createTestDinosaur(ctx, service)
	dino2 := createTestDinosaur(ctx, service)
	dino3 := createTestDinosaur(ctx, service)
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", fmt.Sprintf("/dinosaur?species=%s", dino1.Species), nil)
	router.ServeHTTP(w, req)

	logger.Info().Msgf("%+v", w.Body)

	var newDinos []models.DinosaurResponse
	err := json.Unmarshal(w.Body.Bytes(), &newDinos)
	if err != nil {
		logger.Error().Err(err).Msg("failed to unmarshal resp body")
	}

	var species bool
	for _, newdino := range newDinos {
		if newdino.Name == dino1.Name {
			species = true
		}
	}

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, true, len(newDinos) >= 3)
	assert.Equal(t, true, species)

	// cleanup test data
	deleteTestDinosaur(ctx, service, dino1.ID)
	deleteTestDinosaur(ctx, service, dino2.ID)
	deleteTestDinosaur(ctx, service, dino3.ID)
}

func TestJurassicCageService_AddDinosaurToCage(t *testing.T) {
	config := loadEnv(logger)
	service := integrationTestSetup(config)
	router := service.SetupRouter()
	ctx := context.Background()
	dino := createTestDinosaur(ctx, service)
	cage := createTestCage(ctx, service)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("PUT", fmt.Sprintf("/dinosaur/%s/%s", dino.ID.String(), cage.ID.String()), nil)
	router.ServeHTTP(w, req)

	logger.Info().Msgf("%+v", w.Body)

	var newDino models.DinosaurResponse
	err := json.Unmarshal(w.Body.Bytes(), &newDino)
	if err != nil {
		logger.Error().Err(err).Msg("failed to unmarshal resp body")
	}
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, cage.ID.String(), newDino.CageId.String())
	assert.Equal(t, cage.Name, newDino.CageName)

	// cleanup test data
	deleteTestDinosaur(ctx, service, dino.ID)
	deleteTestCage(ctx, service, cage.ID)
}

func TestJurassicCageService_AddDinosaurToCage_EatingHabitError(t *testing.T) {
	config := loadEnv(logger)
	service := integrationTestSetup(config)
	router := service.SetupRouter()
	ctx := context.Background()
	herb := createTestHerbivore(ctx, service)
	carn := createTestCarnivore(ctx, service)
	cage := createTestCage(ctx, service)
	err := service.DbAddDinosaurToCage(ctx, cage.ID, herb)
	if err != nil {
		logger.Error().Err(err).Msg("failed to add herb dino to cage")
	}

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("PUT", fmt.Sprintf("/dinosaur/%s/%s", carn.ID.String(), cage.ID.String()), nil)
	router.ServeHTTP(w, req)

	var errorResp models.ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &errorResp)
	if err != nil {
		logger.Error().Err(err).Msg("failed to unmarshal resp body")
	}

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, fmt.Sprintf("%s cannot be added to cage %s due to conflict eating habit", carn.Name, cage.ID.String()), errorResp.Error)

	// cleanup test data
	deleteTestCage(ctx, service, cage.ID)
	deleteTestDinosaur(ctx, service, herb.ID)
	deleteTestDinosaur(ctx, service, carn.ID)
}
func TestJurassicCageService_AddDinosaurToCage_FullCageError(t *testing.T) {
	config := loadEnv(logger)
	service := integrationTestSetup(config)
	router := service.SetupRouter()
	ctx := context.Background()
	cage := createTestCage(ctx, service)
	var herbIds []uuid.UUID
	for i := 0; i < 4; i++ {
		herb := createTestHerbivore(ctx, service)
		err := service.DbAddDinosaurToCage(ctx, cage.ID, herb)
		if err != nil {
			logger.Error().Err(err).Msg("failed to add herb dino to cage")
		}
		herbIds = append(herbIds, herb.ID)
	}

	herb2 := createTestHerbivore(ctx, service)
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("PUT", fmt.Sprintf("/dinosaur/%s/%s", herb2.ID.String(), cage.ID.String()), nil)
	router.ServeHTTP(w, req)

	var errorResp models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResp)
	if err != nil {
		logger.Error().Err(err).Msg("failed to unmarshal resp body")
	}

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, "cage dinosaur count is already at max capacity", errorResp.Error)

	// cleanup test data
	deleteTestCage(ctx, service, cage.ID)
	deleteTestDinosaur(ctx, service, herb2.ID)
	for _, herbId := range herbIds {
		deleteTestDinosaur(ctx, service, herbId)
	}
}
func TestJurassicCageService_AddDinosaurToCage_CageStatusError(t *testing.T) {
	config := loadEnv(logger)
	service := integrationTestSetup(config)
	router := service.SetupRouter()
	ctx := context.Background()
	cage := createTestCage(ctx, service)
	herb := createTestHerbivore(ctx, service)

	cageDown, err := service.DbUpdateCage(ctx, models.Cage{
		Name:   cage.Name,
		Status: models.DOWN,
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to update cage status")
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/dinosaur/%s/%s", herb.ID.String(), cageDown.ID.String()), nil)
	router.ServeHTTP(w, req)

	var errorResp models.ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &errorResp)
	if err != nil {
		logger.Error().Err(err).Msg("failed to unmarshal resp body")
	}

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, "cage Status is not in an ACTIVE state", errorResp.Error)

	// cleanup test data
	deleteTestCage(ctx, service, cage.ID)
	deleteTestDinosaur(ctx, service, herb.ID)
}
func TestJurassicCageService_MoveDinosaurFromCageToCage(t *testing.T) {
	config := loadEnv(logger)
	service := integrationTestSetup(config)
	router := service.SetupRouter()
	ctx := context.Background()
	cage := createTestCage(ctx, service)
	herb := createTestHerbivore(ctx, service)
	err := service.DbAddDinosaurToCage(ctx, cage.ID, herb)
	if err != nil {
		logger.Error().Err(err).Msg("failed to add herb dino to cage")
	}
	cage2 := createTestCage(ctx, service)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("PUT", fmt.Sprintf("/dinosaur/%s/%s", herb.ID.String(), cage2.ID.String()), nil)
	router.ServeHTTP(w, req)

	logger.Info().Msgf("%+v", w.Body)

	var newDino models.DinosaurResponse
	err = json.Unmarshal(w.Body.Bytes(), &newDino)
	if err != nil {
		logger.Error().Err(err).Msg("failed to unmarshal resp body")
	}
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, cage2.ID.String(), newDino.CageId.String())
	assert.Equal(t, cage2.Name, newDino.CageName)

	// cleanup test data
	deleteTestDinosaur(ctx, service, herb.ID)
	deleteTestCage(ctx, service, cage.ID)
	deleteTestCage(ctx, service, cage2.ID)
}
