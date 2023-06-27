ALTER TABLE users 
DROP CONSTRAINT uk_users_email;

ALTER TABLE users 
DROP CONSTRAINT uk_users_phone_number;

ALTER TABLE passengers 
DROP CONSTRAINT uk_passengers_user_id_national_code;

ALTER TABLE ticket_passengers 
DROP CONSTRAINT uk_ticket_passengers_ticket_id_passenger_ide;
