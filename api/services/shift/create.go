package shift

import (
	"context"
	"time"
)

func (s *Shift) CreateNewShiftSchedule(ctx context.Context, roleId int, startTime, endTime time.Time) (int, error) {
	return s.storage.CreateNewShiftSchedule(ctx, roleId, startTime, endTime)
}
