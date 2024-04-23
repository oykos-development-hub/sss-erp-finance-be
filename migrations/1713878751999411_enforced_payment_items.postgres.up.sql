
	
CREATE TABLE IF NOT EXISTS enforced_payment_items (
    id serial PRIMARY KEY,
    title VARCHAR ( 255 ) NOT NULL,
	payment_order_id INTEGER NOT NULL REFERENCES enforced_payments(id) ON DELETE CASCADE,
	amount FLOAT,
	invoice_id INTEGER REFERENCES invoices(id),
	account_id INTEGER,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
