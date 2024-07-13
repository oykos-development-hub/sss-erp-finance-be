CREATE TABLE IF NOT EXISTS logs (
    id serial PRIMARY KEY,
    operation VARCHAR(10),
    entity TEXT,
    old_state JSONB,
    new_state JSONB,
    user_id INTEGER,
    item_id INTEGER,
    changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE FUNCTION log_changes() RETURNS TRIGGER AS $$
DECLARE
    user_id INTEGER;
BEGIN
    BEGIN
        SELECT current_setting('myapp.user_id')::INTEGER INTO user_id;
    EXCEPTION
        WHEN others THEN
            user_id := 0;  
    END;

    IF TG_OP = 'INSERT' THEN
        INSERT INTO logs (operation, new_state, user_id, item_id, entity)
        VALUES ('INSERT', row_to_json(NEW)::jsonb, user_id, NEW.id, TG_TABLE_NAME);
        RETURN NEW;
    ELSIF TG_OP = 'UPDATE' THEN
        INSERT INTO logs (operation, old_state, new_state, user_id, item_id, entity)
        VALUES ('UPDATE', row_to_json(OLD)::jsonb, row_to_json(NEW)::jsonb, user_id, NEW.id, TG_TABLE_NAME);
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        INSERT INTO logs (operation, old_state, user_id, item_id, entity)
        VALUES ('DELETE', row_to_json(OLD)::jsonb, user_id, OLD.id, TG_TABLE_NAME);
        RETURN OLD;
    END IF;
END;
$$ LANGUAGE plpgsql;


CREATE TRIGGER budgets_insert
AFTER INSERT ON budgets
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER budgets_update
AFTER UPDATE ON budgets
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER budgets_delete
AFTER DELETE ON budgets
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER accounting_entries_insert
AFTER INSERT ON accounting_entries
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER accounting_entries_update
AFTER UPDATE ON accounting_entries
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER accounting_entries_delete
AFTER DELETE ON accounting_entries
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER activities_insert
AFTER INSERT ON activities
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER activities_update
AFTER UPDATE ON activities
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER activities_delete
AFTER DELETE ON activities
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER budget_requests_insert
AFTER INSERT ON budget_requests
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER budget_requests_update
AFTER UPDATE ON budget_requests
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER budget_requests_delete
AFTER DELETE ON budget_requests
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER current_budgets_insert
AFTER INSERT ON current_budgets
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER current_budgets_update
AFTER UPDATE ON current_budgets
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER current_budgets_delete
AFTER DELETE ON current_budgets
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER deposit_payment_orders_insert
AFTER INSERT ON deposit_payment_orders
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER deposit_payment_orders_update
AFTER UPDATE ON deposit_payment_orders
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER deposit_payment_orders_delete
AFTER DELETE ON deposit_payment_orders
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER deposit_payments_insert
AFTER INSERT ON deposit_payments
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER deposit_payments_update
AFTER UPDATE ON deposit_payments
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER deposit_payments_delete
AFTER DELETE ON deposit_payments
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER enforced_payments_insert
AFTER INSERT ON enforced_payments
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER enforced_payments_update
AFTER UPDATE ON enforced_payments
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER enforced_payments_delete
AFTER DELETE ON enforced_payments
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER external_reallocations_insert
AFTER INSERT ON external_reallocations
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER external_reallocations_update
AFTER UPDATE ON external_reallocations
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER external_reallocations_delete
AFTER DELETE ON external_reallocations
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER fee_payments_insert
AFTER INSERT ON fee_payments
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER fee_payments_update
AFTER UPDATE ON fee_payments
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER fee_payments_delete
AFTER DELETE ON fee_payments
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER fees_insert
AFTER INSERT ON fees
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER fees_update
AFTER UPDATE ON fees
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER fees_delete
AFTER DELETE ON fees
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER filled_financial_budgets_insert
AFTER INSERT ON filled_financial_budgets
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER filled_financial_budgets_update
AFTER UPDATE ON filled_financial_budgets
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER filled_financial_budgets_delete
AFTER DELETE ON filled_financial_budgets
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER financial_budgets_insert
AFTER INSERT ON financial_budgets
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER financial_budgets_update
AFTER UPDATE ON financial_budgets
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER financial_budgets_delete
AFTER DELETE ON financial_budgets
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER fine_payments_insert
AFTER INSERT ON fine_payments
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER fine_payments_update
AFTER UPDATE ON fine_payments
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER fine_payments_delete
AFTER DELETE ON fine_payments
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER fines_insert
AFTER INSERT ON fines
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER fines_update
AFTER UPDATE ON fines
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER fines_delete
AFTER DELETE ON fines
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER fixed_deposit_dispatches_insert
AFTER INSERT ON fixed_deposit_dispatches
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER fixed_deposit_dispatches_update
AFTER UPDATE ON fixed_deposit_dispatches
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER fixed_deposit_dispatches_delete
AFTER DELETE ON fixed_deposit_dispatches
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER fixed_deposit_items_insert
AFTER INSERT ON fixed_deposit_items
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER fixed_deposit_items_update
AFTER UPDATE ON fixed_deposit_items
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER fixed_deposit_items_delete
AFTER DELETE ON fixed_deposit_items
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER fixed_deposit_will_dispatches_insert
AFTER INSERT ON fixed_deposit_will_dispatches
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER fixed_deposit_will_dispatches_update
AFTER UPDATE ON fixed_deposit_will_dispatches
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER fixed_deposit_will_dispatches_delete
AFTER DELETE ON fixed_deposit_will_dispatches
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER fixed_deposit_wills_insert
AFTER INSERT ON fixed_deposit_wills
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER fixed_deposit_wills_update
AFTER UPDATE ON fixed_deposit_wills
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER fixed_deposit_wills_delete
AFTER DELETE ON fixed_deposit_wills
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER fixed_deposits_insert
AFTER INSERT ON fixed_deposits
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER fixed_deposits_update
AFTER UPDATE ON fixed_deposits
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER fixed_deposits_delete
AFTER DELETE ON fixed_deposits
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER flat_rate_payments_insert
AFTER INSERT ON flat_rate_payments
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER flat_rate_payments_update
AFTER UPDATE ON flat_rate_payments
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER flat_rate_payments_delete
AFTER DELETE ON flat_rate_payments
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER flat_rates_insert
AFTER INSERT ON flat_rates
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER flat_rates_update
AFTER UPDATE ON flat_rates
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER flat_rates_delete
AFTER DELETE ON flat_rates
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER goal_indicators_insert
AFTER INSERT ON goal_indicators
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER goal_indicators_update
AFTER UPDATE ON goal_indicators
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER goal_indicators_delete
AFTER DELETE ON goal_indicators
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER internal_reallocations_insert
AFTER INSERT ON internal_reallocations
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER internal_reallocations_update
AFTER UPDATE ON internal_reallocations
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER internal_reallocations_delete
AFTER DELETE ON internal_reallocations
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER invoices_insert
AFTER INSERT ON invoices
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER invoices_update
AFTER UPDATE ON invoices
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER invoices_delete
AFTER DELETE ON invoices
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER models_of_accounting_insert
AFTER INSERT ON models_of_accounting
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER models_of_accounting_update
AFTER UPDATE ON models_of_accounting
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER models_of_accounting_delete
AFTER DELETE ON models_of_accounting
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER non_financial_budget_goals_insert
AFTER INSERT ON non_financial_budget_goals
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER non_financial_budget_goals_update
AFTER UPDATE ON non_financial_budget_goals
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER non_financial_budget_goals_delete
AFTER DELETE ON non_financial_budget_goals
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER non_financial_budgets_insert
AFTER INSERT ON non_financial_budgets
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER non_financial_budgets_update
AFTER UPDATE ON non_financial_budgets
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER non_financial_budgets_delete
AFTER DELETE ON non_financial_budgets
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER payment_orders_insert
AFTER INSERT ON payment_orders
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER payment_orders_update
AFTER UPDATE ON payment_orders
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER payment_orders_delete
AFTER DELETE ON payment_orders
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER procedure_cost_payments_insert
AFTER INSERT ON procedure_cost_payments
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER procedure_cost_payments_update
AFTER UPDATE ON procedure_cost_payments
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER procedure_cost_payments_delete
AFTER DELETE ON procedure_cost_payments
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER procedure_costs_insert
AFTER INSERT ON procedure_costs
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER procedure_costs_update
AFTER UPDATE ON procedure_costs
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER procedure_costs_delete
AFTER DELETE ON procedure_costs
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER programs_insert
AFTER INSERT ON programs
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER programs_update
AFTER UPDATE ON programs
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER programs_delete
AFTER DELETE ON programs
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER property_benefits_confiscations_insert
AFTER INSERT ON property_benefits_confiscations
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER property_benefits_confiscations_update
AFTER UPDATE ON property_benefits_confiscations
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER property_benefits_confiscations_delete
AFTER DELETE ON property_benefits_confiscations
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER property_benefits_confiscations_payments_insert
AFTER INSERT ON property_benefits_confiscations_payments
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER property_benefits_confiscations_payments_update
AFTER UPDATE ON property_benefits_confiscations_payments
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER property_benefits_confiscations_payments_delete
AFTER DELETE ON property_benefits_confiscations_payments
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER salaries_insert
AFTER INSERT ON salaries
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER salaries_update
AFTER UPDATE ON salaries
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER salaries_delete
AFTER DELETE ON salaries
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER spending_dynamic_entries_insert
AFTER INSERT ON spending_dynamic_entries
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER spending_dynamic_entries_update
AFTER UPDATE ON spending_dynamic_entries
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER spending_dynamic_entries_delete
AFTER DELETE ON spending_dynamic_entries
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER spending_releases_insert
AFTER INSERT ON spending_releases
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER spending_releases_update
AFTER UPDATE ON spending_releases
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER spending_releases_delete
AFTER DELETE ON spending_releases
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER tax_authority_codebooks_insert
AFTER INSERT ON tax_authority_codebooks
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER tax_authority_codebooks_update
AFTER UPDATE ON tax_authority_codebooks
FOR EACH ROW
EXECUTE FUNCTION log_changes();

CREATE TRIGGER tax_authority_codebooks_delete
AFTER DELETE ON tax_authority_codebooks
FOR EACH ROW
EXECUTE FUNCTION log_changes();
