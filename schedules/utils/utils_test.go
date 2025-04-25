package utils

import (
	"strconv"
	"testing"
)

func TestRoundUp(t *testing.T) {
	tests := []struct {
		num, multiple int
		expected      int
	}{
		{num: 0, multiple: 15, expected: 0},
		{num: 1, multiple: 15, expected: 15},
		{num: 2, multiple: 15, expected: 15},
		{num: 100000000, multiple: 15, expected: 100000005},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			result := RoundUp(test.num, test.multiple)
			if result != test.expected {
				t.Errorf("Test %d failed: expected %d, got %d", i, test.expected, result)
			}
		})
	}
}

func TestMinuteToTime(t *testing.T) {
	tests := []struct {
		min      int
		expected string
	}{
		{min: 0, expected: "00:00"},
		{min: 60, expected: "01:00"},
		{min: 1440, expected: "24:00"},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			result := MinuteToTime(test.min)
			if result != test.expected {
				t.Errorf("Test %d failed: expected %s, got %s", i, test.expected, result)
			}
		})
	}
}
