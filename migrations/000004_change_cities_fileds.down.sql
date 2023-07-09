ALTER TABLE flights
RENAME COLUMN from_city_id TO from_city;

ALTER TABLE flights
RENAME COLUMN to_city_id TO to_city;

ALTER TABLE flights
RENAME COLUMN finished_at TO ended_at;