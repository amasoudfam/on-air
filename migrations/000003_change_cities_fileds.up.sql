ALTER TABLE flights
RENAME COLUMN from_city TO from_city_id;

ALTER TABLE flights
RENAME COLUMN to_city TO to_city_id;

ALTER TABLE flights
RENAME COLUMN ended_at TO finished_at;