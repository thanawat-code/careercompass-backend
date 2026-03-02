-- Add columns required by auth_func branch
DO $$ BEGIN
  CREATE TYPE gender_enum AS ENUM ('male', 'female', 'other', 'prefer_not_to_say');
EXCEPTION
  WHEN duplicate_object THEN null;
END $$;

ALTER TABLE users
  ADD COLUMN IF NOT EXISTS display_name VARCHAR(255),
  ADD COLUMN IF NOT EXISTS password_hash VARCHAR(255),
  ADD COLUMN IF NOT EXISTS gender gender_enum;

-- Copy existing name -> display_name for backward compatibility
UPDATE users SET display_name = name WHERE display_name IS NULL;

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
