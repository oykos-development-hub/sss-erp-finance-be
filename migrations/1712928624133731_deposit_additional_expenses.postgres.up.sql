CREATE TABLE IF NOT EXISTS deposit_additional_expenses (
    id serial PRIMARY KEY,
    price float not null,
    account_id integer not null,
    title text,
    status text not null,
    subject_id integer not null,
    organization_unit_id integer not null,
    bank_account text not null,
    payment_order_id integer not null references deposit_payment_orders(id) on delete cascade,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
