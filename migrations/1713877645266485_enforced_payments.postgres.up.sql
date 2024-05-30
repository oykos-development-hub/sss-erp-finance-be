CREATE TABLE IF NOT EXISTS enforced_payments (
    id serial PRIMARY KEY,
    organization_unit_id INTEGER NOT NULL,
    supplier_id INTEGER NOT NULL,
    sap_id TEXT,
    date_of_sap TIMESTAMP,
    date_of_payment TIMESTAMP,
    date_of_order TIMESTAMP,
    description TEXT,
    file_id INTEGER,
    bank_account TEXT,
    amount FLOAT,
    amount_for_lawyer FLOAT,
    amount_for_agent FLOAT,
    amount_for_bank FLOAT,
    agent_id INTEGER,
    execution_number TEXT,
    registred BOOLEAN default false,
    registred_return BOOLEAN default false,
    return_date TIMESTAMP,
    return_file_id INTEGER,
    return_amount FLOAT,
    status TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE SEQUENCE statement_id_seq;

ALTER TABLE payment_orders ADD COLUMN id_of_statement BIGINT;
ALTER TABLE enforced_payments ADD COLUMN id_of_statement BIGINT;

CREATE OR REPLACE FUNCTION increment_id_of_statement()
RETURNS TRIGGER AS $$
BEGIN
    NEW.id_of_statement := nextval('statement_id_seq');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER payment_orders_id_trigger
BEFORE INSERT ON payment_orders
FOR EACH ROW
EXECUTE FUNCTION increment_id_of_statement();

CREATE TRIGGER enforced_payments_id_trigger
BEFORE INSERT ON enforced_payments
FOR EACH ROW
EXECUTE FUNCTION increment_id_of_statement();
