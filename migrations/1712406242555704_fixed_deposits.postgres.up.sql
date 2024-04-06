CREATE TABLE IF NOT EXISTS fixed_deposits (
    id serial PRIMARY KEY,
    organization_unit_id INTEGER NOT NULL,
    judge_id INTEGER,
    subject TEXT,
    case_number TEXT,
    date_of_recipiet TIMESTAMP, --datum prijema akta
    date_of_case TIMESTAMP, --datum predemta
    date_of_finality TIMESTAMP, --datum pravosnaznosti
    date_of_enforceability TIMESTAMP, --datum izvrsnosti
    date_of_end TIMESTAMP, -- datum zakljucenja
    account_id INTEGER,
    file_id INTEGER,
    status TEXT,
    type TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
