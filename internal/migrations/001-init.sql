CREATE TABLE currencies (
    code TEXT PRIMARY KEY,
    name TEXT NOT NULL
);

INSERT INTO currencies (code, name) VALUES
    ('RUB', 'Russian Ruble'),
    ('USD', 'US Dollar'),
    ('EUR', 'Euro'),
    ('JPY', 'Japanese Yen')

CREATE TABLE rates (
    id SERIAL PRIMARY KEY,
    currency TEXT REFERENCES currencies(code),
    date DATE NOT NULL,
    nominal INT NOT NULL,
    value_scaled BIGINT NOT NULL,
    UNIQUE(currency, date)
    );

INSERT INTO rates(currency, date, nominal, value_scaled) VALUES
('RUB', '2025-10-21', 1, 10000),
('USD', '2025-10-21', 1, 751234),
('EUR', '2025-10-21', 1, 805678),
('JPY', '2025-10-21', 100, 501234)

