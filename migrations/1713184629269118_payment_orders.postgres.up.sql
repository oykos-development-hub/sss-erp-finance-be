CREATE TABLE IF NOT EXISTS payment_orders (
    id serial PRIMARY KEY,
    organization_unit_id INTEGER NOT NULL,
    supplier_id INTEGER NOT NULL,
    id_of_statement TEXT,
    sap_id TEXT,
    date_of_sap TIMESTAMP,
    date_of_payment TIMESTAMP,
    description TEXT,
    file_id INTEGER,
    bank_account TEXT,
    amount FLOAT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
