delete from model_of_accounting_items;

INSERT INTO model_of_accounting_items(title, model_id, debit_account_id, credit_account_id, created_at, updated_at) VALUES
-- knjizenje racuna
('Korektivni račun', 1, 1, 0, NOW(), NOW()),
('Dobavljač', 1, 0, 1, NOW(), NOW()),
-- knjizenje resenja
('Korektivni račun', 2, 1, 0, NOW(), NOW()),
('Dobavljač', 2, 0, 1, NOW(), NOW()),
('Porez', 2, 0, 1, NOW(), NOW()),
('Prirez', 2, 0, 1, NOW(), NOW()),
('Doprinos za PIO', 2, 0, 1, NOW(), NOW()),
('Doprinos za nezaposlenost', 2, 0, 1, NOW(), NOW()),
('Doprinos za PIO (zaposleni)', 2, 0, 1, NOW(), NOW()),
('Doprinos za nezaposlenost (zaposleni)', 2, 0, 1, NOW(), NOW()),
('Doprinos za PIO (poslodavac)', 2, 0, 1, NOW(), NOW()),
('Doprinos za nezaposlenost (poslodavac)', 2, 0, 1, NOW(), NOW()),
('Doprinos za Fond rada', 2, 0, 1, NOW(), NOW()),
--knjizenje ugovora
('Korektivni račun', 3, 1, 0, NOW(), NOW()),
('Dobavljač', 3, 0, 1, NOW(), NOW()),
('Porez', 3, 0, 1, NOW(), NOW()),
('Prirez', 3, 0, 1, NOW(), NOW()),
('Doprinos za PIO', 2, 0, 1, NOW(), NOW()),
('Doprinos za nezaposlenost', 2, 0, 1, NOW(), NOW()),
('Doprinos za PIO (zaposleni)', 3, 0, 1, NOW(), NOW()),
('Doprinos za nezaposlenost (zaposleni)', 3, 0, 1, NOW(), NOW()),
('Doprinos za PIO (poslodavac)', 3, 0, 1, NOW(), NOW()),
('Doprinos za nezaposlenost (poslodavac)', 3, 0, 1, NOW(), NOW()),
('Doprinos za Fond rada', 3, 0, 1, NOW(), NOW()),
--knjizenje zarada
('Korektivni račun', 4, 1, 0, NOW(), NOW()),
('Doprinos za PIO (zaposleni)', 4, 0, 1, NOW(), NOW()),
('Doprinos za nezaposlenost (zaposleni)', 4, 0, 1, NOW(), NOW()),
('Doprinos za PIO (poslodavac)', 4, 0, 1, NOW(), NOW()),
('Doprinos za nezaposlenost (poslodavac)', 4, 0, 1, NOW(), NOW()),
('Doprinos za Fond rada', 4, 0, 1, NOW(), NOW()),
('Porez', 4, 0, 1, NOW(), NOW()),
('Prirez', 4, 0, 1, NOW(), NOW()),
('Banka', 4, 0, 1, NOW(), NOW()),
('Obustave', 4, 0, 1, NOW(), NOW()),
--knjizenje naloga za placanje
('Dobavljač', 5, 1, 0, NOW(), NOW()),
('Korektivni račun', 5, 0, 1, NOW(), NOW()),
('Izdatak', 5, 1, 0, NOW(), NOW()),
('Rezervisana sredstva', 5, 0, 1, NOW(), NOW()),
--knjizenje prinudnih naplata
('Dobavljač', 6, 1, 0, NOW(), NOW()),
('Trošak izvršenja', 6, 1, 0, NOW(), NOW()),
('Trošak advokata', 6, 1, 0, NOW(), NOW()),
('Naknada za Centralnu banku', 6, 1, 0, NOW(), NOW()),
('Prinudna naplata', 6, 0, 1, NOW(), NOW()),
--knjizenje povracaja
('Prinudna naplata', 7, 1, 0, NOW(), NOW()),
('Dobavljač', 7, 0, 1, NOW(), NOW());