alter table spending_dynamic_entries add column total_amount DECIMAL(10, 2);

UPDATE spending_dynamic_entries sd
SET total_amount = cb.current_amount
FROM current_budgets cb
WHERE sd.current_budget_id = cb.id;

CREATE OR REPLACE FUNCTION update_total_amount()
RETURNS TRIGGER AS $$
BEGIN
    -- Ažuriramo polje total_amount na osnovu vrijednosti iz current_budgets
    NEW.total_amount := (
        SELECT cb.current_amount 
        FROM current_budgets cb 
        WHERE cb.id = NEW.current_budget_id
    );

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_total_amount
BEFORE INSERT ON spending_dynamic_entries
FOR EACH ROW
EXECUTE FUNCTION update_total_amount();
