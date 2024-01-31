CREATE TABLE IF NOT EXISTS non_financial_budgets (
    id serial PRIMARY KEY,
    budget_id INTEGER NOT NULL REFERENCES budgets(id),
    organization_unit_id INTEGER NOT NULL,
    impl_contact_fullname TEXT NOT NULL,
    impl_contact_working_place TEXT NOT NULL,
    impl_contact_phone TEXT NOT NULL,
    contact_fullname TEXT NOT NULL,
    contact_working_place TEXT NOT NULL,
    contact_phone TEXT NOT NULL,
    contact_email TEXT NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
