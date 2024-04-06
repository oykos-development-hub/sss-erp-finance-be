CREATE TABLE IF NOT EXISTS fixed_deposit_judges (
    id serial PRIMARY KEY,
    judge_id INTEGER NOT NULL,
    deposit_id INTEGER REFERENCES fixed_deposits(id) ON DELETE CASCADE,
    will_id INTEGER, -- REFERENCES fixed_deposit_wills(id) ON DELETE CASCADE,
    date_of_start TIMESTAMP,
    date_of_end TIMESTAMP,
    file_id INTEGER,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
