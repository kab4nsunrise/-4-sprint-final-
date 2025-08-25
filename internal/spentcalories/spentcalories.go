package spentcalories

import (
  "errors"
  "fmt"
  "strconv"
  "strings"
  "time"
)

const (
  LenStep = 0.65
  MInKm = 1000
  MinInH = 60
  StepLengthCoefficient = 0.45
  WalkingCaloriesCoefficient = 0.5
  RunningCaloriesCoefficient = 0.029
)

func ParseTraining(data string) (int, string, time.Duration, error) {
  parts := strings.Split(data, ",")
  if len(parts) != 3 {
    return 0, "", 0, errors.New("invalid data format")
  }

  steps, err := strconv.Atoi(strings.TrimSpace(parts[0]))
  if err != nil {
    return 0, "", 0, fmt.Errorf("invalid steps format: %w", err)
  }

  activity := strings.TrimSpace(parts[1])
  if activity != "Бег" && activity != "Ходьба" {
    return 0, "", 0, errors.New("unknown training type")
  }

  duration, err := time.ParseDuration(strings.TrimSpace(parts[2]))
  if err != nil {
    return 0, "", 0, fmt.Errorf("invalid duration format: %w", err)
  }

  return steps, activity, duration, nil
}

func Distance(steps int, height float64) float64 {
  return float64(steps) * height * StepLengthCoefficient / MInKm
}

func MeanSpeed(steps int, height float64, duration time.Duration) float64 {
  if duration <= 0 {
    return 0
  }
  return Distance(steps, height) / duration.Hours()
}

func TrainingInfo(data string, weight, height float64) (string, error) {
  steps, activity, duration, err := ParseTraining(data)
  if err != nil {
    return "", err
  }

  var calories float64
  switch activity {
  case "Бег":
    calories, err = RunningSpentCalories(steps, weight, height, duration)
  case "Ходьба":
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
    "Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f",
    activity,
    duration.Hours(),
    dist,
    speed,
    calories,
  ), nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
  if steps <= 0  weight <= 0  height <= 0  duration <= 0 {
    return 0, fmt.Errorf("invalid parameters: steps=%d, weight=%.2f, height=%.2f, duration=%v", steps, weight, height, duration)
  }
  speed := MeanSpeed(steps, height, duration)
  return (weight * speed * duration.Minutes()) / MinInH, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
  if steps <= 0  weight <= 0  height <= 0  duration <= 0 {
    return 0, errors.New("invalid parameters")
  }
  speed := MeanSpeed(steps, height, duration)
  calories := (weight * speed * duration.Minutes()) / MinInH
  return calories * WalkingCaloriesCoefficient, nil
}
