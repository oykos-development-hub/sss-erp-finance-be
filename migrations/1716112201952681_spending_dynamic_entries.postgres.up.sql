CREATE TABLE IF NOT EXISTS spending_dynamic_entries (
    id serial PRIMARY KEY,
    spending_dynamic_id INTEGER NOT NULL REFERENCES spending_dynamics(id) ON DELETE CASCADE,
    username TEXT NOT NULL,
    january DECIMAL(10, 2) NOT NULL NOT NULL,
    february DECIMAL(10, 2) NOT NULL,
    march DECIMAL(10, 2) NOT NULL,
    april DECIMAL(10, 2) NOT NULL,
    may DECIMAL(10, 2) NOT NULL,
    june DECIMAL(10, 2) NOT NULL,
    july DECIMAL(10, 2) NOT NULL,
    august DECIMAL(10, 2) NOT NULL,
    september DECIMAL(10, 2) NOT NULL,
    october DECIMAL(10, 2) NOT NULL,
    november DECIMAL(10, 2) NOT NULL,
    december DECIMAL(10, 2) NOT NULL,
    version INTEGER NOT NULL,
    created_at TIMESTAMP
);
