CREATE TABLE IF NOT EXISTS items (
    id INTEGER PRIMARY KEY,
    numbers INTEGER[] NOT NULL -- лучше чем делать таблицу для массива
);