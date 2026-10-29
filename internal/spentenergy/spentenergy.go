package spentenergy

import (
	"errors"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе.
)

// Статические ошибки для валидции входных данных.
var (
	ErrInvalidSteps    = errors.New("недопустимое количество шагов: должно быть > 0")
	ErrInvalidWeight   = errors.New("недопустимый вес (кг): должен быть > 0")
	ErrInvalidHeight   = errors.New("недопустимый рост (м): должен быть > 0")
	ErrInvalidDuration = errors.New("недопустимая продолжительность: должна быть > 0")
)

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	cal, err := baseCalories(steps, weight, height, duration)
	if err != nil {
		return 0, err
	}
	cal *= walkingCaloriesCoefficient
	return cal, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	return baseCalories(steps, weight, height, duration)
}

func MeanSpeed(steps int, height float64, duration time.Duration) float64 {
	if steps < 0 {
		return 0
	}

	if duration <= 0 {
		return 0
	}

	distance := Distance(steps, height)
	hours := duration.Hours()
	if hours == 0 {
		return 0
	}
	return distance / hours
}

func Distance(steps int, height float64) float64 {
	stepLen := height * stepLengthCoefficient
	meters := float64(steps) * stepLen

	return meters / float64(mInKm)
}

func baseCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, ErrInvalidSteps
	}
	if weight <= 0 {
		return 0, ErrInvalidWeight
	}
	if height <= 0 {
		return 0, ErrInvalidHeight
	}
	if duration <= 0 {
		return 0, ErrInvalidDuration
	}

	avg := MeanSpeed(steps, height, duration)
	durMin := duration.Minutes()
	cal := (weight * avg * durMin) / float64(minInH)
	return cal, nil
}
