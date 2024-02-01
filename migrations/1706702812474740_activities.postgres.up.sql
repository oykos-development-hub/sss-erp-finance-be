CREATE TABLE IF NOT EXISTS activities (
    id serial PRIMARY KEY,
    sub_program_id INTEGER NOT NULL references programs(id),
    organization_unit_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    code text NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
