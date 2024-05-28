CREATE TABLE IF NOT EXISTS payment_order_items (
    id serial PRIMARY KEY,
    payment_order_id INTEGER NOT NULL REFERENCES payment_orders(id) ON DELETE CASCADE,
    invoice_id INTEGER REFERENCES invoices(id),
    additional_expense_id INTEGER REFERENCES additional_expenses(id),
    salary_additional_expense_id INTEGER REFERENCES salary_additional_expenses(id),
    account_id INTEGER,
    source_account TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
