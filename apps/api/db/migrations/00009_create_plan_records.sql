-- +goose Up

CREATE TABLE sightseeings (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    travel_plan_id BIGINT UNSIGNED NOT NULL,
    name VARCHAR(255) NOT NULL,
    location_text TEXT NULL,
    map_url TEXT NULL,
    cost INT UNSIGNED NULL,
    sort_order INT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    INDEX idx_sightseeings_plan (travel_plan_id),
    FOREIGN KEY (travel_plan_id) REFERENCES travel_plans(id) ON DELETE CASCADE
);

CREATE TABLE stays (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    travel_plan_id BIGINT UNSIGNED NOT NULL,
    name VARCHAR(255) NOT NULL,
    location_text TEXT NULL,
    map_url TEXT NULL,
    cost INT UNSIGNED NULL,
    sort_order INT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    INDEX idx_stays_plan (travel_plan_id),
    FOREIGN KEY (travel_plan_id) REFERENCES travel_plans(id) ON DELETE CASCADE
);

CREATE TABLE meals (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    travel_plan_id BIGINT UNSIGNED NOT NULL,
    name VARCHAR(255) NOT NULL,
    meal_type VARCHAR(16) NULL,
    location_text TEXT NULL,
    map_url TEXT NULL,
    cost INT UNSIGNED NULL,
    sort_order INT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    INDEX idx_meals_plan (travel_plan_id),
    CONSTRAINT chk_meal_type CHECK (meal_type IS NULL OR meal_type IN ('breakfast', 'lunch', 'dinner', 'snack', 'other')),
    FOREIGN KEY (travel_plan_id) REFERENCES travel_plans(id) ON DELETE CASCADE
);

CREATE TABLE others (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    travel_plan_id BIGINT UNSIGNED NOT NULL,
    name VARCHAR(255) NOT NULL,
    location_text TEXT NULL,
    map_url TEXT NULL,
    cost INT UNSIGNED NULL,
    sort_order INT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    INDEX idx_others_plan (travel_plan_id),
    FOREIGN KEY (travel_plan_id) REFERENCES travel_plans(id) ON DELETE CASCADE
);

-- +goose Down

DROP TABLE others;
DROP TABLE meals;
DROP TABLE stays;
DROP TABLE sightseeings;
