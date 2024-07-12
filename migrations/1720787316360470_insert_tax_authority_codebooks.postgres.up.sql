DELETE FROM tax_authority_codebooks;

INSERT INTO tax_authority_codebooks (
    title, code, active, tax_percentage, tax_supplier_id, release_percentage, 
    pio_percentage, pio_supplier_id, pio_percentage_employer_percentage, pio_employer_supplier_id, 
    pio_percentage_employee_percentage, pio_employee_supplier_id, unemployment_percentage, 
    unemployment_supplier_id, unemployment_employer_percentage, unemployment_employer_supplier_id, 
    unemployment_employee_percentage, unemployment_employee_supplier_id, labor_fund, 
    labor_fund_supplier_id, previous_income_percentage_less_than_700, 
    previous_income_percentage_less_than_1000, previous_income_percentage_more_than_1000, 
    coefficient, created_at, updated_at, coefficient_less_700, coefficient_less_1000, 
    coefficient_more_1000, include_subtax, amount_less_700, amount_less_1000, amount_more_1000, 
    release_amount
) VALUES
('036 - NAKNADA STECAJNOG UPRAVNIKA', '036', 't', 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 9, 15, 0, NOW(), NOW(), 0, 0, 0, 'f', 0, 0, 0, 0),
('047 - NAKNADA SAMO POREZ', '047', 't', 15, 1, 30, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, NOW(), NOW(), 0, 0, 0, 'f', 0, 0, 0, 0),
('Depozit - parnica', '047 - Depozit', 't', 15, 1, 30, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, NOW(), NOW(), 0, 0, 0, 'f', 0, 0, 0, 0),
('048 - UGOVORENA NAKNADA', '048', 't', 15, 1, 30, 20.5, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, NOW(), NOW(), 0, 0, 0, 'f', 0, 0, 0, 0),
('NAKNADE - SAMO POREZ', '', 't', 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 9, 15, 0, NOW(), NOW(), 0, 0, 0, 't', 0, 0, 0, 0),
('UGOVOR O PRIVREMENO POVREMENOM VRSENJU POSLOVA', '', 't', 0, 1, 0, 0, 1, 15, 0, 5.5, 0, 0, 1, 0.5, 0, 0.5, 0, 0.2, 1, 0, 9, 15, 0, NOW(), NOW(), 0, 0, 0, 't', 0, 0, 0, 0),
('NAKNADA SVI DOPRINOSI', '', 't', 0, 1, 0, 0, 1, 15, 0, 5.5, 0, 0, 1, 0.5, 0, 0.5, 0, 0.2, 1, 0, 9, 15, 0, NOW(), NOW(), 0, 0, 0, 't', 0, 0, 0, 0);
