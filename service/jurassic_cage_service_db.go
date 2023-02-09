package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/hayesmp/jurassic-cage-service/internal"
	"github.com/hayesmp/jurassic-cage-service/internal/models"
	"github.com/hayesmp/jurassic-cage-service/postgres"
)

/* Db Methods */

func (s *JurassicCageService) DbGetCage(c context.Context, cageId uuid.UUID) (models.Cage, error) {
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

	if err != nil {
		msg := "failed to retrieve dinos from local db"
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

func (s *JurassicCageService) DbGetCageAndDinosaurs(c context.Context, cageId uuid.UUID) (models.Cage, error) {
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

	dinos, err := s.DbGetDinosaursByCageId(c, cageId)
	if err != nil && err == sql.ErrNoRows {
		msg := "no dinos found on local db"
		s.logger.Error().Err(err).Msg(msg)
		return models.Cage{}, err
	}
	if err != nil {
		msg := "failed to retrieve dinos from local db"
		s.logger.Error().Err(err).Msg(msg)
		return models.Cage{}, err
	}

	return models.Cage{
		ID:                     cage.ID,
		Name:                   cage.Name.String,
		Status:                 models.Status(cage.Status.Int32),
		PredominateEatingHabit: models.EatingHabit(cage.PredominateEatingHabit.Int32),
		Dinosaurs:              dinos,
	}, nil
}

func (s *JurassicCageService) DbGetAllCages(c context.Context, status string) ([]models.Cage, error) {
	var cages []postgres.Cage

	// filter on status
	if len(status) > 0 {
		statusFilter := models.ParseStatus(status)

		filteredCages, err := s.db.GetCagesByStatus(c, sql.NullInt32{statusFilter.Int32(), internal.ValidateInt32(statusFilter.Int32())})
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
		cages = filteredCages
	} else {
		nonFilteredCages, err := s.db.GetCages(c)
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
		cages = nonFilteredCages
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

func (s *JurassicCageService) DbCreateCage(c context.Context, cageRequest models.Cage) (models.Cage, error) {
	upsertCageParams := postgres.UpsertCageParams{
		Name:                   sql.NullString{cageRequest.Name, internal.ValidateString(cageRequest.Name)},
		Status:                 sql.NullInt32{models.ACTIVE.Int32(), internal.ValidateInt32(models.ACTIVE.Int32())},
		PredominateEatingHabit: sql.NullInt32{models.UNKNOWN.Int32(), internal.ValidateInt32(models.UNKNOWN.Int32())},
	}
	cage, err := s.db.UpsertCage(c, upsertCageParams)
	if err != nil {
		msg := "failed to save cage to local db"
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

func (s *JurassicCageService) DbUpdateCage(c context.Context, cageRequest models.Cage) (models.Cage, error) {
	cage, err := s.db.UpsertCage(c, postgres.UpsertCageParams{
		Name:                   sql.NullString{cageRequest.Name, internal.ValidateString(cageRequest.Name)},
		Status:                 sql.NullInt32{cageRequest.Status.Int32(), internal.ValidateInt32(cageRequest.Status.Int32())},
		PredominateEatingHabit: sql.NullInt32{cageRequest.PredominateEatingHabit.Int32(), internal.ValidateInt32(cageRequest.PredominateEatingHabit.Int32())},
	})
	if err != nil {
		msg := "failed to save update to local db"
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

func (s *JurassicCageService) DbDeleteCage(c context.Context, cageId uuid.UUID) error {
	err := s.db.DeleteCage(c, cageId)
	if err != nil {
		msg := "failed to delete cage"
		s.logger.Error().Err(err).Msg(msg)
		return err
	}
	return nil
}

func (s *JurassicCageService) DbGetDinosaur(c context.Context, dinosaurId uuid.UUID) (models.Dinosaur, error) {
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

func (s *JurassicCageService) DbCreateDinosaur(c context.Context, dinoRequest models.DinosaurRequest) (models.Dinosaur, error) {
	species := models.ParseSpecies(dinoRequest.Species)
	eatingHabit := species.EatingHabit()

	upsertDinoParams := postgres.UpsertDinosaurParams{
		Name:        sql.NullString{dinoRequest.Name, internal.ValidateString(dinoRequest.Name)},
		Species:     sql.NullInt32{species.Int32(), internal.ValidateInt32(species.Int32())},
		EatingHabit: sql.NullInt32{eatingHabit.Int32(), internal.ValidateInt32(eatingHabit.Int32())},
	}

	dinosaur, err := s.db.UpsertDinosaur(c, upsertDinoParams)
	if err != nil {
		msg := "failed to save dinosaur to local db"
		s.logger.Error().Err(err).Msg(msg)
		return models.Dinosaur{}, err
	}
	return models.Dinosaur{
		ID:          dinosaur.ID,
		Name:        dinosaur.Name.String,
		EatingHabit: models.EatingHabit(dinosaur.EatingHabit.Int32),
		Species:     models.Species(dinosaur.Species.Int32),
		CageId:      dinosaur.CageID.UUID,
	}, nil
}

func (s *JurassicCageService) DbGetAllDinosaurs(c context.Context, species string) ([]models.Dinosaur, error) {
	var dinosaurs []postgres.Dinosaur

	// filter on species if the parameter is passed
	if len(species) > 0 {
		speciesFilter := models.ParseSpecies(species)

		filterDinosaurs, err := s.db.GetDinosaursBySpecies(c, sql.NullInt32{speciesFilter.Int32(), internal.ValidateInt32(speciesFilter.Int32())})
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
		dinosaurs = filterDinosaurs
	} else {
		notFilterDinosaurs, err := s.db.GetDinosaurs(c)
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
		dinosaurs = notFilterDinosaurs
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

func (s *JurassicCageService) DbGetDinosaursByCageId(c context.Context, cageId uuid.UUID) ([]models.Dinosaur, error) {
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

func (s *JurassicCageService) DbGetCageDinosaurCount(c context.Context, cageId uuid.UUID) (int64, error) {
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

func (s *JurassicCageService) DbAddDinosaurToCage(c context.Context, cageId uuid.UUID, dino models.Dinosaur) error {
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

func (s *JurassicCageService) DbRemoveDinosaurFromCage(c context.Context, oldCageId uuid.UUID, dinoId uuid.UUID) error {
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

func (s *JurassicCageService) DbDeleteDinosaur(c context.Context, dinoId uuid.UUID) error {
	err := s.db.DeleteDinosaur(c, dinoId)
	if err != nil {
		msg := "failed to delete dinosaur"
		s.logger.Error().Err(err).Msg(msg)
		return err
	}
	return nil
}
