BEGIN;

-- Update the decision_number column to TEXT type
ALTER TABLE fines
ALTER COLUMN decision_number TYPE TEXT USING decision_number::TEXT;

COMMIT;
