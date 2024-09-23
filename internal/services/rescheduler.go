package services

import (
	"errors"
	"slices"
	"strconv"
	"strings"
	"time"
)

// Переназначает задачу на следующую дату.
func NextDate(now time.Time, date string, repeat string) (string, error) {
	var result string

	dateParsed, err := time.Parse("20060102", date)
	if err != nil {
		return "", err
	}

	// Проверка на длину правила
	repeatArray := strings.Split(repeat, " ")
	if repeatArray[0] != "y" && (len(repeatArray) < 2 || len(repeatArray) > 3) {
		return "", errors.New("некорректный формат правила")
	}

	// Проверка информации о месяцах для переноса по месяцам
	repeatMonths := "1,2,3,4,5,6,7,8,9,10,11,12"
	if len(repeatArray) == 3 {
		repeatMonths = repeatArray[2]
	}

	switch repeatArray[0] {
	case "y":
		result, err = repeat_y(now, dateParsed)
	case "d":
		result, err = repeat_d(now, dateParsed, repeatArray[1])
	case "w":
		result, err = repeat_w(now, dateParsed, repeatArray[1])
	case "m":
		result, err = repeat_m(now, dateParsed, repeatArray[1], repeatMonths)
	default:
		return "", errors.New("правило не поддерживается")
	}

	return result, err
}

func repeat_y(now time.Time, dateParsed time.Time) (string, error) {
	// Первое переназначение необходимо даже в случае, если дата уже больше, чем now:
	dateParsed = dateParsed.AddDate(1, 0, 0)

	// При необходимости переназначает задачу, пока дата не будет больше, чем now:
	for !dateParsed.After(now) {
		dateParsed = dateParsed.AddDate(1, 0, 0)
	}

	return dateParsed.Format("20060102"), nil
}

func repeat_d(now time.Time, dateParsed time.Time, repeat string) (string, error) {
	// Конвертация строковой информации о днях в число
	days, err := strconv.Atoi(repeat)
	if err != nil {
		return "", err
	}

	// Максимально допустимое число равно 400
	if days > 400 {
		return "", errors.New("превышено максимально допустимое число дней")
	}

	// Первое переназначение необходимо даже в случае, если дата уже больше, чем now:
	dateParsed = dateParsed.AddDate(0, 0, days)

	// При необходимости переназначает задачу, пока дата не будет больше, чем now:
	for !dateParsed.After(now) {
		dateParsed = dateParsed.AddDate(0, 0, days)
	}

	return dateParsed.Format("20060102"), nil
}

func repeat_w(now time.Time, dateParsed time.Time, repeat string) (string, error) {
	// Переводит строковую информацию о днях в последовательность чисел
	daysStringArray := strings.Split(repeat, ",")
	var daysArray []int

	// Проверяет, что все элементы массива - корректные дни недели
	for _, d := range daysStringArray {
		n, err := strconv.Atoi(d)
		if err != nil || n < 1 || n > 7 {
			return "", errors.New("недопустимый день недели")
		}
		if n == 7 {
			n = 0 // Для соответствия формату int(dateParsed.Weekday())
		}
		daysArray = append(daysArray, n)
	}

	// Переназначает задачу на соответствующий день недели после now
	invalidDate := true
	for !dateParsed.After(now) || invalidDate {
		dateParsed = dateParsed.AddDate(0, 0, 1)

		weekday := int(dateParsed.Weekday())
		if slices.Contains(daysArray, weekday) {
			invalidDate = false
		} else {
			invalidDate = true
		}
	}

	return dateParsed.Format("20060102"), nil
}

func repeat_m(now time.Time, dateParsed time.Time, repeatDays string, repeatMonths string) (string, error) {
	// Переводит строковую информацию о днях в последовательность чисел
	daysStringArray := strings.Split(repeatDays, ",")
	var daysArray []int

	// Проверяет, что все элементы массива - корректные числа месяца
	for _, d := range daysStringArray {
		n, err := strconv.Atoi(d)
		if err != nil || n < -2 || n > 31 || n == 0 {
			return "", errors.New("недопустимый день месяца")
		}
		daysArray = append(daysArray, n)
	}

	var monthsArray []int
	monthsStringArray := strings.Split(repeatMonths, ",")

	// Проверяет, что все элементы массива - корректные месяцы
	for _, d := range monthsStringArray {
		n, err := strconv.Atoi(d)
		if err != nil || n < 1 || n > 12 {
			return "", errors.New("недопустимый день месяца")
		}
		monthsArray = append(monthsArray, n)
	}

	// Проверка случая, когда укзаны числа, которые отсутствуют в феврале
	if len(monthsArray) == 1 && monthsArray[0] == 2 {
		invalidDate := true
		for _, d := range daysArray {
			if d < 30 {
				invalidDate = false
				break
			}
		}
		if invalidDate {
			return "", errors.New("недопустимый день месяца")
		}
	}

	// Проверка случая, когда укзаны числа, которые отсутствуют во всех нужных месяцах
	if len(daysArray) == 1 && daysArray[0] == 31 {
		invalidDate := true

		for _, m := range monthsArray {
			if slices.Contains([]int{1, 3, 5, 7, 8, 10, 12}, m) {
				invalidDate = false
				break
			}
		}

		if invalidDate {
			return "", errors.New("недопустимый день месяца")
		}
	}

	// Переназначает задачу на соответствующий день соответствующего месяца после now
	invalidDate := true
	for !dateParsed.After(now) || invalidDate {
		dateParsed = dateParsed.AddDate(0, 0, 1)

		monthday := int(dateParsed.Day())
		month := int(dateParsed.Month())
		var convertedDaysArray []int
		for _, d := range daysArray {
			// Конвертация дней вида -1, -2 в конкретные даты
			if d == -1 {
				daysInMonth := time.Date(dateParsed.Year(), dateParsed.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
				convertedDaysArray = append(convertedDaysArray, daysInMonth)
			} else if d == -2 {
				daysInMonth := time.Date(dateParsed.Year(), dateParsed.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
				convertedDaysArray = append(convertedDaysArray, daysInMonth-1)
			} else {
				convertedDaysArray = append(convertedDaysArray, d)
			}
		}

		if slices.Contains(convertedDaysArray, monthday) && slices.Contains(monthsArray, month) {
			invalidDate = false
		} else {
			invalidDate = true
		}
	}

	return dateParsed.Format("20060102"), nil
}
