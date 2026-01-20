ALTER TABLE IF EXISTS markets
DROP COLUMN resolved_outcome_id,
DROP COLUMN resolved_at,
DROP COLUMN resolved_by;