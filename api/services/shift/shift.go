package shift

import (
	"context"
	"time"
)

type storage interface {
	CreateNewShiftSchedule(ctx context.Context, roleId int, startTime, endTime time.Time) (int, error)
}

type ShiftInterface interface {
	CreateNewShiftSchedule(ctx context.Context, roleId int, startTime time.Time, endTime time.Time) (int, error)
}

type Shift struct {
	storage storage
}

func NewShift(storage storage) *Shift {
	return &Shift{
		storage: storage,
	}
}
