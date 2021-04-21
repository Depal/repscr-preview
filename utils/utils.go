package utils

import (
	"time"
)

const (
	SwitchNightToDay = 9  // hour
	SwitchDayToNight = 21 // hour
	MaxShiftDuration = 12 // hours
)

func TimeToShiftName(moment time.Time) (shiftName string) {
	if moment.Hour() < SwitchNightToDay {
		moment = moment.AddDate(0, 0, -1)
		return moment.Format("02.01.2006_ночная")
	} else if moment.Hour() < SwitchDayToNight {
		return moment.Format("02.01.2006_дневная")
	} else {
		return moment.Format("02.01.2006_ночная")
	}
}

func ShiftOffsetToShiftName(moment time.Time, offset int) (shiftName string) {
	durationOffset := time.Duration(12*offset) * time.Hour
	moment = moment.Add(durationOffset)
	return TimeToShiftName(moment)
}

func IsDayShift(moment time.Time) bool {
	hour := moment.Hour()

	return hour > SwitchNightToDay && hour <= SwitchDayToNight
}

func NewShiftStarted(now time.Time, previous time.Time) bool {
	return (now.Hour() > SwitchNightToDay && previous.Hour() <= SwitchNightToDay) ||
		(now.Hour() > SwitchDayToNight && previous.Hour() <= SwitchDayToNight) ||
		(IsDayShift(now) && !IsDayShift(previous)) ||
		now.Sub(previous).Hours() > MaxShiftDuration
}
