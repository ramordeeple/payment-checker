CREATE TABLE currencies
(
    code     TEXT PRIMARY KEY, -- "USD"
    name     TEXT NOT NULL,    -- "Доллар США"
    num_code TEXT NOT NULL,    -- "840"
    cbr_id   TEXT NOT NULL     -- "R01235"
);

INSERT INTO currencies (code, name_ru, num_code, cbr_id)
VALUES ('RUB', 'Российский рубль', '643', 'R00000'),
       ('USD', 'Доллар США', '840', 'R01235'),
       ('EUR', 'Евро', '978', 'R01239'),
       ('JPY', 'Японская иена', '392', 'R01820');


CREATE TABLE rates
(
    id           SERIAL PRIMARY KEY,
    currency     TEXT REFERENCES currencies (code),
    date         DATE   NOT NULL,
    nominal      INT    NOT NULL,
    value_scaled BIGINT NOT NULL,
    UNIQUE (currency, date)
);


INSERT INTO rates(currency, date, nominal, value_scaled)
VALUES ('RUB', '2025-10-21', 1, 10000),
       ('USD', '2025-10-21', 1, 751234),
       ('EUR', '2025-10-21', 1, 805678),
       ('JPY', '2025-10-21', 100, 501234);
