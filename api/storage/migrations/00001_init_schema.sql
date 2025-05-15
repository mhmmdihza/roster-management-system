-- +goose Up
CREATE TABLE roles (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL
);

INSERT INTO roles (id,name) VALUES 
(1,'Cashier'),
(2,'Cook'),
(3,'Waiter');

CREATE TABLE employees (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL CHECK (char_length(name) > 0),
    role_id INTEGER NOT NULL REFERENCES roles(id),
    status TEXT NOT NULL CHECK (status IN ('ACTIVE', 'INACTIVE')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE shifts (
    id SERIAL PRIMARY KEY,
    role_id INTEGER NOT NULL REFERENCES roles(id),
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE shift_requests (
    id SERIAL PRIMARY KEY,
    employee_id INTEGER NOT NULL REFERENCES employees(id),
    shift_id INTEGER NOT NULL REFERENCES shifts(id),
    status TEXT NOT NULL CHECK (status IN ('PENDING', 'APPROVED', 'REJECTED')),
    requested_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    reviewed_at TIMESTAMP,
    reviewed_by INTEGER REFERENCES employees(id)
);

-- +goose Down
DROP TABLE IF EXISTS shift_requests;
DROP TABLE IF EXISTS shifts;
DROP TABLE IF EXISTS employees;
DROP TABLE IF EXISTS roles;
