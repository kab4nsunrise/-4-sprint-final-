package spentcalories

import (
	"time"
)

const (
	cmInM         = 100
	mInKm         = 1000
	minInH        = 60
	stepLength    = 0.65
	walkingSpeed  = 4.0
	caloriesRate  = 0.035
	conversionFac = 0.029
)

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) float64 {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0.0
	}

	heightM := height / cmInM
	distance := float64(steps) * stepLength / mInKm
	speed := distance / duration.Hours()
	met := speed / walkingSpeed * 4.0

	calories := met * weight * duration.Hours() * caloriesRate

	return calories
}
