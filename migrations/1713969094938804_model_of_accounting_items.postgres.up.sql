CREATE TABLE IF NOT EXISTS model_of_accounting_items (
    id serial PRIMARY KEY,
    title VARCHAR ( 255 ) NOT NULL,
    model_id INTEGER NOT NULL REFERENCES models_of_accounting(id) ON DELETE CASCADE,
    debit_account_id INTEGER,
    credit_account_id INTEGER,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
