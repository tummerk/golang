package entities

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/tummerk/golang/schedules/logger"
	"log/slog"
	"testing"
)

func TestScheduleOnDay(t *testing.T) {
	ctx := context.WithValue(context.Background(), "traceID", "test")
	logger := logger.NewLogger("schedule.log", slog.LevelDebug)
	tests := []struct {
		name     string
		schedule Schedule
		expected []int
	}{
		{name: "positive", schedule: Schedule{ReceptionsPerDay: 1}, expected: []int{480}},
		{name: "positive", schedule: Schedule{ReceptionsPerDay: 2}, expected: []int{480, 1320}},
		{name: "positive", schedule: Schedule{ReceptionsPerDay: 3}, expected: []int{480, 900, 1320}},
		{name: "positive", schedule: Schedule{ReceptionsPerDay: 15}, expected: []int{480, 540, 600, 660, 720, 780, 840, 900, 960, 1020, 1080, 1140, 1200, 1260, 1320}},
		{name: "negative", schedule: Schedule{ReceptionsPerDay: 20}, expected: []int{}},
		{name: "negative", schedule: Schedule{ReceptionsPerDay: -1}, expected: []int{}},
	}
	for i, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.schedule.ScheduleOnDay(ctx, logger)
			if !assert.Equal(t, result, test.expected) {
				t.Errorf("Test %d failed: expected %d, got %d", i, test.expected, result)
			}
		})
	}
}
