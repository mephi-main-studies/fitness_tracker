package trainings

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	personaldata "github.com/Yandex-Practicum/tracker/internal/personaldata"
	spentenergy "github.com/Yandex-Practicum/tracker/internal/spentenergy"
)

// Статические ошибки пакета trainings
var (
	ErrInvalidInputFormat   = errors.New("неверный формат ввода")
	ErrInvalidStepsValue    = errors.New("некорректное значение шагов")
	ErrInvalidDurationValue = errors.New("некорректное значение продолжительности")
	ErrUnknownTrainingType  = errors.New("неизвестный тип тренировки")
)

type Training struct {
	Steps        int
	TrainingType string
	Duration     time.Duration
	personaldata.Personal
}

func (t *Training) Parse(datastring string) (err error) {
	// Ожидаемый формат: "<steps>,<type>,<duration>"
	parts := strings.Split(datastring, ",")
	if len(parts) != 3 {
		return ErrInvalidInputFormat
	}

	stepsStr := strings.TrimSpace(parts[0])
	typeStr := strings.TrimSpace(parts[1])
	durStr := strings.TrimSpace(parts[2])

	// steps: положительное целое число
	steps, err := strconv.Atoi(stepsStr)
	if err != nil || steps <= 0 {
		return ErrInvalidStepsValue
	}

	// duration: положительная продолжительность
	dur, err := time.ParseDuration(durStr)
	if err != nil || dur <= 0 {
		return ErrInvalidDurationValue
	}

	// запись значений в структуру
	t.Steps = steps
	t.TrainingType = typeStr
	t.Duration = dur
	return nil
}

func (t Training) ActionInfo() (string, error) {
	// Расчёт базовых метрик
	distance := spentenergy.Distance(t.Steps, t.Height)
	speed := spentenergy.MeanSpeed(t.Steps, t.Height, t.Duration)

	// Расчёт калорий в зависимости от типа тренировки
	var (
		cal float64
		err error
	)
	switch strings.ToLower(strings.TrimSpace(t.TrainingType)) {
	case "бег":
		cal, err = spentenergy.RunningSpentCalories(t.Steps, t.Weight, t.Height, t.Duration)
	case "ходьба":
		cal, err = spentenergy.WalkingSpentCalories(t.Steps, t.Weight, t.Height, t.Duration)
	default:
		return "", ErrUnknownTrainingType
	}
	if err != nil {
		// Пробрасываем ошибку наверх без изменения текста, чтобы сохранить единый стиль сообщений
		return "", err
	}

	// Формирование строки отчёта
	info := fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		t.TrainingType,
		t.Duration.Hours(),
		distance,
		speed,
		cal,
	)
	return info, nil
}
