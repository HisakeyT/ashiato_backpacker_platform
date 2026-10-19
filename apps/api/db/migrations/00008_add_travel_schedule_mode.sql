-- +goose Up

ALTER TABLE travels
    ADD COLUMN schedule_mode VARCHAR(16) NOT NULL DEFAULT 'planning' AFTER end_date,
    ADD CONSTRAINT chk_travels_schedule_mode CHECK (schedule_mode IN ('planning', 'recorded'));

UPDATE travels SET schedule_mode = 'recorded';

ALTER TABLE travel_plans MODIFY COLUMN start_date DATE NULL;

UPDATE travel_plans AS p
JOIN (
    SELECT id, ROW_NUMBER() OVER (PARTITION BY travel_id ORDER BY start_date, id) AS position
    FROM travel_plans
) AS ordered ON ordered.id = p.id
SET p.sort_order = ordered.position;

-- +goose Down

UPDATE travel_plans AS p
JOIN travels AS t ON t.id = p.travel_id
SET p.start_date = t.start_date
WHERE p.start_date IS NULL;

ALTER TABLE travel_plans MODIFY COLUMN start_date DATE NOT NULL;
ALTER TABLE travels DROP CHECK chk_travels_schedule_mode;
ALTER TABLE travels DROP COLUMN schedule_mode;
