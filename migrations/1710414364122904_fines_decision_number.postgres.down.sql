BEGIN;

-- Update the decision_number column to INTEGER type
ALTER TABLE fines
ALTER COLUMN decision_number TYPE INTEGER USING decision_number::INTEGER;

COMMIT;
