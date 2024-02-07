CREATE TABLE IF NOT EXISTS articles (
    id serial PRIMARY KEY,
    title VARCHAR ( 255 ) NOT NULL,
    net_price NUMERIC(10, 2) NOT NULL,
    vat_price NUMERIC(10, 2) NOT NULL,
    description TEXT,
    invoice_id INT NOT NULL,
    account_id INT NOT NULL,
    cost_account_id INT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
