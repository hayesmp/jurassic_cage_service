package service

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hayesmp/jurassic-cage-service/internal"
	"github.com/hayesmp/jurassic-cage-service/internal/models"
	"github.com/hayesmp/jurassic-cage-service/postgres"
)

/* Db Methods */

func (s *JurassicCageService) DbGetCage(c *gin.Context, cageId uuid.UUID) (models.Cage, error) {
	cage, err := s.db.GetCage(c, cageId)
	if err != nil && err == sql.ErrNoRows {
		msg := "no cage found on local db"
		s.logger.Error().Err(err).Msg(msg)
		return models.Cage{}, err
	}
	if err != nil {
		msg := "failed to retrieve cage from local db"
		s.logger.Error().Err(err).Msg(msg)
		return models.Cage{}, err
	}
	return models.Cage{
		ID:                     cage.ID,
		Name:                   cage.Name.String,
		Status:                 models.Status(cage.Status.Int32),
		PredominateEatingHabit: models.EatingHabit(cage.PredominateEatingHabit.Int32),
	}, nil
}

func (s *JurassicCageService) DbGetAllCages(c *gin.Context) ([]models.Cage, error) {
	cages, err := s.db.GetCages(c)
	if err != nil && err == sql.ErrNoRows {
		msg := "no cages found on local db"
		s.logger.Error().Err(err).Msg(msg)
		return []models.Cage{}, err
	}
	if err != nil {
		msg := "failed to retrieve cages from local db"
		s.logger.Error().Err(err).Msg(msg)
		return []models.Cage{}, err
	}

	var retCages []models.Cage
	for _, cage := range cages {
		dinos, err := s.DbGetDinosaursByCageId(c, cage.ID)
		if err != nil {
			msg := fmt.Sprintf("failed to retrieve dinos for cage %s from local db", cage.ID.String())
			s.logger.Warn().Err(err).Msg(msg)
		}

		retCages = append(retCages, models.Cage{
			ID:                     cage.ID,
			Name:                   cage.Name.String,
			Status:                 models.Status(cage.Status.Int32),
			PredominateEatingHabit: models.EatingHabit(cage.PredominateEatingHabit.Int32),
			Capacity:               int32(len(dinos)),
			Dinosaurs:              dinos,
		})
	}
	return retCages, nil
}

func (s *JurassicCageService) DbGetDinosaur(c *gin.Context, dinosaurId uuid.UUID) (models.Dinosaur, error) {
	dinosaur, err := s.db.GetDinosaurAndCage(c, dinosaurId)
	if err != nil && err == sql.ErrNoRows {
		msg := "dinosaur not found"
		s.logger.Error().Err(err).Msg(msg)
		return models.Dinosaur{}, err
	}
	if err != nil {
		msg := "error retrieving dinosaur from local db"
		s.logger.Error().Err(err).Msg(msg)
		return models.Dinosaur{}, err
	}
	return models.Dinosaur{
		ID:          dinosaur.ID,
		Name:        dinosaur.Name.String,
		Species:     models.Species(dinosaur.Species.Int32),
		EatingHabit: models.EatingHabit(dinosaur.EatingHabit.Int32),
		CageId:      dinosaur.CageID.UUID,
		Cage: models.Cage{
			ID:   dinosaur.ID_2.UUID,
			Name: dinosaur.Name_2.String,
		},
	}, nil
}

//func (s *JurassicCageService) DbCreateDinosaur(c *gin.Context)

func (s *JurassicCageService) DbGetAllDinosaurs(c *gin.Context) ([]models.Dinosaur, error) {
	dinosaurs, err := s.db.GetDinosaurs(c)
	if err != nil && err == sql.ErrNoRows {
		msg := "no dinosaurs found on local db"
		s.logger.Error().Err(err).Msg(msg)
		return []models.Dinosaur{}, err
	}
	if err != nil {
		msg := "failed to retrieve dinosaurs from local db"
		s.logger.Error().Err(err).Msg(msg)
		return []models.Dinosaur{}, err
	}

	var retDinos []models.Dinosaur
	for _, dino := range dinosaurs {
		cage, err := s.DbGetCage(c, dino.CageID.UUID)
		if err != nil {
			s.logger.Warn().Err(err).Msgf("failed to retrieve cage for dino %s", dino.ID.String())
		}

		retDinos = append(retDinos, models.Dinosaur{
			ID:          dino.ID,
			Name:        dino.Name.String,
			EatingHabit: models.EatingHabit(dino.EatingHabit.Int32),
			Species:     models.Species(dino.Species.Int32),
			CageId:      dino.CageID.UUID,
			Cage:        cage,
		})
	}
	return retDinos, nil
}

func (s *JurassicCageService) DbGetDinosaursByCageId(c *gin.Context, cageId uuid.UUID) ([]models.Dinosaur, error) {
	dinosaurs, err := s.db.GetDinosaursByCage(c, uuid.NullUUID{cageId, internal.ValidateUuid(cageId)})
	if err != nil && err == sql.ErrNoRows {
		msg := "dinosaurs not found"
		s.logger.Error().Err(err).Msg(msg)
		return []models.Dinosaur{}, err
	}
	if err != nil {
		msg := "error retrieving dinosaurs from local db"
		s.logger.Error().Err(err).Msg(msg)
		return []models.Dinosaur{}, err
	}

	var dinosaursRet []models.Dinosaur
	for _, dinosaur := range dinosaurs {
		cage, err := s.DbGetCage(c, dinosaur.CageID.UUID)
		if err != nil {
			s.logger.Warn().Err(err).Msgf("failed to retrieve cage for dino %s", dinosaur.ID.String())
		}
		dinosaursRet = append(dinosaursRet, models.Dinosaur{
			ID:          dinosaur.ID,
			Name:        dinosaur.Name.String,
			Species:     models.Species(dinosaur.Species.Int32),
			EatingHabit: models.EatingHabit(dinosaur.EatingHabit.Int32),
			CageId:      cage.ID,
			Cage: models.Cage{
				ID:   cage.ID,
				Name: cage.Name,
			},
		})
	}

	return dinosaursRet, nil
}

func (s *JurassicCageService) DbGetCageDinosaurCount(c *gin.Context, cageId uuid.UUID) (int64, error) {
	dinoCount, err := s.db.GetCageDinosaurCount(c, uuid.NullUUID{cageId, internal.ValidateUuid(cageId)})
	if err != nil && err == sql.ErrNoRows {
		msg := "no cage found on local db"
		s.logger.Error().Err(err).Msg(msg)
		return 0, err
	}
	if err != nil {
		msg := "failed to retrieve cage from local db"
		s.logger.Error().Err(err).Msg(msg)
		return 0, err
	}
	return dinoCount, nil
}

func (s *JurassicCageService) DbAddDinosaurToCage(c *gin.Context, cageId uuid.UUID, dino models.Dinosaur) error {
	_, err := s.db.UpdateDinosaurCage(c, postgres.UpdateDinosaurCageParams{
		ID:     dino.ID,
		CageID: uuid.NullUUID{cageId, internal.ValidateUuid(cageId)},
	})
	if err != nil {
		msg := fmt.Sprintf("failed to update cage_id on dinosaur %s", dino.ID.String())
		s.logger.Error().Err(err).Msg(msg)
		return err
	}

	err = s.db.UpdateCagePredominateEatingHabit(c, postgres.UpdateCagePredominateEatingHabitParams{
		PredominateEatingHabit: sql.NullInt32{dino.EatingHabit.Int32(), internal.ValidateInt32(dino.EatingHabit.Int32())},
		ID:                     cageId,
	})
	if err != nil {
		msg := fmt.Sprintf("failed to update predominate_eating_habit on cage %s", cageId.String())
		s.logger.Error().Err(err).Msg(msg)
		return err
	}

	return nil
}

func (s *JurassicCageService) DbRemoveDinosaurFromCage(c *gin.Context, oldCageId uuid.UUID, dinoId uuid.UUID) error {
	_, err := s.db.UpdateDinosaurCage(c, postgres.UpdateDinosaurCageParams{
		ID:     dinoId,
		CageID: uuid.NullUUID{uuid.Nil, internal.ValidateUuid(uuid.Nil)},
	})
	if err != nil {
		msg := fmt.Sprintf("failed to update cage_id on dinosaur %s", dinoId.String())
		s.logger.Error().Err(err).Msg(msg)
		return err
	}

	cage, err := s.DbGetCage(c, oldCageId)
	if err != nil {
		msg := "failed to get cage from local db"
		s.logger.Error().Err(err).Msg(msg)
		return err
	}

	cageCount, err := s.DbGetCageDinosaurCount(c, oldCageId)
	if err != nil {
		msg := "failed to get cage dino count from local db"
		s.logger.Error().Err(err).Msg(msg)
		return err
	}

	s.logger.Info().Msgf("cageCount %d", cageCount)

	if cageCount == 0 {
		eatingHabit := models.UNKNOWN.Int32()
		err = s.db.UpdateCagePredominateEatingHabit(c, postgres.UpdateCagePredominateEatingHabitParams{
			PredominateEatingHabit: sql.NullInt32{eatingHabit, internal.ValidateInt32(eatingHabit)},
			ID:                     oldCageId,
		})
		if err != nil {
			msg := fmt.Sprintf("failed to update predominate_eating_habit on cage %s", oldCageId.String())
			s.logger.Error().Err(err).Msg(msg)
			return err
		}
	} else {
		eatingHabit := cage.PredominateEatingHabit.Int32()
		err = s.db.UpdateCagePredominateEatingHabit(c, postgres.UpdateCagePredominateEatingHabitParams{
			PredominateEatingHabit: sql.NullInt32{eatingHabit, internal.ValidateInt32(eatingHabit)},
			ID:                     oldCageId,
		})
		if err != nil {
			msg := fmt.Sprintf("failed to update predominate_eating_habit on cage %s", oldCageId.String())
			s.logger.Error().Err(err).Msg(msg)
			return err
		}

	}

	return nil
}
