-- +goose Up

ALTER TABLE travel_plans
    ADD COLUMN end_date DATE NULL AFTER start_date;

-- +goose Down

ALTER TABLE travel_plans
    DROP COLUMN end_date;
