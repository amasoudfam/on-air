CREATE TABLE users (
  id int PRIMARY KEY,
  first_name varchar(50),
  last_name varchar(50),
  email varchar(50),
  password varchar(128),
  phone_number varchar(15),
  created_at timestamp with time zone,
  updated_at timestamp with time zone,
  deleted_at timestamp with time zone
);

CREATE TABLE countries (
  id int PRIMARY KEY,
  name varchar(50),
  created_at timestamp with time zone,
  updated_at timestamp with time zone,
  deleted_at timestamp with time zone
);

CREATE TABLE cities (
  id int PRIMARY KEY,
  name varchar(50),
  country_id int,
  created_at timestamp with time zone,
  updated_at timestamp with time zone,
  deleted_at timestamp with time zone
);
ALTER TABLE cities ADD FOREIGN KEY (country_id) REFERENCES countries (id);

CREATE TABLE flights (
  id int PRIMARY KEY,
  number varchar(20),
  from_city int,
  to_city int,
  airplane varchar(50),
  airline varchar(50),
  started_at date,
  ended_at date,
  created_at timestamp with time zone,
  updated_at timestamp with time zone,
  deleted_at timestamp with time zone
);
ALTER TABLE flights ADD FOREIGN KEY (from_city) REFERENCES cities (id);
ALTER TABLE flights ADD FOREIGN KEY (to_city) REFERENCES cities (id);

CREATE TABLE tickets (
  id int PRIMARY KEY,
  user_id int,
  unit_price int,
  count int,
  flight_id int,
  status varchar(10),
  created_at timestamp with time zone,
  updated_at timestamp with time zone,
  deleted_at timestamp with time zone
);

ALTER TABLE tickets ADD FOREIGN KEY (user_id) REFERENCES users (id);
ALTER TABLE tickets ADD FOREIGN KEY (flight_id) REFERENCES flights (id);

CREATE TABLE passengers (
  id int PRIMARY KEY,
  national_code varchar(10),
  first_name varchar(50),
  last_name varchar(50),
  gender varchar(5),
  user_id int,
  created_at timestamp with time zone,
  updated_at timestamp with time zone,
  deleted_at timestamp with time zone
);
ALTER TABLE passengers ADD FOREIGN KEY (user_id) REFERENCES users (id);

CREATE TABLE ticket_passengers (
  ticket_id int,
  passenger_id int
);
ALTER TABLE ticket_passengers ADD FOREIGN KEY (ticket_id) REFERENCES tickets (id);
ALTER TABLE ticket_passengers ADD FOREIGN KEY (passenger_id) REFERENCES passengers (id);

CREATE TABLE payments (
  id int PRIMARY KEY,
  amount int,
  status varchar(20),
  payed_at date,
  ticket_id int,
  created_at timestamp with time zone,
  updated_at timestamp with time zone,
  deleted_at timestamp with time zone
);
ALTER TABLE payments ADD FOREIGN KEY (ticket_id) REFERENCES tickets (id);
