ALTER TABLE flights rename column from_city_id TO from_city;
ALTER TABLE flights rename column to_city_id TO to_city;
ALTER TABLE flights ADD FOREIGN KEY (from_city) REFERENCES cities (id);
ALTER TABLE flights ADD FOREIGN KEY (to_city) REFERENCES cities (id);

