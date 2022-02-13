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

func (m *ScheduleMaker) CreateTodaySchedule(lessons []models.Lesson, group string, day int, is_numerator bool) string {
	sort.SliceStable(lessons, func(i, j int) bool {
		return lessons[i].StartAt < lessons[j].StartAt
	})

	var lessonsstr = fmt.Sprintf("Группа %s. %s %s\n", group, m.dayOfWeekRus(time.Now().Weekday()), m.currentDate())

	count := 0
	if is_numerator {
		for _, lesson := range lessons {
			if lesson.Day == day && lesson.IsNumerator {
				count += 1
				lessonsstr += fmt.Sprintf("%d. %s - %s (%s)\n\t%s\n\tАуд: %s\n", count, lesson.StartAt[:5], lesson.EndAt[:5],
					lesson.Type, lesson.Name, lesson.Cabinet)
			}
		}
	} else {
		for _, lesson := range lessons {
			if lesson.Day == day && !lesson.IsNumerator {
				count += 1
				lessonsstr += fmt.Sprintf("%d. %s - %s (%s)\n\t%s\n\tАуд: %s\n", count, lesson.StartAt[:5], lesson.EndAt[:5],
					lesson.Type, lesson.Name, lesson.Cabinet)
			}
		}
	}

	if count == 0 {
		lessonsstr += "\n Приказано чилибасить."
	}

	return lessonsstr
}

func (m *ScheduleMaker) CreateTomorrowSchedule(lessons []models.Lesson, group string, is_numerator bool) string {
	sort.SliceStable(lessons, func(i, j int) bool {
		return lessons[i].StartAt < lessons[j].StartAt
	})

	tmdate := m.tomorrowDate()
	var lessonsstr = fmt.Sprintf("Группа %s. %s %s\n", group, m.dayOfWeekRus(time.Now().AddDate(0, 0, 1).Weekday()), tmdate)

	count := 0
	day := int(time.Now().AddDate(0, 0, 1).Weekday())
	log.Print(day)
	if day == 1 {
		is_numerator = !is_numerator
	}

	if is_numerator {
		for _, lesson := range lessons {
			if lesson.Day == day && lesson.IsNumerator {
				count += 1
				lessonsstr += fmt.Sprintf("%d. %s - %s (%s)\n\t%s\n\tАуд: %s\n", count, lesson.StartAt[:5], lesson.EndAt[:5],
					lesson.Type, lesson.Name, lesson.Cabinet)
			}
		}
	} else {
		for _, lesson := range lessons {
			if lesson.Day == day && !lesson.IsNumerator {
				count += 1
				lessonsstr += fmt.Sprintf("%d. %s - %s (%s)\n\t%s\n\tАуд: %s\n", count, lesson.StartAt[:5], lesson.EndAt[:5],
					lesson.Type, lesson.Name, lesson.Cabinet)
			}
		}
	}

	if count == 0 {
		lessonsstr += "\n Приказано чилибасить."
	}

	return lessonsstr
}

func (m *ScheduleMaker) CreateDailySchedule(lessons []models.Lesson, group string, day int) string {
	sort.SliceStable(lessons, func(i, j int) bool {
		return lessons[i].StartAt < lessons[j].StartAt
	})

	var lessons_even = fmt.Sprintf("Группа %s. \n%s", group, "\U0001F976Числитель:\n")
	var lessons_odd = "\n\U0001F975Знаменатель:\n"

	count_e := 1
	count_o := 1

	for _, lesson := range lessons {
		if lesson.Day == day && lesson.IsNumerator {
			lessons_even += fmt.Sprintf("%d. %s - %s (%s)\n\t%s\n\tАуд: %s\n", count_e, lesson.StartAt[:5], lesson.EndAt[:5],
				lesson.Type, lesson.Name, lesson.Cabinet)
			count_e = count_e + 1
		}
		if lesson.Day == day && !lesson.IsNumerator {
			lessons_odd += fmt.Sprintf("%d. %s - %s (%s)\n\t%s\n\tАуд: %s\n", count_o, lesson.StartAt[:5], lesson.EndAt[:5],
				lesson.Type, lesson.Name, lesson.Cabinet)
			count_o = count_o + 1
		}
	}

	return lessons_even + lessons_odd
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
