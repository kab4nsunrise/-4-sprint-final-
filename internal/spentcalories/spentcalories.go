package spentcalories

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	StepLengthRatio            = 0.45
	MetersInKm                 = 1000
	MinutesInHour              = 60
	WalkingCaloriesCoefficient = 0.5
)

func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("invalid data format: %q, expected 'steps,activity,duration'", data)
	}

	stepsStr := strings.TrimSpace(parts[0])
	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("invalid steps format: %q, must be integer: %w", stepsStr, err)
	}
	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("steps count must be positive, got %d", steps)
	}

	activity := strings.TrimSpace(parts[1])
	if activity != "Бег" && activity != "Ходьба" {
		return 0, "", 0, fmt.Errorf("unknown training type: %q, expected 'Бег' or 'Ходьба'", activity)
	}

	durationStr := strings.TrimSpace(parts[2])
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("invalid duration format: %q, must be valid time duration: %w", durationStr, err)
	}
	if duration <= 0 {
		return 0, "", 0, fmt.Errorf("duration must be positive, got %v", duration)
	}

	return steps, activity, duration, nil
}

func distance(steps int, height float64) float64 {
	if height <= 0 {
		return 0
	}
	stepLength := height * StepLengthRatio
	distanceMeters := float64(steps) * stepLength
	return distanceMeters / MetersInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 || height <= 0 || steps <= 0 {
		return 0
	}
	dist := distance(steps, height)
	durationHours := duration.Hours()
	if durationHours == 0 {
		return 0
	}
	return dist / durationHours
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("steps must be positive")
	}
	if weight <= 0 {
		return 0, errors.New("weight must be positive")
	}
	if height <= 0 {
		return 0, errors.New("height must be positive")
	}
	if duration <= 0 {
		return 0, errors.New("duration must be positive")
	}

	speed := meanSpeed(steps, height, duration)
	durationMinutes := duration.Minutes()
	if durationMinutes == 0 {
		return 0, errors.New("duration too short")
	}

	calories := (weight * speed * durationMinutes) / MinutesInHour
	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("steps must be positive")
	}
	if weight <= 0 {
		return 0, errors.New("weight must be positive")
	}
	if height <= 0 {
		return 0, errors.New("height must be positive")
	}
	if duration <= 0 {
		return 0, errors.New("duration must be positive")
	}

	speed := meanSpeed(steps, height, duration)
	durationMinutes := duration.Minutes()
	if durationMinutes == 0 {
		return 0, errors.New("duration too short")
	}

	calories := (weight * speed * durationMinutes) / MinutesInHour
	calories *= WalkingCaloriesCoefficient
	return calories, nil
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	if weight <= 0 {
		return "", errors.New("weight must be positive")
	}
	if height <= 0 {
		return "", errors.New("height must be positive")
	}

	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		return "", fmt.Errorf("failed to parse training data: %w", err)
	}

	var calories float64
	switch activity {
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
	default:
		return "", fmt.Errorf("unsupported training type: %q", activity)
	}
	if err != nil {
		return "", fmt.Errorf("calories calculation failed: %w", err)
	}

	dist := distance(steps, height)
	speed := meanSpeed(steps, height, duration)

	return fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожжено калорий: %.2f",
		activity,
		duration.Hours(),
		dist,
		speed,
		calories,
	), nil
}
