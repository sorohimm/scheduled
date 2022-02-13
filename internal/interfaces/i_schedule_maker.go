package interfaces

import "schbot/internal/models"

type IScheduleMaker interface {
	CreateTodaySchedule(lessons []models.Lesson, group string, day int, is_numerator bool) string
	CreateTomorrowSchedule(lessons []models.Lesson, group string, is_numerator bool) string
	CreateDailySchedule(lessons []models.Lesson, group string, day int) string
}
