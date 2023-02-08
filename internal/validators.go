package internal

import (
	"github.com/google/uuid"
	"time"
)

func ValidateInt32(intValue int32) bool {
	var intPointer *int32
	intPointer = &intValue
	if intPointer == nil {
		return false
	}
	return true
}

func ValidateBool(boolValue bool) bool {
	var boolPointer *bool
	boolPointer = &boolValue
	if boolPointer == nil {
		return false
	}
	return true
}

func ValidateString(stringValue string) bool {
	var stringPointer *string
	stringPointer = &stringValue
	if stringPointer == nil {
		return false
	}
	return true
}

func ValidateTime(timeValue time.Time) bool {
	var timePointer *time.Time
	timePointer = &timeValue
	if timePointer == nil {
		return false
	}
	return true
}

func ValidateUuid(uuidValue uuid.UUID) bool {
	var uuidPointer *uuid.UUID
	uuidPointer = &uuidValue
	if uuidPointer == nil {
		return false
	}
	return true
}
