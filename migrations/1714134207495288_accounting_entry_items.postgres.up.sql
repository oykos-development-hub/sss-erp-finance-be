CREATE TABLE IF NOT EXISTS accounting_entry_items (
    id serial PRIMARY KEY,
    title VARCHAR ( 255 ) NOT NULL,
    entry_id INTEGER NOT NULL REFERENCES accounting_entries(id) ON DELETE CASCADE,
    account_id INTEGER,
    credit_amount FLOAT,
    debit_amount FLOAT,
    invoice_id INTEGER REFERENCES invoices(id),
    salary_id INTEGER REFERENCES salaries(id), 
    type TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
