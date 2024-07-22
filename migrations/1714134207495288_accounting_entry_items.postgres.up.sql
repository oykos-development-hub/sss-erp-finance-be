CREATE TABLE IF NOT EXISTS accounting_entry_items (
    id serial PRIMARY KEY,
    title VARCHAR ( 255 ) NOT NULL,
    entry_id INTEGER NOT NULL REFERENCES accounting_entries(id) ON DELETE CASCADE,
    account_id INTEGER,
    credit_amount FLOAT,
    debit_amount FLOAT,
    invoice_id INTEGER REFERENCES invoices(id),
    salary_id INTEGER REFERENCES salaries(id),
    payment_order_id INTEGER REFERENCES payment_orders(id),
    enforced_payment_id INTEGER REFERENCES enforced_payments(id),
    return_enforced_payment_id INTEGER REFERENCES enforced_payments(id), 
    supplier_id INTEGER,
    type TEXT,
    date TIMESTAMP,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
