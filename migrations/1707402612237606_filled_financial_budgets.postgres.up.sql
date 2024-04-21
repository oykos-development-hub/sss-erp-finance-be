CREATE TABLE IF NOT EXISTS filled_financial_budgets (
    id serial PRIMARY KEY,
    budget_request_id INTEGER NOT NULL REFERENCES budget_requests(id) ON DELETE CASCADE,
    account_id INTEGER NOT NULL,
    current_year DECIMAL(10, 2) NOT NULL,
    next_year DECIMAL(10, 2) NOT NULL,
    year_after_next DECIMAL(10, 2) NOT NULL,
    actual DECIMAL(10, 2),
    description TEXT NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
