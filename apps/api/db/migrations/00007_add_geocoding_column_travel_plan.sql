-- +goose Up

ALTER TABLE travel_plans
    ADD COLUMN country VARCHAR(100) NULL AFTER place,
    ADD COLUMN state VARCHAR(100) NULL AFTER country,
    ADD COLUMN city VARCHAR(100) NULL AFTER state,
    ADD COLUMN latitude DECIMAL(9,6) NULL AFTER city,
    ADD COLUMN longitude DECIMAL(10,6) NULL AFTER latitude;

-- +goose Down

ALTER TABLE travel_plans
    DROP COLUMN longitude,
    DROP COLUMN latitude,
    DROP COLUMN city,
    DROP COLUMN state,
    DROP COLUMN country;
