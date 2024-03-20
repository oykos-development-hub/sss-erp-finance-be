CREATE TABLE invoices (
    id serial PRIMARY KEY,
    invoice_number TEXT,
    type TEXT NOT NULL,
    type_of_subject INTEGER,
    type_of_contract INTEGER,
    status TEXT,
    source_of_funding INTEGER,
    gross_price FLOAT,
    vat_price FLOAT,
    supplier TEXT,
    supplier_id INTEGER,
    order_id INTEGER,
    organization_unit_id INTEGER,
    activity_id INTEGER,
    date_of_invoice TIMESTAMP,             --datum fakture/resenja, 
    receipt_date TIMESTAMP,                --datum prijema,
    date_of_payment TIMESTAMP,             --datum placanja,
    sss_invoice_receipt_date TIMESTAMP,    --datum sss, 
    date_of_start TIMESTAMP,               --datum pocetka ugovora
    file_id INTEGER,
    bank_account TEXT,
    description TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
