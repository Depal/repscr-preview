package utils

import (
	"testing"
	"time"
)

type FormatMap struct {
	moment    time.Time
	shiftName string
}

type OffsetMap struct {
	moment    time.Time
	offset    int
	shiftName string
}

func TestTimeToShiftName(t *testing.T) {
	tests := []FormatMap{
		{time.Date(2021, 1, 1, 1, 0, 0, 0, time.Local), "31.12.2020_ночная"},
		{time.Date(2021, 1, 5, 3, 0, 0, 0, time.Local), "04.01.2021_ночная"},
		{time.Date(2021, 2, 5, 11, 30, 12, 0, time.Local), "05.02.2021_дневная"},
		{time.Date(2021, 3, 6, 23, 10, 10, 0, time.Local), "06.03.2021_ночная"},
		{time.Date(2021, 4, 7, 9, 0, 0, 0, time.Local), "07.04.2021_дневная"},
		{time.Date(2021, 5, 8, 21, 0, 0, 0, time.Local), "08.05.2021_ночная"},
	}

	for _, test := range tests {
		moment := test.moment
		expected := test.shiftName
		got := TimeToShiftName(moment)

		if got != expected {
			t.Logf("FAIL (expected %v; got %v)", expected, got)
			t.Fail()
		} else {
			t.Logf("%v OK", got)
		}
	}
}

func TestShiftOffsetToShiftName(t *testing.T) {
	tests := []OffsetMap{
		{time.Date(2021, 1, 30, 15, 0, 0, 0, time.Local), -1, "29.01.2021_ночная"},
		{time.Date(2021, 1, 30, 15, 0, 0, 0, time.Local), -2, "29.01.2021_дневная"},
		{time.Date(2021, 1, 30, 15, 0, 0, 0, time.Local), -3, "28.01.2021_ночная"},
		{time.Date(2021, 1, 30, 15, 0, 0, 0, time.Local), 1, "30.01.2021_ночная"},
		{time.Date(2021, 1, 30, 15, 0, 0, 0, time.Local), 2, "31.01.2021_дневная"},
		{time.Date(2021, 1, 30, 15, 0, 0, 0, time.Local), 0, "30.01.2021_дневная"},
	}

	for _, test := range tests {
		moment := test.moment
		offset := test.offset
		expected := test.shiftName
		got := ShiftOffsetToShiftName(moment, offset)

		if got != expected {
			t.Logf("FAIL (expected %v; got %v)", expected, got)
			t.Fail()
		} else {
			t.Logf("%v OK", got)
		}
	}
}

type shiftStartTest struct {
	lastSentAt      time.Time
	now             time.Time
	newShiftStarted bool
}

func TestNewShiftStarted(t *testing.T) {
	tests := []shiftStartTest{
		{
			time.Date(2021, 2, 2, 4, 0, 0, 0, time.Local),
			time.Date(2021, 2, 2, 15, 0, 0, 0, time.Local),
			true,
		},

		{
			time.Date(2021, 2, 2, 16, 0, 0, 0, time.Local),
			time.Date(2021, 2, 2, 22, 0, 0, 0, time.Local),
			true,
		},

		{
			time.Date(2020, 12, 31, 23, 0, 0, 0, time.Local),
			time.Date(2021, 1, 1, 12, 0, 0, 0, time.Local),
			true,
		},

		{
			time.Date(2020, 12, 31, 23, 0, 0, 0, time.Local),
			time.Date(2021, 1, 1, 2, 0, 0, 0, time.Local),
			false,
		},

		{
			time.Date(2021, 2, 2, 4, 0, 0, 0, time.Local),
			time.Date(2021, 2, 2, 9, 59, 0, 0, time.Local),
			false,
		},

		{
			time.Date(2021, 2, 2, 16, 0, 0, 0, time.Local),
			time.Date(2021, 2, 2, 21, 15, 0, 0, time.Local),
			false,
		},

		{
			time.Date(2021, 3, 3, 4, 0, 0, 0, time.Local),
			time.Date(2021, 3, 3, 6, 0, 0, 0, time.Local),
			false,
		},
		{
			time.Date(2021, 4, 4, 21, 0, 0, 0, time.Local),
			time.Date(2021, 4, 5, 22, 0, 0, 0, time.Local),
			true,
		},

		{
			time.Date(2021, 5, 5, 18, 0, 0, 0, time.Local),
			time.Date(2021, 5, 5, 22, 0, 1, 0, time.Local),
			true,
		},

		{
			time.Date(2021, 5, 5, 6, 0, 0, 0, time.Local),
			time.Date(2021, 5, 5, 10, 0, 0, 0, time.Local),
			true,
		},

		{
			time.Date(2021, 6, 6, 8, 0, 0, 0, time.Local),
			time.Date(2021, 6, 6, 8, 0, 0, 0, time.Local),
			false,
		},
	}

	for _, test := range tests {
		result := NewShiftStarted(test.now, test.lastSentAt)
		if test.newShiftStarted != result {
			t.Logf("failed for %v -> %v (got %v, expected %v)", test.lastSentAt, test.now, result, test.newShiftStarted)
			t.Fail()
		}
	}
}
