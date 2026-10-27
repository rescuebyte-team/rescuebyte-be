ALTER TABLE users
    ADD COLUMN IF NOT EXISTS phone_number VARCHAR(20) NULL,
    ADD COLUMN IF NOT EXISTS gender SMALLINT NULL,
    ADD COLUMN IF NOT EXISTS date_of_birth DATE NULL;

CREATE INDEX IF NOT EXISTS idx_users_phone_number ON users(phone_number);
