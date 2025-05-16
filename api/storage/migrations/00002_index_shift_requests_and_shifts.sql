-- +goose Up
CREATE INDEX idx_shift_requests_covering_filter ON shift_requests (employee_id, shift_id, status);
CREATE INDEX idx_shifts_covering_filter ON shifts (role_id, start_time);


-- +goose Down
DROP INDEX IF EXISTS idx_shift_requests_covering_filter;
DROP INDEX IF EXISTS idx_shifts_covering_filter;
