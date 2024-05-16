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
('Model knjiženja zarada', 'salaries', NOW(), NOW()),
('Model knjiženja naloga za plaćanje', 'payment_orders', NOW(), NOW()),
('Model knjiženja prinudnih naplata', 'enforced_payments', NOW(), NOW()),
('Model knjiženja povraćaja prinudnih naplata', 'return_enforced_payment', NOW(), NOW());

INSERT INTO model_of_accounting_items(title, model_id, debit_account_id, credit_account_id, created_at, updated_at) VALUES
-- knjizenje racuna
('Korektivni račun', 1, 1310, 0, NOW(), NOW()),
('Dobavljač', 1, 0, 1310, NOW(), NOW()),
-- knjizenje resenja
('Korektivni račun', 2, 1310, 0, NOW(), NOW()),
('Dobavljač', 2, 0, 1310, NOW(), NOW()),
('Porez', 2, 0, 1310, NOW(), NOW()),
('Prirez', 2, 0, 1310, NOW(), NOW()),
('Doprinos za PIO', 2, 0, 1310, NOW(), NOW()),
('Doprinos za nezaposlenost', 2, 0, 1310, NOW(), NOW()),
('Doprinos za PIO (zaposleni)', 2, 0, 1310, NOW(), NOW()),
('Doprinos za nezaposlenost (zaposleni)', 2, 0, 1310, NOW(), NOW()),
('Doprinos za PIO (poslodavac)', 2, 0, 1310, NOW(), NOW()),
('Doprinos za nezaposlenost (poslodavac)', 2, 0, 1310, NOW(), NOW()),
('Doprinos za Fond rada', 2, 0, 1310, NOW(), NOW()),
--knjizenje ugovora
('Korektivni račun', 3, 1310, 0, NOW(), NOW()),
('Dobavljač', 3, 0, 1310, NOW(), NOW()),
('Porez', 3, 0, 1310, NOW(), NOW()),
('Prirez', 3, 0, 1310, NOW(), NOW()),
('Doprinos za PIO', 2, 0, 1310, NOW(), NOW()),
('Doprinos za nezaposlenost', 2, 0, 1310, NOW(), NOW()),
('Doprinos za PIO (zaposleni)', 3, 0, 1310, NOW(), NOW()),
('Doprinos za nezaposlenost (zaposleni)', 3, 0, 1310, NOW(), NOW()),
('Doprinos za PIO (poslodavac)', 3, 0, 1310, NOW(), NOW()),
('Doprinos za nezaposlenost (poslodavac)', 3, 0, 1310, NOW(), NOW()),
('Doprinos za Fond rada', 3, 0, 1310, NOW(), NOW()),
--knjizenje zarada
('Korektivni račun', 4, 1310, 0, NOW(), NOW()),
('Doprinos za PIO (zaposleni)', 4, 0, 1310, NOW(), NOW()),
('Doprinos za nezaposlenost (zaposleni)', 4, 0, 1310, NOW(), NOW()),
('Doprinos za PIO (poslodavac)', 4, 0, 1310, NOW(), NOW()),
('Doprinos za nezaposlenost (poslodavac)', 4, 0, 1310, NOW(), NOW()),
('Doprinos za Fond rada', 4, 0, 1310, NOW(), NOW()),
('Porez', 4, 0, 1310, NOW(), NOW()),
('Prirez', 4, 0, 1310, NOW(), NOW()),
('Banka', 4, 0, 1310, NOW(), NOW()),
('Obustave', 4, 0, 1310, NOW(), NOW()),
--knjizenje naloga za placanje
('Dobavljač', 5, 1310, 0, NOW(), NOW()),
('Korektivni račun', 5, 0, 1310, NOW(), NOW()),
('Izdatak', 5, 1310, 0, NOW(), NOW()),
('Rezervisana sredstva', 5, 0, 1310, NOW(), NOW()),
--knjizenje prinudnih naplata
('Dobavljač', 6, 1310, 0, NOW(), NOW()),
('Trošak izvršenja', 6, 1310, 0, NOW(), NOW()),
('Trošak advokata', 6, 1310, 0, NOW(), NOW()),
('Prinudna naplata', 6, 0, 1310, NOW(), NOW()),
--knjizenje povracaja
('Prinudna naplata', 7, 1310, 0, NOW(), NOW()),
('Dobavljač', 7, 0, 1310, NOW(), NOW());

