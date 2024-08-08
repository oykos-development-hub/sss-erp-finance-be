delete from fine_payments;
delete from fines;
alter table fines add column organization_unit_id integer not null;

delete from fee_payments;
delete from fees;
alter table fees add column organization_unit_id integer not null;

delete from flat_rate_payments;
delete from flat_rates;
alter table flat_rates add column organization_unit_id integer not null;

delete from procedure_cost_payments;
delete from procedure_costs;
alter table procedure_costs add column organization_unit_id integer not null;

delete from property_benefits_confiscations_payments;
delete from property_benefits_confiscations;
alter table property_benefits_confiscations add column organization_unit_id integer not null;