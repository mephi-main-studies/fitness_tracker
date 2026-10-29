package daysteps

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	personaldata "github.com/Yandex-Practicum/tracker/internal/personaldata"
	spentenergy "github.com/Yandex-Practicum/tracker/internal/spentenergy"
)

// Статические ошибки пакета daysteps
var (
	ErrInvalidInputFormat   = errors.New("неверный формат ввода")
	ErrInvalidStepsValue    = errors.New("некорректное значение шагов")
	ErrInvalidDurationValue = errors.New("некорректное значение продолжительности")
)

type DaySteps struct {
	Steps    int
	Duration time.Duration
	personaldata.Personal
}

// Parse разбирает строку формата "<steps>,<duration>" и заполняет структуру
func (ds *DaySteps) Parse(datastring string) (err error) {
	parts := strings.Split(datastring, ",")
	if len(parts) != 2 {
		return ErrInvalidInputFormat
	}

	stepsStr := parts[0]
	durStr := parts[1]

	// Строгие требования из тестов: без ведущих/замыкающих пробелов
	if strings.TrimSpace(stepsStr) != stepsStr || strings.TrimSpace(durStr) != durStr || stepsStr == "" || durStr == "" {
		return ErrInvalidInputFormat
	}

	// steps: целое > 0, допускается ведущий '+'
	steps, err := strconv.Atoi(stepsStr)
	if err != nil || steps <= 0 {
		return ErrInvalidStepsValue
	}

	// duration: корректный time.ParseDuration и > 0
	dur, err := time.ParseDuration(durStr)
	if err != nil || dur <= 0 {
		return ErrInvalidDurationValue
	}

	ds.Steps = steps
	ds.Duration = dur
	return nil
}

// ActionInfo формирует строку отчёта о прогулке
func (ds DaySteps) ActionInfo() (string, error) {
	distance := spentenergy.Distance(ds.Steps, ds.Height)
	calories, err := spentenergy.WalkingSpentCalories(ds.Steps, ds.Weight, ds.Height, ds.Duration)
	if err != nil {
		return "", err
	}

	info := fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		ds.Steps,
		distance,
		calories,
	)
	return info, nil
}
