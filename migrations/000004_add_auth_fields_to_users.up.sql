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

-- Copy legacy name -> display_name only when a legacy name column exists.
DO $$
BEGIN
  IF EXISTS (
    SELECT 1
    FROM information_schema.columns
    WHERE table_name = 'users' AND column_name = 'name'
  ) THEN
    EXECUTE 'UPDATE users SET display_name = name WHERE display_name IS NULL';
  END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
