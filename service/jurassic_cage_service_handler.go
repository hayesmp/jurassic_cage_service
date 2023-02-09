package service

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hayesmp/jurassic-cage-service/internal/models"
	"net/http"
)

func (s *JurassicCageService) GetCage(c *gin.Context) {
	id := c.Param("id")

	uuid, err := uuid.Parse(id)
	if err != nil {
		msg := "failed to parsed uuid"
		s.logger.Error().Err(err).Msg(msg)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	cage, err := s.DbGetCage(c, uuid)
	if err != nil {
		msg := "failed to retrieve cage from local db"
		s.logger.Error().Err(err).Msg(msg)
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": msg})
		return
	}

	var retDinos []models.DinosaurResponse
	for _, dino := range cage.Dinosaurs {
		retDinos = append(retDinos, models.DinosaurResponse{
			ID:          dino.ID,
			Name:        dino.Name,
			EatingHabit: dino.EatingHabit.String(),
			Species:     dino.Species.String(),
			CageId:      dino.Cage.ID,
			CageName:    dino.Cage.Name,
		})
	}
	c.IndentedJSON(http.StatusOK, models.CageResponse{
		ID:                     cage.ID,
		Name:                   cage.Name,
		Status:                 cage.Status.String(),
		PredominateEatingHabit: cage.PredominateEatingHabit.String(),
		Dinosaurs:              retDinos,
	})
}

func (s *JurassicCageService) CreateCage(c *gin.Context) {
	var newCage models.Cage

	err := c.BindJSON(&newCage)
	if err != nil {
		msg := "failed to bind json newCage"
		s.logger.Error().Err(err).Msg(msg)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	cage, err := s.DbCreateCage(c, newCage)
	if err != nil {
		msg := "failed to save cage to local db"
		s.logger.Error().Err(err).Msg(msg)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	c.IndentedJSON(http.StatusOK, models.CageResponse{
		ID:                     cage.ID,
		Name:                   cage.Name,
		Status:                 cage.Status.String(),
		PredominateEatingHabit: cage.PredominateEatingHabit.String(),
	})
}

func (s *JurassicCageService) GetAllCages(c *gin.Context) {
	status := c.Query("status")

	cages, err := s.DbGetAllCages(c, status)
	if err != nil {
		msg := "failed to get cages"
		s.logger.Error().Err(err).Msg(msg)
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": msg})
		return
	}

	var retCages []models.CageResponse
	for _, cage := range cages {

		var retDinos []models.DinosaurResponse
		for _, dino := range cage.Dinosaurs {
			retDinos = append(retDinos, models.DinosaurResponse{
				ID:          dino.ID,
				Name:        dino.Name,
				EatingHabit: dino.EatingHabit.String(),
				Species:     dino.Species.String(),
				CageId:      dino.Cage.ID,
				CageName:    dino.Cage.Name,
			})
		}
		retCages = append(retCages, models.CageResponse{
			ID:                     cage.ID,
			Name:                   cage.Name,
			Status:                 cage.Status.String(),
			PredominateEatingHabit: cage.PredominateEatingHabit.String(),
			Capacity:               cage.Capacity,
			Dinosaurs:              retDinos,
		})
	}

	c.IndentedJSON(http.StatusOK, retCages)
}

func (s *JurassicCageService) UpdateCage(c *gin.Context) {
	var updateCage models.CageRequest
	id := c.Param("id")

	err := c.BindJSON(&updateCage)
	if err != nil {
		msg := "failed to bind json newCage"
		s.logger.Error().Err(err).Msg(msg)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	cageUuid, err := uuid.Parse(id)
	if err != nil {
		msg := "failed to parsed uuid"
		s.logger.Error().Err(err).Msg(msg)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	cage, err := s.DbGetCageAndDinosaurs(c, cageUuid)
	if err != nil {
		msg := "failed to retrieve cage from local db"
		s.logger.Error().Err(err).Msg(msg)
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": msg})
		return
	}

	requestedStatus := models.ParseStatus(updateCage.Status)
	if requestedStatus == models.DOWN && cage.Status == models.ACTIVE {
		if len(cage.Dinosaurs) > 0 {
			msg := fmt.Sprintf("cannot set cage to DOWN with %d dinosaurs in cage", len(cage.Dinosaurs))
			s.logger.Error().Err(err).Msg(msg)
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}
	}

	// set the cage with the requested status
	cage.Status = requestedStatus

	cage, err = s.DbUpdateCage(c, cage)
	if err != nil {
		msg := "failed to update cage on local db"
		s.logger.Error().Err(err).Msg(msg)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	var dinosResponse []models.DinosaurResponse
	for _, dino := range cage.Dinosaurs {
		dinosResponse = append(dinosResponse, models.DinosaurResponse{
			ID:          dino.ID,
			Name:        dino.Name,
			Species:     dino.Species.String(),
			EatingHabit: dino.EatingHabit.String(),
			CageId:      dino.CageId,
			CageName:    dino.Cage.Name,
		})
	}

	c.IndentedJSON(http.StatusOK, models.CageResponse{
		ID:                     cage.ID,
		Name:                   cage.Name,
		Status:                 cage.Status.String(),
		PredominateEatingHabit: cage.PredominateEatingHabit.String(),
		Dinosaurs:              dinosResponse,
	})
}

func (s *JurassicCageService) GetDinosaur(c *gin.Context) {
	id := c.Param("id")

	uuid, err := uuid.Parse(id)
	if err != nil {
		msg := "failed to parsed uuid"
		s.logger.Error().Err(err).Msg(msg)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	dinosaur, err := s.DbGetDinosaur(c, uuid)
	if err != nil {
		msg := "error retrieving dinosaur from local db"
		s.logger.Error().Err(err).Msg(msg)
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": msg})
		return
	}

	c.IndentedJSON(http.StatusOK, models.DinosaurResponse{
		ID:          dinosaur.ID,
		Name:        dinosaur.Name,
		Species:     dinosaur.Species.String(),
		EatingHabit: dinosaur.EatingHabit.String(),
		CageName:    dinosaur.Cage.Name,
		CageId:      dinosaur.CageId,
	})
}

func (s *JurassicCageService) CreateDinosaur(c *gin.Context) {
	var newDinosaurRequest models.DinosaurRequest

	err := c.BindJSON(&newDinosaurRequest)
	if err != nil {
		msg := "failed to bind json newDinosaur"
		s.logger.Error().Err(err).Msg(msg)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	dinosaur, err := s.DbCreateDinosaur(c, newDinosaurRequest)
	if err != nil {
		msg := "failed to save dinosaur to local db"
		s.logger.Error().Err(err).Msg(msg)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	dinoResponse := &models.DinosaurResponse{
		ID:          dinosaur.ID,
		Name:        dinosaur.Name,
		Species:     dinosaur.Species.String(),
		EatingHabit: dinosaur.EatingHabit.String(),
	}

	if dinosaur.CageId != uuid.Nil {
		cage, err := s.db.GetCage(c, dinosaur.CageId)
		if err != nil && err == sql.ErrNoRows {
			msg := "no cage found on local db"
			s.logger.Error().Err(err).Msg(msg)
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}
		if err != nil {
			msg := "failed to retrieve cage from local db"
			s.logger.Error().Err(err).Msg(msg)
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}
		dinoResponse.CageId = cage.ID
		dinoResponse.CageName = cage.Name.String
	}

	c.IndentedJSON(http.StatusOK, dinoResponse)
}

func (s *JurassicCageService) GetAllDinosaurs(c *gin.Context) {
	species := c.Query("species")

	dinosaurs, err := s.DbGetAllDinosaurs(c, species)
	if err != nil {
		msg := "failed to get dinosaurs"
		s.logger.Error().Err(err).Msg(msg)
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": msg})
		return
	}

	var retDinos []models.DinosaurResponse
	for _, dino := range dinosaurs {

		retDinos = append(retDinos, models.DinosaurResponse{
			ID:          dino.ID,
			Name:        dino.Name,
			EatingHabit: dino.EatingHabit.String(),
			Species:     dino.Species.String(),
			CageId:      dino.Cage.ID,
			CageName:    dino.Cage.Name,
		})
	}

	c.IndentedJSON(http.StatusOK, retDinos)
}

// MAX_CAPACITY for cage
const MAX_CAPACITY = 4

// Also moves a dinosaur between cages
func (s *JurassicCageService) AddDinosaurToCage(c *gin.Context) {
	id := c.Param("id")
	cage_id := c.Param("cage_id")

	dinoUuid, err := uuid.Parse(id)
	if err != nil {
		msg := "failed to parse dino uuid"
		s.logger.Error().Err(err).Msg(msg)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}
	cageUuid, err := uuid.Parse(cage_id)
	if err != nil {
		msg := "failed to parse cage uuid"
		s.logger.Error().Err(err).Msg(msg)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	dinosaur, err := s.DbGetDinosaur(c, dinoUuid)
	if err != nil {
		msg := "failed to get dinosaur from local db"
		s.logger.Error().Err(err).Msg(msg)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	cage, err := s.DbGetCage(c, cageUuid)
	if err != nil {
		msg := "failed to get cage from local db"
		s.logger.Error().Err(err).Msg(msg)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	// Check if cage predominant eating habit is compatible
	if cage.PredominateEatingHabit != models.UNKNOWN {
		if dinosaur.EatingHabit != cage.PredominateEatingHabit {
			msg := fmt.Sprintf("%s cannot be added to cage %s due to conflict eating habit", dinosaur.Name, cageUuid.String())
			s.logger.Error().Err(err).Msg(msg)
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}
	}

	// Check if cage is already at max capacity
	dinoCount, err := s.DbGetCageDinosaurCount(c, cageUuid)
	if err != nil {
		msg := fmt.Sprintf("error retrieving dino count for cage %s", cage.ID.String())
		s.logger.Error().Err(err).Msg(msg)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}
	if dinoCount+1 > MAX_CAPACITY {
		msg := "cage dinosaur count is already at max capacity"
		s.logger.Error().Err(err).Msg(msg)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	// Check if the cage is ACTIVE or DOWN
	if cage.Status != models.ACTIVE {
		msg := "cage Status is not in an ACTIVE state"
		s.logger.Error().Err(err).Msg(msg)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	// Remove dinosaur from old cage if in old cage
	if dinosaur.CageId != uuid.Nil {
		err = s.DbRemoveDinosaurFromCage(c, dinosaur.CageId, dinosaur.ID)
		if err != nil {
			msg := fmt.Sprintf("failed to remove dinosaur %s from old cage %s", dinosaur.Name, dinosaur.Cage.Name)
			s.logger.Error().Err(err).Msg(msg)
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}
	}

	// If all checks succeed, add dinosaur to cage
	err = s.DbAddDinosaurToCage(c, cageUuid, dinosaur)
	if err != nil {
		msg := fmt.Sprintf("failed to add dinosaur %s to cage %s", dinosaur.Name, cage.Name)
		s.logger.Error().Err(err).Msg(msg)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	// Get dino count again
	dinoCount, err = s.DbGetCageDinosaurCount(c, cageUuid)
	if err != nil {
		msg := fmt.Sprintf("error retrieving dino count for cage %s", cage.ID.String())
		s.logger.Error().Err(err).Msg(msg)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	c.IndentedJSON(http.StatusOK, models.DinosaurResponse{
		ID:          dinosaur.ID,
		Name:        dinosaur.Name,
		Species:     dinosaur.Species.String(),
		EatingHabit: dinosaur.EatingHabit.String(),
		CageId:      cage.ID,
		CageName:    cage.Name,
	})
}
