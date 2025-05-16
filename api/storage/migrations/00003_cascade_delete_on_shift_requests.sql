-- +goose Up
ALTER TABLE shift_requests DROP CONSTRAINT shift_requests_shift_id_fkey;

ALTER TABLE shift_requests
ADD CONSTRAINT shift_requests_shift_id_fkey
FOREIGN KEY (shift_id) REFERENCES shifts(id) ON DELETE CASCADE;

-- +goose Down
ALTER TABLE shift_requests DROP CONSTRAINT shift_requests_shift_id_fkey;

ALTER TABLE shift_requests
ADD CONSTRAINT shift_requests_shift_id_fkey
FOREIGN KEY (shift_id) REFERENCES shifts(id);
