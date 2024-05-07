CREATE TABLE IF NOT EXISTS model_of_accounting_items (
    id serial PRIMARY KEY,
    title VARCHAR ( 255 ) NOT NULL,
    model_id INTEGER NOT NULL REFERENCES models_of_accounting(id) ON DELETE CASCADE,
    debit_account_id INTEGER,
    credit_account_id INTEGER,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);


INSERT INTO models_of_accounting(title, type, created_at, updated_at) VALUES
('Model knjiženja računa', 'invoices', NOW(), NOW()),
('Model knjiženja rješenja', 'decisions', NOW(), NOW()),
('Model knjiženja ugovora', 'contracts', NOW(), NOW()),
('Model knjiženja zarada', 'salaries', NOW(), NOW());

INSERT INTO model_of_accounting_items(title, model_id, debit_account_id, credit_account_id, created_at, updated_at) VALUES
('Korektivni račun', 1, 1175, 0, NOW(), NOW()),
('Dobavljač', 1, 0, 1175, NOW(), NOW()),
('Korektivni račun', 2, 1175, 0, NOW(), NOW()),
('Dobavljač', 2, 0, 1175, NOW(), NOW()),
('Porez', 2, 0, 1175, NOW(), NOW()),
('Prirez', 2, 0, 1175, NOW(), NOW()),
('Korektivni račun', 3, 1175, 0, NOW(), NOW()),
('Dobavljač', 3, 0, 1175, NOW(), NOW()),
('Porez', 3, 0, 1175, NOW(), NOW()),
('Prirez', 3, 0, 1175, NOW(), NOW()),
('Korektivni račun', 4, 1175, 0, NOW(), NOW()),
('Doprinos za PIO (zaposleni)', 4, 0, 1175, NOW(), NOW()),
('Doprinos za nezaposlenost (zaposleni)', 4, 0, 1175, NOW(), NOW()),
('Doprinos za PIO (poslodavac)', 4, 0, 1175, NOW(), NOW()),
('Doprinos za nezaposlenost (poslodavac)', 4, 0, 1175, NOW(), NOW()),
('Doprinos za Fond rada', 4, 0, 1175, NOW(), NOW()),
('Porez', 4, 0, 1175, NOW(), NOW()),
('Prirez', 4, 0, 1175, NOW(), NOW()),
('Banka', 4, 0, 1175, NOW(), NOW()),
('Obustave', 4, 0, 1175, NOW(), NOW());

