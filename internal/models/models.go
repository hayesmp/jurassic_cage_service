package models

import (
	"github.com/google/uuid"
)

type Cage struct {
	ID     uuid.UUID
	Name   string
	Status Status
}

type Dinosaur struct {
	ID          uuid.UUID
	Name        string
	EatingHabit EatingHabit
	Species     Species
	CageID      uuid.NullUUID
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

type EatingHabit int32

const (
	CARNIVORE EatingHabit = iota
	HERBIVORE
	UNKNOWN
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

func (s Species) EatingHabit() EatingHabit {
	switch s {
	case TYRANNOSAURUS, VELOCIRAPTOR, SPINOSAURUS, MEGALOSAURUS:
		return CARNIVORE
	case BRACHIOSAURUS, STEGOSAURUS, ANKYLOSAURUS, TRICERATOPS:
		return HERBIVORE
	}
	return UNKNOWN
}
