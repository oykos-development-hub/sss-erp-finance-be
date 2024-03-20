CREATE TABLE IF NOT EXISTS additional_expenses (
    id serial PRIMARY KEY,
    price float not null,
    account_id integer not null,
    title text,
    status integer not null,
    subject_id integer not null,
    bank_account integer not null,
    invoice_id integer not null references invoices(id) on delete cascade,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
