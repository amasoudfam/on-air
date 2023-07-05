ALTER TABLE passengers
ADD CONSTRAINT uk_passengers_user_id_national_code UNIQUE (user_id, national_code);

ALTER TABLE users
ADD CONSTRAINT uk_users_phone_number UNIQUE (phone_number);

ALTER TABLE ticket_passengers
ADD CONSTRAINT uk_ticket_passengers_ticket_id_passenger_ide UNIQUE (ticket_id, passenger_id);