CREATE TABLE IF NOT EXISTS fixed_deposit_wills (
    id serial PRIMARY KEY,
    organization_unit_id INTEGER NOT NULL,
    subject TEXT,
    father_name TEXT,
    date_of_birth TIMESTAMP,
    jmbg TEXT,
    case_number_si TEXT,
    case_number_rs TEXT,
    date_of_receipt_si TIMESTAMP,
    date_of_receipt_rs TIMESTAMP,
    date_of_end TIMESTAMP,
    status TEXT,
    file_id INTEGER,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
