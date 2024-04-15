CREATE TABLE IF NOT EXISTS deposit_payments (
    id serial PRIMARY KEY,
    payer TEXT,
    organization_unit_id INTEGER NOT NULL,
    case_number TEXT,
    party_name TEXT,
    number_of_bank_statement TEXT,
    date_of_bank_statement TEXT,
    account_id INTEGER,
    amount FLOAT,
    main_bank_account BOOLEAN,
    date_of_transfer_main_account TIMESTAMP,
    file_id INTEGER,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
