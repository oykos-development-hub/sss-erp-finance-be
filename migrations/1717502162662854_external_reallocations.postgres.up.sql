CREATE TABLE IF NOT EXISTS external_reallocations (
    id serial PRIMARY KEY,
    title TEXT,
    budget_id INTEGER NOT NULL REFERENCES budgets(id) ON DELETE CASCADE,
    source_organization_unit_id INTEGER NOT NULL,
    destination_organization_unit_id INTEGER NOT NULL,
    date_of_request TIMESTAMP NOT NULL,
    date_of_action_dest_org_unit TIMESTAMP,
    date_of_action_sss TIMESTAMP,
    requested_by INTEGER NOT NULL,
    accepted_by INTEGER,
    file_id INTEGER,
    destination_org_unit_file_id INTEGER,
    sss_file_id INTEGER,
    status TEXT NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
