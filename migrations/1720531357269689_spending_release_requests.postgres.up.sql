CREATE TABLE IF NOT EXISTS spending_release_requests (
    id serial PRIMARY KEY,
    year INTEGER NOT NULL,
    month INTEGER NOT NULL,
    organization_unit_file_id INTEGER,
    sss_file_id INTEGER,
    organization_unit_id INTEGER NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMP
);
