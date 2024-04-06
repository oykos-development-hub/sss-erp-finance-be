CREATE TABLE IF NOT EXISTS fixed_deposit_will_dispatches (
    id serial PRIMARY KEY,
    will_id INTEGER REFERENCES fixed_deposit_wills(id) ON DELETE CASCADE,
    dispatch_type_id INTEGER NOT NULL,
    judge_id INTEGER NOT NULL,
    case_number TEXT,
    date_of_dispatch TIMESTAMP,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
