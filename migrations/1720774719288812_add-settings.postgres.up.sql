INSERT INTO programs (parent_id, title, description, code, created_at, updated_at) VALUES
(NULL, 'Rad sudova', 'Rad sudova', 'RS', NOW(), NOW()),
(1, 'Nadleznost Osnovnih sudova', 'Nadleznost Osnovnih sudova', 'NOS', NOW(), NOW()),
(1, 'Rad Visih sudova', 'Rad Visih sudova', '12 003 007', NOW(), NOW()),
(1, 'Nadleznost Vrhovnog i apelacionog suda', 'Nadleznost Vrhovnog i apelacionog suda', '12 003 004', NOW(), NOW());

INSERT INTO activities (title, description, code, created_at, updated_at, organization_unit_id, sub_program_id) VALUES
('Vodjenje sudskih postupaka u ...', 'description', 'code', NOW(), NOW(), 1, 2),
('Vodjenje sudskih postupaka u ...', 'description', 'code', NOW(), NOW(), 2, 2),
('Vodjenje sudskih postupaka u ...', 'description', 'code', NOW(), NOW(), 3, 2),
('Vodjenje sudskih postupaka u ...', 'description', 'code', NOW(), NOW(), 4, 2),
('Vodjenje sudskih postupaka u ...', 'description', 'code', NOW(), NOW(), 5, 2),
('Vodjenje sudskih postupaka u ...', 'description', 'code', NOW(), NOW(), 6, 2),
('Vodjenje sudskih postupaka u ...', 'description', 'code', NOW(), NOW(), 7, 2),
('Vodjenje sudskih postupaka u ...', 'description', 'code', NOW(), NOW(), 8, 2),
('Vodjenje sudskih postupaka u ...', 'description', 'code', NOW(), NOW(), 9, 2),
('Vodjenje sudskih postupaka u ...', 'description', 'code', NOW(), NOW(), 10, 2),
('Vodjenje sudskih postupaka u ...', 'description', 'code', NOW(), NOW(), 11, 2),
('Vodjenje sudskih postupaka u ...', 'description', 'code', NOW(), NOW(), 12, 2),
('Vodjenje sudskih postupaka u ...', 'description', 'code', NOW(), NOW(), 13, 2),
('Vodjenje sudskih postupaka u ...', 'description', 'code', NOW(), NOW(), 14, 2),
('Vodjenje sudskih postupaka u ...', 'description', 'code', NOW(), NOW(), 15, 2),
('Vodjenje sudskih postupaka u ...', 'description', 'code', NOW(), NOW(), 16, 2),
('Vodjenje sudskih postupaka u ...', 'description', 'code', NOW(), NOW(), 17, 2),
('Vodjenje sudskih postupaka u ...', 'description', 'code', NOW(), NOW(), 18, 2),
('Vodjenje sudskih postupaka u ...', 'description', 'code', NOW(), NOW(), 19, 2),
('Vodjenje sudskih postupaka u ...', 'description', 'code', NOW(), NOW(), 20, 2),
('Vodjenje sudskih postupaka u ...', 'description', 'code', NOW(), NOW(), 21, 2),
('Vodjenje sudskih postupaka u ...', 'description', 'code', NOW(), NOW(), 22, 2),
('Vodjenje sudskih postupaka u ...', 'description', 'code', NOW(), NOW(), 23, 2),
('Vodjenje sudskih postupaka u ...', 'description', 'code', NOW(), NOW(), 24, 2),
('Vodjenje sudskih postupaka u ...', 'description', 'code', NOW(), NOW(), 25, 2),
('Vodjenje sudskih postupaka u ...', 'description', 'code', NOW(), NOW(), 26, 2);


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
('036 - NAKNADA STECAJNOG UPRAVNIKA', '036', 't', 0, 36, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 9, 15, 0, NOW(), NOW(), 0, 0, 0, 'f', 0, 0, 0, 0),
('047 - NAKNADA SAMO POREZ', '047', 't', 15, 36, 30, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, NOW(), NOW(), 0, 0, 0, 'f', 0, 0, 0, 0),
('Depozit - parnica', '047 - Depozit', 't', 15, 47, 30, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, NOW(), NOW(), 0, 0, 0, 'f', 0, 0, 0, 0),
('048 - UGOVORENA NAKNADA', '048', 't', 15, 36, 30, 20.5, 36, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, NOW(), NOW(), 0, 0, 0, 'f', 0, 0, 0, 0),
('NAKNADE - SAMO POREZ', '', 't', 0, 47, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 9, 15, 0, NOW(), NOW(), 0, 0, 0, 't', 0, 0, 0, 0),
('UGOVOR O PRIVREMENO POVREMENOM VRSENJU POSLOVA', '', 't', 0, 36, 0, 0, 36, 15, 0, 5.5, 0, 0, 36, 0.5, 0, 0.5, 0, 0.2, 36, 0, 9, 15, 0, NOW(), NOW(), 0, 0, 0, 't', 0, 0, 0, 0),
('NAKNADA SVI DOPRINOSI', '', 't', 0, 36, 0, 0, 36, 15, 0, 5.5, 0, 0, 36, 0.5, 0, 0.5, 0, 0.2, 36, 0, 9, 15, 0, NOW(), NOW(), 0, 0, 0, 't', 0, 0, 0, 0);
