package handles

import (
	"fmt"
	"schbot/internal/models"
	"sort"
	"time"
)

func (h *Handles) createTodaySchedule(lessons []models.Lesson, group string, day int, is_numerator bool) string {
	sort.SliceStable(lessons, func(i, j int) bool {
		return lessons[i].StartAt < lessons[j].StartAt
	})

	var lessonsstr = fmt.Sprintf("Группа %s. %s %s\n", group, dayOfWeekRus(), date())

	count := 1
	if is_numerator {
		for _, lesson := range lessons {
			if lesson.Day == day && lesson.IsNumerator {
				lessonsstr += fmt.Sprintf("%d. %s - %s (%s)\n\t%s\n\tАуд: %s\n", count, lesson.StartAt[:5], lesson.EndAt[:5],
					lesson.Type, lesson.Name, lesson.Cabinet)
				count = count + 1
			}
		}
	} else {
		for _, lesson := range lessons {
			if lesson.Day == day && !lesson.IsNumerator {
				lessonsstr += fmt.Sprintf("%d. %s - %s (%s)\n\t%s\n\tАуд: %s\n", count, lesson.StartAt[:5], lesson.EndAt[:5],
					lesson.Type, lesson.Name, lesson.Cabinet)
				count = count + 1
			}
		}
	}

	return lessonsstr
}

func (h *Handles) createDailySchedule(lessons []models.Lesson, group string, day int) string {
	sort.SliceStable(lessons, func(i, j int) bool {
		return lessons[i].StartAt < lessons[j].StartAt
	})

	var lessons_even = fmt.Sprintf("Группа %s. %s %s\n%s", group, dayOfWeekRus(), date(), "\U0001F976Числитель:\n")
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

func date() string {
	yy, mm, dd := time.Now().Date()
	if mm < 10 {
		return fmt.Sprintf("%d.%s.%d", dd, fmt.Sprintf("0%d", mm), yy)
	}
	return fmt.Sprintf("%d.%d.%d", dd, mm, yy)
}

func dayOfWeekRus() string {
	switch time.Now().Weekday().String() {
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
