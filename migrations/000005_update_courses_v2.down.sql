-- Reverse of 000005: remove v2 course columns and clear v2 data
ALTER TABLE courses DROP COLUMN IF EXISTS provider;
ALTER TABLE courses DROP COLUMN IF EXISTS page_type;
ALTER TABLE courses DROP COLUMN IF EXISTS level;
ALTER TABLE courses DROP COLUMN IF EXISTS direct_link_note;
ALTER TABLE courses DROP COLUMN IF EXISTS relevance_reason;
ALTER TABLE courses DROP COLUMN IF EXISTS course_identifier;
ALTER TABLE courses DROP COLUMN IF EXISTS source_code;
ALTER TABLE courses DROP COLUMN IF EXISTS source_name;
ALTER TABLE courses DROP COLUMN IF EXISTS cost_type;
ALTER TABLE courses DROP COLUMN IF EXISTS notes;
ALTER TABLE stages DROP COLUMN IF EXISTS stage_identifier;
ALTER TABLE learning_paths DROP COLUMN IF EXISTS career_identifier;

-- Clear v2 course data (v1 data is not recoverable from this migration)
DELETE FROM courses;
