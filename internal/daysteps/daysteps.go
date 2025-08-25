package daysteps

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/kab4nsunrise/4-sprint-final--/internal/spentcalories"
)

const (
	stepLength = 0.65 // meters
	mInKm      = 1000 // meters in kilometer
)

func parsePackage(data string) (int, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid data format: %q, expected 'steps,duration'", data)
	}

	stepsStr := strings.TrimSpace(parts[0])
	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid steps format: %q, must be integer: %w", stepsStr, err)
	}
	if steps <= 0 {
		return 0, 0, fmt.Errorf("steps count must be positive, got %d", steps)
	}

	durationStr := strings.TrimSpace(parts[1])
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid duration format: %q, must be valid time duration: %w", durationStr, err)
	}
	if duration <= 0 {
		return 0, 0, fmt.Errorf("duration must be positive, got %v", duration)
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) (string, error) {
	if weight <= 0 {
		return "", errors.New("weight must be positive")
	}
	if height <= 0 {
		return "", errors.New("height must be positive")
	}

	steps, duration, err := parsePackage(data)
	if err != nil {
		return "", fmt.Errorf("failed to parse data: %w", err)
	}

	distanceMeters := float64(steps) * stepLength
	distanceKm := distanceMeters / mInKm

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		return "", fmt.Errorf("calories calculation failed: %w", err)
	}

	return fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.",
		steps, distanceKm, calories,
	), nil
}
