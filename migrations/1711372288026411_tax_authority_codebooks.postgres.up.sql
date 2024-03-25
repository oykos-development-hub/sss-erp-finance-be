CREATE TABLE IF NOT EXISTS tax_authority_codebooks (
    id serial PRIMARY KEY,
    title VARCHAR ( 255 ) NOT NULL,
    code TEXT NOT NULL,
    pergentage FLOAT NOT NULL,
    previous_income_percentage_less_than_700 FLOAT,
    previous_income_percentage_less_than_1000 FLOAT,
    previous_income_percentage_more_than_1000 FLOAT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
