package models

import (
	"github.com/google/uuid"
	"strings"
)

type Cage struct {
	ID                     uuid.UUID   `json:"id"`
	Name                   string      `json:"name"`
	Status                 Status      `json:"status"`
	Dinosaurs              []Dinosaur  `json:"dinosaurs"`
	PredominateEatingHabit EatingHabit `json:"predominate_eating_habit"`
	Capacity               int32       `json:"capacity"`
}

type CageResponse struct {
	ID                     uuid.UUID          `json:"id"`
	Name                   string             `json:"name"`
	Status                 string             `json:"status"`
	Dinosaurs              []DinosaurResponse `json:"dinosaurs"`
	PredominateEatingHabit string             `json:"predominate_eating_habit"`
	Capacity               int32              `json:"capacity"`
}

type CageRequest struct {
	Status string `json:"status"`
}

type Dinosaur struct {
	ID          uuid.UUID   `json:"id"`
	Name        string      `json:"name"`
	EatingHabit EatingHabit `json:"eating_habit"`
	Species     Species     `json:"species"`
	CageId      uuid.UUID   `json:"cage_id"`
	Cage        Cage        `json:"cage"`
}

type DinosaurResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	EatingHabit string    `json:"eating_habit"`
	Species     string    `json:"species"`
	CageId      uuid.UUID `json:"cage_id"`
	CageName    string    `json:"cage_name"`
}

type DinosaurRequest struct {
	Name        string `json:"name"`
	EatingHabit string `json:"eating_habit"`
	Species     string `json:"species"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type Status int32

const (
	ACTIVE Status = iota
	DOWN
)

func (s Status) String() string {
	switch s {
	case ACTIVE:
		return "ACTIVE"
	case DOWN:
		return "DOWN"
	}
	return "DOWN"
}

func (s Status) Int32() int32 {
	switch s {
	case ACTIVE:
		return 0
	case DOWN:
		return 1
	}
	return 1 // Down in the default ** For Safety **
}

func ParseStatus(s string) Status {
	switch strings.ToLower(s) {
	case "active":
		return ACTIVE
	case "down":
		return DOWN
	}
	return ACTIVE
}

type EatingHabit int32

const (
	UNKNOWN EatingHabit = iota
	CARNIVORE
	HERBIVORE
)

func (eh EatingHabit) String() string {
	switch eh {
	case CARNIVORE:
		return "Carnivore"
	case HERBIVORE:
		return "Herbivore"
	}
	return "Unknown"
}

func (eh EatingHabit) Int32() int32 {
	switch eh {
	case CARNIVORE:
		return 1
	case HERBIVORE:
		return 2
	}
	return 0 // UNKNOWN
}

type Species int32

const (
	TYRANNOSAURUS Species = iota
	VELOCIRAPTOR
	SPINOSAURUS
	MEGALOSAURUS
	BRACHIOSAURUS
	STEGOSAURUS
	ANKYLOSAURUS
	TRICERATOPS
	SPECIES_UNKNOWN
)

func (s Species) String() string {
	switch s {
	case TYRANNOSAURUS:
		return "Tyrannosaurus"
	case VELOCIRAPTOR:
		return "Velociraptor"
	case SPINOSAURUS:
		return "Spinosaurus"
	case MEGALOSAURUS:
		return "Megalosaurus"
	case BRACHIOSAURUS:
		return "Brachiosaurus"
	case STEGOSAURUS:
		return "Stegosaurus"
	case ANKYLOSAURUS:
		return "Ankylosaurus"
	case TRICERATOPS:
		return "Triceratops"
	}
	return "Unknown Species"
}

func (s Species) Int32() int32 {
	switch s {
	case TYRANNOSAURUS:
		return 0
	case VELOCIRAPTOR:
		return 1
	case SPINOSAURUS:
		return 2
	case MEGALOSAURUS:
		return 3
	case BRACHIOSAURUS:
		return 4
	case STEGOSAURUS:
		return 5
	case ANKYLOSAURUS:
		return 6
	case TRICERATOPS:
		return 7
	}
	return 8 // SPECIES_UNKNOW
}

func ParseSpecies(species string) Species {
	switch strings.ToLower(species) {
	case "tyrannosaurus":
		return TYRANNOSAURUS
	case "velociraptor":
		return VELOCIRAPTOR
	case "spinosaurus":
		return SPINOSAURUS
	case "megalosaurus":
		return MEGALOSAURUS
	case "brachiosaurus":
		return BRACHIOSAURUS
	case "stegosaurus":
		return STEGOSAURUS
	case "ankylosaurus":
		return ANKYLOSAURUS
	case "triceratops":
		return TRICERATOPS
	}
	return SPECIES_UNKNOWN
}

func (s Species) EatingHabit() EatingHabit {
	switch s {
	case TYRANNOSAURUS, VELOCIRAPTOR, SPINOSAURUS, MEGALOSAURUS:
		return CARNIVORE
	case BRACHIOSAURUS, STEGOSAURUS, ANKYLOSAURUS, TRICERATOPS:
		return HERBIVORE
	}
	return UNKNOWN
}

func ParseEatingHabit(s string) string {
	switch s {
	case "brachiosaurus", "stegosaurus", "ankylosaurus", "triceratops":
		return "herbivore"
	case "tyrannosaurus", "velociraptor", "spinosaurus", "megalosaurus":
		return "carnivore"
	}
	return ""
}
