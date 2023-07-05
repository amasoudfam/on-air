ALTER TABLE passengers
ADD CONSTRAINT uk_passengers_user_id_national_code UNIQUE (user_id, national_code);

