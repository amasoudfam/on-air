CREATE TABLE users (
  id int PRIMARY KEY,
  first_name text,
  last_name text,
  email text,
  password text,
  phone_number text
);

CREATE TABLE countries (
  id int PRIMARY KEY,
  name text
);

CREATE TABLE cities (
  id int PRIMARY KEY,
  name text,
  country_id int
);
ALTER TABLE cities ADD FOREIGN KEY (country_id) REFERENCES countries (id);

CREATE TABLE flights (
  id int PRIMARY KEY,
  number text,
  from_city int,
  to_city int,
  airplane text,
  airline text,
  started_at date,
  ended_at date
);
ALTER TABLE flights ADD FOREIGN KEY (from_city) REFERENCES cities (id);
ALTER TABLE flights ADD FOREIGN KEY (to_city) REFERENCES cities (id);

CREATE TABLE tickets (
  id int PRIMARY KEY,
  user_id int,
  unit_price int,
  count int,
  flight_id int,
  created_at date,
  status text
);
ALTER TABLE tickets ADD FOREIGN KEY (user_id) REFERENCES users (id);
ALTER TABLE tickets ADD FOREIGN KEY (flight_id) REFERENCES flights (id);

CREATE TABLE passengers (
  id int PRIMARY KEY,
  ticket_id int,
  national_code text,
  first_name text,
  last_name text,
  gender text
);
ALTER TABLE passengers ADD FOREIGN KEY (ticket_id) REFERENCES tickets (id);

CREATE TABLE payments (
  id int PRIMARY KEY,
  amount int,
  status text,
  payed_at date,
  ticket_id int
);
ALTER TABLE payments ADD FOREIGN KEY (ticket_id) REFERENCES tickets (id);


