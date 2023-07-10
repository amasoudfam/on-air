ALTER TABLE flights rename column from_city TO from_city_id;
ALTER TABLE flights rename column to_city TO to_city_id;
ALTER TABLE flights ADD FOREIGN KEY (from_city_id) REFERENCES cities (id);
ALTER TABLE flights ADD FOREIGN KEY (from_city_id) REFERENCES cities (id);

