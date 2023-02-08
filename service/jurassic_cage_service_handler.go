package service

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hayesmp/jurassic-cage-service/internal"
	"github.com/hayesmp/jurassic-cage-service/internal/models"
	"github.com/hayesmp/jurassic-cage-service/postgres"
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

	cage, err := s.db.GetCageAndDinosaurs(c, uuid)
	if err != nil {
		msg := "failed to retrieve cage from local db"
		s.logger.Error().Err(err).Msg(msg)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	c.IndentedJSON(http.StatusOK, models.Cage{
		ID:                     cage.ID,
		Name:                   cage.Name.String,
		Status:                 models.Status(cage.Status.Int32),
		PredominateEatingHabit: models.EatingHabit(cage.PredominateEatingHabit.Int32),
		Dinosaurs:              []models.Dinosaur{},
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
	cage, err := s.db.UpsertCage(c, postgres.UpsertCageParams{
		Name:                   sql.NullString{newCage.Name, internal.ValidateString(newCage.Name)},
		Status:                 sql.NullInt32{models.ACTIVE.Int32(), internal.ValidateInt32(models.ACTIVE.Int32())},
		PredominateEatingHabit: sql.NullInt32{models.UNKNOWN.Int32(), internal.ValidateInt32(models.UNKNOWN.Int32())},
	})
	if err != nil {
		msg := "failed to save cage to local db"
		s.logger.Error().Err(err).Msg(msg)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	c.IndentedJSON(http.StatusOK, cage)
}

func (s *JurassicCageService) GetAllCages(c *gin.Context) {
	cages, err := s.DbGetAllCages(c)
	if err != nil {
		msg := "failed to get cages"
		s.logger.Error().Err(err).Msg(msg)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": msg})
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

func (s *JurassicCageService) GetDinosaur(c *gin.Context) {
	id := c.Param("id")

	uuid, err := uuid.Parse(id)
	if err != nil {
		msg := "failed to parsed uuid"
		s.logger.Error().Err(err).Msg(msg)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	dinosaur, err := s.db.GetDinosaurAndCage(c, uuid)
	if err != nil && err == sql.ErrNoRows {
		msg := "dinosaur not found"
		s.logger.Error().Err(err).Msg(msg)
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": msg})
		return
	}
	if err != nil {
		msg := "error retrieving dinosaur from local db"
		s.logger.Error().Err(err).Msg(msg)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	c.IndentedJSON(http.StatusOK, dinosaur)
}

func (s *JurassicCageService) CreateDinosaur(c *gin.Context) {
	var newDinosaurRequest models.DinosaurRequest
	var newDinosaur models.Dinosaur

	err := c.BindJSON(&newDinosaurRequest)
	if err != nil {
		msg := "failed to bind json newDinosaur"
		s.logger.Error().Err(err).Msg(msg)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	newDinosaur.Name = newDinosaurRequest.Name
	newDinosaur.Species = models.ParseSpecies(newDinosaurRequest.Species)
	newDinosaur.EatingHabit = newDinosaur.Species.EatingHabit()

	upsertDinoParams := postgres.UpsertDinosaurParams{
		Name:        sql.NullString{newDinosaur.Name, internal.ValidateString(newDinosaur.Name)},
		Species:     sql.NullInt32{newDinosaur.Species.Int32(), internal.ValidateInt32(newDinosaur.Species.Int32())},
		EatingHabit: sql.NullInt32{newDinosaur.Species.EatingHabit().Int32(), internal.ValidateInt32(newDinosaur.Species.EatingHabit().Int32())},
	}

	/*
		if len(newDinosaurRequest.CageId) > 0 {
			cageUuid, err := uuid.Parse(newDinosaurRequest.CageId)
			if err != nil {
				msg := "failed to parsed uuid"
				s.logger.Error().Err(err).Msg(msg)
				c.IndentedJSON(http.StatusBadRequest, gin.H{"error": msg})
				return
			}
			upsertDinoParams.CageID = uuid.NullUUID{cageUuid, internal.ValidateUuid(cageUuid)}
		} else {
			upsertDinoParams.CageID = uuid.NullUUID{uuid.Nil, true}
		}
	*/
	dinosaur, err := s.db.UpsertDinosaur(c, upsertDinoParams)
	if err != nil {
		msg := "failed to save dinosaur to local db"
		s.logger.Error().Err(err).Msg(msg)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	dinoResponse := &models.Dinosaur{
		ID:          dinosaur.ID,
		Name:        dinosaur.Name.String,
		Species:     models.Species(dinosaur.Species.Int32),
		EatingHabit: models.EatingHabit(dinosaur.EatingHabit.Int32),
	}

	if dinosaur.CageID.UUID != uuid.Nil {
		cage, err := s.db.GetCage(c, dinosaur.CageID.UUID)
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
		dinoResponse.Cage = models.Cage{
			Name:   cage.Name.String,
			Status: models.Status(cage.Status.Int32),
			ID:     cage.ID,
		}
	}

	c.IndentedJSON(http.StatusOK, dinoResponse)
}

func (s *JurassicCageService) GetAllDinosaurs(c *gin.Context) {
	dinosaurs, err := s.DbGetAllDinosaurs(c)
	if err != nil {
		msg := "failed to get dinosaurs"
		s.logger.Error().Err(err).Msg(msg)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": msg})
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

//func (s *JurassicCageService) RemoveDinosaurFromCage(c *gin.Context) {
//	id := c.Param("id")
//
//	uuId, err := uuid.Parse(id)
//	if err != nil {
//		msg := "failed to parsed uuid"
//		s.logger.Error().Err(err).Msg(msg)
//		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": msg})
//		return
//	}
//	// set cage id to nil
//	dinosaur, err := s.db.UpdateDinosaurCage(c, postgres.UpdateDinosaurCageParams{
//		CageID: uuid.NullUUID{uuid.Nil, true},
//		ID:     uuId,
//	})
//	if err != nil && err == sql.ErrNoRows {
//		msg := "dinosaur not found"
//		s.logger.Error().Err(err).Msg(msg)
//		c.IndentedJSON(http.StatusNotFound, gin.H{"error": msg})
//		return
//	}
//	if err != nil {
//		msg := "error updating dinosaur on local db"
//		s.logger.Error().Err(err).Msg(msg)
//		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": msg})
//		return
//	}
//
//	c.IndentedJSON(http.StatusOK, dinosaur)
//}
