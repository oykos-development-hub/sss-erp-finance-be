CREATE TABLE IF NOT EXISTS salaries (
    id serial PRIMARY KEY,
    activity_id INTEGER NOT NULL,
    month TEXT NOT NULL,
    date_of_calculation TIMESTAMP,
    description TEXT,
    status TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
