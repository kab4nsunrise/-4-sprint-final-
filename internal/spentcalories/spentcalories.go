package spentcalories

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	StepLength       = 0.65
	MInKm            = 1000
	MinInH           = 60
	StepLengthRatio  = 0.45
	WalkingCalorieC  = 0.035
	RunningCalorieC  = 0.029
)

func ParseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, errors.New("invalid data format: expected 'steps,activity,duration'")
	}

	steps, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, "", 0, fmt.Errorf("invalid steps format: %w", err)
	}

	activity := strings.TrimSpace(parts[1])
	activity = strings.ToLower(activity)
	if activity != "running" && activity != "walking" {
		return 0, "", 0, errors.New("unknown training type: expected 'running' or 'walking'")
	}

	duration, err := time.ParseDuration(strings.TrimSpace(parts[2]))
	if err != nil {
		return 0, "", 0, fmt.Errorf("invalid duration format: %w", err)
	}
	if steps <= 0 || duration <= 0 {
		return 0, "", 0, errors.New("steps and duration must be positive")
	}

	return steps, activity, duration, nil
}

func Distance(steps int, height float64) float64 {
	return float64(steps) * StepLengthRatio * height / MInKm
}

func MeanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	distance := Distance(steps, height)
	return distance / duration.Hours()
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := ParseTraining(data)
	if err != nil {
		return "", err
	}

	var calories float64
	switch activity {
	case "running":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
	case "walking":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
	default:
		return "", errors.New("unknown training type")
	}
	if err != nil {
		return "", err
	}

	dist := Distance(steps, height)
	speed := MeanSpeed(steps, height, duration)

	return fmt.Sprintf(
		"Training type: %s\nDuration: %.2f h\nDistance: %.2f km\nSpeed: %.2f km/h\nCalories burned: %.2f kcal",
		strings.ToTitle(activity),
		duration.Hours(),
		dist,
		speed,
		calories,
	), nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("invalid parameters: steps, weight, height, and duration must be positive")
	}
	speed := MeanSpeed(steps, height, duration)
	calories := (RunningCalorieC * weight * duration.Minutes())
	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) float64{
	
	return calories, nil
}
