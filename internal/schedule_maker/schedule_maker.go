package schedule_maker

import (
	"fmt"
	"go.uber.org/zap"
	"log"
	"schbot/internal/models"
	"sort"
	"time"
)

type ScheduleMaker struct {
	Log *zap.SugaredLogger
}

func (m *ScheduleMaker) CreateTodaySchedule(lessons []models.Lesson, group string, day int, isNumerator bool) string {
	sort.SliceStable(lessons, func(i, j int) bool {
		return lessons[i].StartAt < lessons[j].StartAt
	})

	var lessonsString = fmt.Sprintf("Группа %s. %s %s\n", group, m.dayOfWeekRus(time.Now().Weekday()), m.currentDate())

	count := 0
	if isNumerator {
		for _, lesson := range lessons {
			if lesson.Day == day && lesson.IsNumerator {
				count += 1
				lessonsString += fmt.Sprintf("%d. %s - %s (%s)\n\t%s\n\tАуд: %s\n", count, lesson.StartAt[:5], lesson.EndAt[:5],
					lesson.Type, lesson.Name, lesson.Cabinet)
			}
		}
	} else {
		for _, lesson := range lessons {
			if lesson.Day == day && !lesson.IsNumerator {
				count += 1
				lessonsString += fmt.Sprintf("%d. %s - %s (%s)\n\t%s\n\tАуд: %s\n", count, lesson.StartAt[:5], lesson.EndAt[:5],
					lesson.Type, lesson.Name, lesson.Cabinet)
			}
		}
	}

	if count == 0 {
		lessonsString += "\n Приказано чилибасить."
	}

	return lessonsString
}

func (m *ScheduleMaker) CreateTomorrowSchedule(lessons []models.Lesson, group string, isNumerator bool) string {
	sort.SliceStable(lessons, func(i, j int) bool {
		return lessons[i].StartAt < lessons[j].StartAt
	})

	tomorrowDate := m.tomorrowDate()
	var lessonsString = fmt.Sprintf("Группа %s. %s %s\n", group, m.dayOfWeekRus(time.Now().AddDate(0, 0, 1).Weekday()), tomorrowDate)

	count := 0
	day := int(time.Now().AddDate(0, 0, 1).Weekday())
	log.Print(day)
	if day == 1 {
		isNumerator = !isNumerator
	}

	if isNumerator {
		for _, lesson := range lessons {
			if lesson.Day == day && lesson.IsNumerator {
				count += 1
				lessonsString += fmt.Sprintf("%d. %s - %s (%s)\n\t%s\n\tАуд: %s\n", count, lesson.StartAt[:5], lesson.EndAt[:5],
					lesson.Type, lesson.Name, lesson.Cabinet)
			}
		}
	} else {
		for _, lesson := range lessons {
			if lesson.Day == day && !lesson.IsNumerator {
				count += 1
				lessonsString += fmt.Sprintf("%d. %s - %s (%s)\n\t%s\n\tАуд: %s\n", count, lesson.StartAt[:5], lesson.EndAt[:5],
					lesson.Type, lesson.Name, lesson.Cabinet)
			}
		}
	}

	if count == 0 {
		lessonsString += "\n Приказано чилибасить."
	}

	return lessonsString
}

func (m *ScheduleMaker) CreateDailySchedule(lessons []models.Lesson, group string, day int) string {
	sort.SliceStable(lessons, func(i, j int) bool {
		return lessons[i].StartAt < lessons[j].StartAt
	})

	var lessonsEven = fmt.Sprintf("Группа %s. \n%s", group, "\U0001F976Числитель:\n")
	var lessonsOdd = "\n\U0001F975Знаменатель:\n"

	countEven := 1
	countOdd := 1

	for _, lesson := range lessons {
		if lesson.Day == day && lesson.IsNumerator {
			countEven += 1
			lessonsEven += fmt.Sprintf("%d. %s - %s (%s)\n\t%s\n\tАуд: %s\n", countEven, lesson.StartAt[:5], lesson.EndAt[:5],
				lesson.Type, lesson.Name, lesson.Cabinet)
		}
		if lesson.Day == day && !lesson.IsNumerator {
			countOdd += 1
			lessonsOdd += fmt.Sprintf("%d. %s - %s (%s)\n\t%s\n\tАуд: %s\n", countOdd, lesson.StartAt[:5], lesson.EndAt[:5],
				lesson.Type, lesson.Name, lesson.Cabinet)
		}
	}

	if countEven == 0 {
		lessonsEven += "\n Приказано чилибасить."
	}

	if countOdd == 0 {
		lessonsOdd += "\n Приказано чилибасить."
	}

	return lessonsEven + lessonsOdd
}

func (m *ScheduleMaker) currentDate() string {
	yy, mm, dd := time.Now().Date()
	if mm < 10 {
		return fmt.Sprintf("%d.%s.%d", dd, fmt.Sprintf("0%d", mm), yy)
	}
	return fmt.Sprintf("%d.%d.%d", dd, mm, yy)
}

func (m *ScheduleMaker) tomorrowDate() string {
	tomorrow := time.Now().AddDate(0, 0, 1)
	yy, mm, dd := tomorrow.Date()
	if mm < 10 {
		return fmt.Sprintf("%d.%s.%d", dd, fmt.Sprintf("0%d", mm), yy)
	}
	return fmt.Sprintf("%d.%d.%d", dd, mm, yy)
}

func (m *ScheduleMaker) dayOfWeekRus(weekday time.Weekday) string {
	switch weekday.String() {
	case "Monday":
		return "Понедельник"
	case "Tuesday":
		return "Вторник"
	case "Wednesday":
		return "Среда"
	case "Thursday":
		return "Четверг"
	case "Friday":
		return "Пятница"
	case "Saturday":
		return "Суббота"
	case "Sunday":
		return "Воскресенье"
	default:
		return ""
	}
}
