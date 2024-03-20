CREATE TABLE IF NOT EXISTS flat_rates (
    id serial PRIMARY KEY,
    flat_rate_type INTEGER NOT NULL,
    decision_number  TEXT NOT NULL,
    decision_date TIMESTAMP,
    subject TEXT,
    jmbg TEXT NOT NULL,
    residence TEXT,
    amount FLOAT,
    payment_reference_number TEXT,
    debit_reference_number TEXT,
    account_id INT, 
    execution_date TIMESTAMP,
    payment_deadline_date TIMESTAMP,
    description TEXT,
    status INTEGER NOT NULL DEFAULT 1,
    court_costs NUMERIC(10, 2) NOT NULL,
    court_account_id INTEGER,
    file INTEGER[],
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS flat_rate_payments (
    id serial PRIMARY KEY,
    flat_rate_id INT NULL REFERENCES flat_rate (id) ON DELETE CASCADE,
    payment_method INTEGER NOT NULL,
    amount NUMERIC(10, 2) NOT NULL,
    payment_date TIMESTAMP,
    payment_due_date TIMESTAMP,
    receipt_number TEXT,  
    payment_reference_number TEXT,
    debit_reference_number TEXT,
    status INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);