CREATE TABLE IF NOT EXISTS deposit_payment_orders (
    id serial PRIMARY KEY,
    organization_unit_id INTEGER NOT NULL,
    case_number TEXT,
    supplier_id INTEGER,
    net_amount FLOAT,
    bank_account TEXT,
    date_of_payment TIMESTAMP,
    date_of_statement TIMESTAMP,
    id_of_statement TIMESTAMP,
    file_id INTEGER,
    tax_authority_codebook_id INTEGER,
    municipality_id INTEGER,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
