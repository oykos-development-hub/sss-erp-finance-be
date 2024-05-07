CREATE TABLE IF NOT EXISTS accounting_entries (
    id serial PRIMARY KEY,
    title VARCHAR ( 255 ) NOT NULL,
    organization_unit_id INTEGER NOT NULL,
    id_of_entry INTEGER NOT NULL,
    date_of_booking TIMESTAMP,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE OR REPLACE FUNCTION increment_entry_id()
RETURNS TRIGGER AS $$
BEGIN
    -- Pronalazi poslednji id_of_entry za organization_unit_id
    SELECT COALESCE(MAX(id_of_entry), 0) + 1 INTO NEW.id_of_entry
    FROM accounting_entries
    WHERE organization_unit_id = NEW.organization_unit_id;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER entry_id_before_insert
BEFORE INSERT ON accounting_entries
FOR EACH ROW
EXECUTE FUNCTION increment_entry_id();
