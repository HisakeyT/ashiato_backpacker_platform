-- +goose Up

ALTER TABLE travel_plans RENAME COLUMN date TO start_date;

-- +goose Down

ALTER TABLE travel_plans RENAME COLUMN start_date TO date;
