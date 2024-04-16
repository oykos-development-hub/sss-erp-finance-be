CREATE TABLE IF NOT EXISTS salary_additional_expenses (
    id serial PRIMARY KEY,
    type TEXT NOT NULL,
    salary_id INTEGER NOT NULL REFERENCES salaries(id) ON DELETE CASCADE,
    account_id INTEGER,
    amount FLOAT NOT NULL,
    subject_id INTEGER,
    bank_account TEXT,
    status TEXT NOT NULL,
    title TEXT,
    debtor_id INTEGER,
    organization_unit_id INTEGER NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
