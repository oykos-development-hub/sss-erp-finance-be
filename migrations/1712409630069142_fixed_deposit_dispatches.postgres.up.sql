CREATE TABLE IF NOT EXISTS fixed_deposit_dispatches (
    id serial PRIMARY KEY,
    deposit_id INTEGER REFERENCES fixed_deposits(id) ON DELETE CASCADE,
    category_id INTEGER,
    type_id INTEGER,
    unit TEXT,
    amount FLOAT,
    serial_number TEXT,
    date_of_action TIMESTAMP,
    subject TEXT,
    action TEXT,
    case_number TEXT,
    judge_id INTEGER,
    file_id INTEGER,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
